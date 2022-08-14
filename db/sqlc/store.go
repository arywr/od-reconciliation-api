package db

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/arywr/od-reconciliation-api/helper"
	_ "github.com/arywr/od-reconciliation-api/helper"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/nulls"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction Error %s %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type ProgressEventTypeFromParent struct {
	ID                           int64  `json:"id"`
	ProgressEventTypeName        string `json:"progress_event_type_name"`
	ProgressEventTypeDescription string `json:"progress_event_type_description"`
}

type UpdateProgressTxParams struct {
	ID         int64   `json:"progress_event_id"`
	Percentage float64 `json:"percentage" binding:"min=0"`
	Status     string  `json:"status"`
}

type UpdateProgressTxResult struct {
	ID                int64                       `json:"id"`
	ProgressEventType ProgressEventTypeFromParent `json:"progress_event_type"`
	ProgressName      string                      `json:"progress_name"`
	Status            string                      `json:"status"`
	Percentage        float64                     `json:"percentage"`
	File              string                      `json:"file"`
	CreatedAt         time.Time                   `json:"created_at"`
	UpdatedAt         time.Time                   `json:"updated_at"`
	DeletedAt         nulls.Time                  `json:"deleted_at"`
}

func (store *Store) UpdateProgressTx(ctx context.Context, arg UpdateProgressTxParams) (UpdateProgressTxResult, error) {
	var result UpdateProgressTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		args := UpdateProgressParams{
			ID:         arg.ID,
			Percentage: arg.Percentage,
			Status:     arg.Status,
		}

		progress, err := q.UpdateProgress(ctx, args)
		if err != nil {
			return err
		}

		evtType, err := q.ViewProgressEventType(ctx, int64(progress.ProgressEventTypeID))
		if err != nil {
			return err
		}

		result.ID = progress.ID
		result.ProgressName = progress.ProgressName
		result.ProgressEventType.ID = evtType.ID
		result.ProgressEventType.ProgressEventTypeName = evtType.ProgressEventTypeName
		result.ProgressEventType.ProgressEventTypeDescription = evtType.ProgressEventTypeDescription
		result.Status = progress.Status
		result.Percentage = progress.Percentage
		result.File = progress.File
		result.CreatedAt = progress.CreatedAt
		result.UpdatedAt = progress.UpdatedAt
		result.DeletedAt = progress.DeletedAt

		return nil
	})

	return result, err
}

type ProductTrxCSVParams struct {
	Day        int                   `form:"day"`
	PlatformId string                `form:"platform_id" binding:"required"`
	File       *multipart.FileHeader `form:"file"`
}

type SaveTrxCSVRequest struct {
	FileName   string
	PlatformID string
	ProgressID int64
	Counter    int64
}

func (store *Store) CreateProductTransactionCSV(ctx *gin.Context, args SaveTrxCSVRequest) {
	currentTime := time.Now()

	csvReader, csvFile, err := helper.ReadCSVFile(args.FileName)
	if err != nil {
		log.Println(err)
		return
	}

	defer csvFile.Close()

	isHeader := true
	var rows []CreateProductTransactionParams

	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		if isHeader {
			isHeader = false
			continue
		}

		trxDate, _ := time.Parse("2006-01-02", row[2])
		trxDatetime, _ := time.Parse("2006-01-02 15:04:05", row[2]+" "+row[3])
		collect, _ := strconv.ParseFloat(row[11], 64)
		settled, _ := strconv.ParseFloat(row[11], 64)

		transaction := CreateProductTransactionParams{
			OwnerID:              args.PlatformID,
			TransactionStatusID:  1,
			TransactionTypeID:    1,
			ProductTransactionID: nulls.String{String: row[4], Valid: true},
			TransactionDate:      trxDate,
			TransactionDatetime:  trxDatetime,
			CollectedAmount:      collect,
			SettledAmount:        settled,
			CreatedAt:            currentTime,
			UpdatedAt:            currentTime,
		}
		rows = append(rows, transaction)
	}

	jobs := generateIndex(rows)
	worker := runtime.NumCPU()
	result := TestingInsert(store, ctx, jobs, worker)

	counterSuccess := 0
	for res := range result {
		if res.ID == 0 {
			log.Println("Has Error")
		} else {
			counterSuccess++
		}

		if counterSuccess%100 == 0 {
			progressArgs := UpdateProgressParams{
				ID:         args.ProgressID,
				Percentage: float64(counterSuccess) / float64(args.Counter) * 100,
			}

			store.Queries.UpdateProgress(ctx, progressArgs)
		}
	}

	progressArgs := UpdateProgressParams{
		ID:         args.ProgressID,
		Status:     "completed",
		Percentage: 100,
	}

	store.Queries.UpdateProgress(ctx, progressArgs)
}

func generateIndex(store []CreateProductTransactionParams) <-chan CreateProductTransactionParams {
	result := make(chan CreateProductTransactionParams)

	go func() {
		for _, transaction := range store {
			result <- transaction
		}

		close(result)
	}()

	return result
}

func TestingInsert(
	store *Store,
	ctx *gin.Context,
	jobs <-chan CreateProductTransactionParams,
	worker int,
) <-chan ProductTransaction {
	result := make(chan ProductTransaction)

	wg := new(sync.WaitGroup)
	wg.Add(worker)

	go func() {
		for i := 0; i < worker; i++ {
			go func() {
				for job := range jobs {
					response := InsertFromExcel(store, ctx, job)
					result <- response
				}
				wg.Done()
			}()
		}
	}()

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func InsertFromExcel(store *Store, ctx *gin.Context, transaction CreateProductTransactionParams) ProductTransaction {
	var response ProductTransaction

	for {
		var outerError error
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					*outerError = fmt.Errorf("%v", err)
					log.Println(err)
				}
			}()

			store.execTx(ctx, func(q *Queries) error {
				txRes, errors := q.CreateProductTransaction(ctx, transaction)
				if errors != nil {
					return errors
				}

				response = txRes
				return nil
			})
		}(&outerError)
		if outerError == nil {
			break
		}
	}

	return response
}
