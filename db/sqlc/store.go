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
	"strings"
	"sync"
	"time"

	"github.com/arywr/od-reconciliation-api/helper"
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
	PlatformId int32                 `form:"platform_id" binding:"required"`
	File       *multipart.FileHeader `form:"file"`
}

type SaveTrxCSVRequest struct {
	FileName   string
	PlatformID int32
	ProgressID int64
	Counter    int64
}

func (store *Store) CreateProductTransactionCSV(ctx *gin.Context, args SaveTrxCSVRequest) {
	currentTime := time.Now()

	model, err := store.Queries.GetTransactionModel(ctx, args.PlatformID)

	var rows []CreateProductTransactionParams

	if err == nil {
		switch model.FileType.String {
		case "XLSX":
			readResult, err := helper.ReadExcelFile(args.FileName)

			if err != nil {
				log.Println(err)
				return
			}

			for i := model.RowStartAt.Int16; int(i) < len(readResult); i++ {
				trxDate, _ := time.Parse(model.TransactionDateFormat.String, readResult[i][model.TransactionDate.Int16])
				trxDatetime, _ := time.Parse(model.TransactionDatetimeFormat.String, readResult[i][model.TransactionDatetime.Int16])
				collect, _ := strconv.ParseFloat(readResult[i][model.TransactionAmount.Int16], 64)
				settled, _ := strconv.ParseFloat(readResult[i][model.SettledAmount.Int16], 64)

				transaction := CreateProductTransactionParams{
					ProductID:            args.PlatformID,
					TransactionStatusID:  1,
					TransactionTypeID:    1,
					ProductTransactionID: nulls.String{String: readResult[i][model.TransactionID.Int16], Valid: true},
					TransactionDate:      trxDate,
					TransactionDatetime:  trxDatetime,
					ChannelCode:          nulls.String{String: readResult[i][model.ChannelCode.Int16], Valid: true},
					ChannelName:          nulls.String{String: readResult[i][model.ChannelName.Int16], Valid: true},
					MerchantCode:         nulls.String{String: readResult[i][model.MerchantCode.Int16], Valid: true},
					MerchantName:         nulls.String{String: readResult[i][model.MerchantName.Int16], Valid: true},
					ProductCode:          nulls.String{String: readResult[i][model.ProductCode.Int16], Valid: true},
					ProductName:          nulls.String{String: readResult[i][model.ProductName.Int16], Valid: true},
					CollectedAmount:      collect,
					SettledAmount:        settled,
					CreatedAt:            currentTime,
					UpdatedAt:            currentTime,
				}
				rows = append(rows, transaction)
			}

		case "CSV":
			csvReader, csvFile, err := helper.ReadCSVFile(args.FileName)
			if err != nil {
				log.Println(err)
				return
			}

			defer csvFile.Close()

			isHeader := true

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

				trxDate, _ := time.Parse("2006-01-02", row[model.TransactionDate.Int16])
				trxDatetime, _ := time.Parse("2006-01-02 15:04:05", row[model.TransactionDatetime.Int16])
				collect, _ := strconv.ParseFloat(row[model.TransactionAmount.Int16], 64)
				settled, _ := strconv.ParseFloat(row[model.SettledAmount.Int16], 64)

				transaction := CreateProductTransactionParams{
					ProductID:            args.PlatformID,
					TransactionStatusID:  1,
					TransactionTypeID:    1,
					ProductTransactionID: nulls.String{String: row[model.TransactionID.Int16], Valid: true},
					TransactionDate:      trxDate,
					TransactionDatetime:  trxDatetime,
					ChannelCode:          nulls.String{String: row[model.ChannelCode.Int16], Valid: true},
					ChannelName:          nulls.String{String: row[model.ChannelName.Int16], Valid: true},
					MerchantCode:         nulls.String{String: row[model.MerchantCode.Int16], Valid: true},
					MerchantName:         nulls.String{String: row[model.MerchantName.Int16], Valid: true},
					ProductCode:          nulls.String{String: row[model.ProductCode.Int16], Valid: true},
					ProductName:          nulls.String{String: row[model.ProductName.Int16], Valid: true},
					CollectedAmount:      collect,
					SettledAmount:        settled,
					CreatedAt:            currentTime,
					UpdatedAt:            currentTime,
				}
				rows = append(rows, transaction)
			}

		default:
			return
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
					log.Println(errors)
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

func (store *Store) CreateMerchantTransactionCSV(ctx *gin.Context, args SaveTrxCSVRequest) {
	currentTime := time.Now()

	model, err := store.Queries.GetTransactionModelMerchant(ctx, args.PlatformID)

	var rows []CreateMerchantTrxParams

	if err == nil {
		switch model.FileType.String {
		case "XLSX":
			readResult, err := helper.ReadExcelFile(args.FileName)

			if err != nil {
				log.Println("Read Error: ", err)
				return
			}

			for i := model.RowStartAt.Int16; int(i) < len(readResult); i++ {
				trxDate, _ := time.Parse(model.TransactionDateFormat.String, readResult[i][model.TransactionDate.Int16])
				trxDatetime, _ := time.Parse(model.TransactionDatetimeFormat.String, readResult[i][model.TransactionDatetime.Int16])
				collect, _ := strconv.ParseFloat(readResult[i][model.TransactionAmount.Int16], 64)
				settled, _ := strconv.ParseFloat(readResult[i][model.SettledAmount.Int16], 64)

				transaction := CreateMerchantTrxParams{
					ProductID:            args.PlatformID,
					TransactionStatusID:  1,
					TransactionTypeID:    1,
					ProductTransactionID: nulls.String{String: readResult[i][model.TransactionID.Int16], Valid: true},
					TransactionDate:      trxDate,
					TransactionDatetime:  trxDatetime,
					ChannelCode:          nulls.String{String: readResult[i][model.ChannelCode.Int16], Valid: true},
					ChannelName:          nulls.String{String: readResult[i][model.ChannelName.Int16], Valid: true},
					MerchantCode:         nulls.String{String: readResult[i][model.MerchantCode.Int16], Valid: true},
					MerchantName:         nulls.String{String: readResult[i][model.MerchantName.Int16], Valid: true},
					ProductCode:          nulls.String{String: readResult[i][model.ProductCode.Int16], Valid: true},
					ProductName:          nulls.String{String: readResult[i][model.ProductName.Int16], Valid: true},
					CollectedAmount:      collect,
					SettledAmount:        settled,
					CreatedAt:            currentTime,
					UpdatedAt:            currentTime,
				}
				rows = append(rows, transaction)
			}

		case "CSV":
			csvReader, csvFile, err := helper.ReadCSVFile(args.FileName)
			if err != nil {
				log.Println(err)
				return
			}

			defer csvFile.Close()

			isHeader := true

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

				trxDate, _ := time.Parse("2006-01-02", row[model.TransactionDate.Int16])
				trxDatetime, _ := time.Parse("2006-01-02 15:04:05", row[model.TransactionDatetime.Int16])
				collect, _ := strconv.ParseFloat(row[model.TransactionAmount.Int16], 64)
				settled, _ := strconv.ParseFloat(row[model.SettledAmount.Int16], 64)

				transaction := CreateMerchantTrxParams{
					ProductID:            args.PlatformID,
					TransactionStatusID:  1,
					TransactionTypeID:    1,
					ProductTransactionID: nulls.String{String: row[model.TransactionID.Int16], Valid: true},
					TransactionDate:      trxDate,
					TransactionDatetime:  trxDatetime,
					ChannelCode:          nulls.String{String: row[model.ChannelCode.Int16], Valid: true},
					ChannelName:          nulls.String{String: row[model.ChannelName.Int16], Valid: true},
					MerchantCode:         nulls.String{String: row[model.MerchantCode.Int16], Valid: true},
					MerchantName:         nulls.String{String: row[model.MerchantName.Int16], Valid: true},
					ProductCode:          nulls.String{String: row[model.ProductCode.Int16], Valid: true},
					ProductName:          nulls.String{String: row[model.ProductName.Int16], Valid: true},
					CollectedAmount:      collect,
					SettledAmount:        settled,
					CreatedAt:            currentTime,
					UpdatedAt:            currentTime,
				}
				rows = append(rows, transaction)
			}

		default:
			return
		}

		jobs := generateIndexM(rows)
		worker := runtime.NumCPU()
		result := TestingInsertM(store, ctx, jobs, worker)

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
}

func generateIndexM(store []CreateMerchantTrxParams) <-chan CreateMerchantTrxParams {
	result := make(chan CreateMerchantTrxParams)

	go func() {
		for _, transaction := range store {
			result <- transaction
		}

		close(result)
	}()

	return result
}

func TestingInsertM(
	store *Store,
	ctx *gin.Context,
	jobs <-chan CreateMerchantTrxParams,
	worker int,
) <-chan MerchantTransaction {
	result := make(chan MerchantTransaction)

	wg := new(sync.WaitGroup)
	wg.Add(worker)

	go func() {
		for i := 0; i < worker; i++ {
			go func() {
				for job := range jobs {
					response := InsertFromExcelM(store, ctx, job)
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

func InsertFromExcelM(store *Store, ctx *gin.Context, transaction CreateMerchantTrxParams) MerchantTransaction {
	var response MerchantTransaction

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
				txRes, errors := q.CreateMerchantTrx(ctx, transaction)
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

type ReconcileRequest struct {
	ProductID  int32
	PlatformID int32
	StartDate  time.Time
	EndDate    time.Time
}

func (store *Store) Reconcile(ctx *gin.Context, args ReconcileRequest) error {
	// Find all products and merchant related
	platforms, err := store.Queries.GetRelatedPlatforms(ctx, int64(args.ProductID))

	if err != nil {
		return err
	}

	for _, platform := range platforms {
		switch platform.ProductHasSub {
		case true:
			// Case if products has sub products
		case false:
			// Case if products has no sub products

			// Get reconcile rules of product & platform related
			// Params: product_id, platform_id

			ruleArgs := GetReconcileRuleParams{
				ProductID:  int64(platform.ProductID),
				PlatformID: int64(platform.MerchantID.Int64),
			}

			rules, err := store.Queries.GetReconcileRule(ctx, ruleArgs)

			if err != nil {
				return err
			}

			if len(rules) == 0 {
				continue
			}

			// Initialize slice of string of query tokens

			var matchQuery []string
			var tokenOperator []string
			var combinedToken []string

			for length, rule := range rules {
				if rule.ProductColumnField.String == "transaction_key" && rule.PlatformColumnField.String == "transaction_key" {
					// If 	-> column used is transaction_key
					// Then -> have join possibility and create token join
				} else {
					if rule.ProductColumnField.String != "" {
						matchQuery = append(matchQuery, "WITH join_recon AS (SELECT t1.id FROM product_transactions t1")

						switch rule.ProductColumnConditions.String {
						case "EQUAL":
							token := fmt.Sprintf("WHERE t1.%s = '%s'", rule.ProductColumnField.String, rule.ProductColumnValue.String)
							tokenOperator = append(tokenOperator, token)

							if length != len(rules) {
								break
							}

							if rule.RuleMandatory.String == "REQUIRED" {
								tokenOperator = append(tokenOperator, "AND")
							} else if rule.RuleMandatory.String == "OPTIONAL" {
								tokenOperator = append(tokenOperator, "OR")
							}
						}

						start := args.StartDate.Format("2006-01-02 15:04:05")
						end := args.EndDate.Format("2006-01-02 15:04:05")

						// Join query tokens
						closingToken := fmt.Sprintf(`
							AND t1.transaction_date BETWEEN '%s' AND '%s')
							UPDATE product_transactions
							SET transaction_status_id = 3, updated_at = now(), platform_id = %d
							FROM join_recon WHERE product_transactions.id = join_recon.id;`, start, end, platform.MerchantID.Int64,
						)

						combinedToken = append(matchQuery, tokenOperator...)
						combinedToken = append(combinedToken, closingToken)

						query := strings.Join(combinedToken, " ")
						store.Queries.db.ExecContext(ctx, query)
					} else {
						matchQuery = append(matchQuery, "WITH join_recon AS (SELECT t1.id FROM merchant_transactions t1")

						switch rule.ProductColumnConditions.String {
						case "EQUAL":
							token := fmt.Sprintf("WHERE t1.%s = '%s'", rule.ProductColumnField.String, rule.ProductColumnValue.String)
							tokenOperator = append(tokenOperator, token)

							if length != len(rules) {
								break
							}

							if rule.RuleMandatory.String == "REQUIRED" {
								tokenOperator = append(tokenOperator, "AND")
							} else if rule.RuleMandatory.String == "OPTIONAL" {
								tokenOperator = append(tokenOperator, "OR")
							}
						}

						start := args.StartDate.Format("2006-01-02 15:04:05")
						end := args.EndDate.Format("2006-01-02 15:04:05")

						// Join query tokens
						closingToken := fmt.Sprintf(`
							t1.transaction_date BETWEEN '%s' AND '%s')
							UPDATE merchant_transactions
							SET transaction_status_id = 3, updated_at = now(), product_id = %d, platform_id = %d
							FROM join_recon WHERE merchant_transactions.id = join_recon.id;`, start, end, platform.ProductID, platform.MerchantID.Int64,
						)

						if len(tokenOperator) != 0 {
							closingToken = ("AND" + string(closingToken))
						} else {
							closingToken = ("WHERE" + string(closingToken))
						}

						combinedToken = append(matchQuery, tokenOperator...)
						combinedToken = append(combinedToken, closingToken)

						query := strings.Join(combinedToken, " ")

						log.Println(query)
						store.Queries.db.ExecContext(ctx, query)
					}
				}
			}
		}
	}

	return nil
}
