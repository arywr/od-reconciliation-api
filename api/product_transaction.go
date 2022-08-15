package api

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	db "github.com/arywr/od-reconciliation-api/db/sqlc"
	"github.com/arywr/od-reconciliation-api/helper"
	"github.com/arywr/od-reconciliation-api/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/gobuffalo/nulls"
	"github.com/google/uuid"
)

type CreateProductTrxRequest struct {
	TransactionStatusID   int16        `json:"transaction_status_id" binding:"required"`
	TransactionTypeID     int16        `json:"transaction_type_id" binding:"required"`
	ProgressEventID       int16        `json:"progress_event_id"`
	ProductTransactionID  nulls.String `json:"product_transaction_id" binding:"required"`
	MerchantTransactionID nulls.String `json:"merchant_transaction_id"`
	ChannelTransactionID  nulls.String `json:"channel_transaction_id"`
	OwnerID               string       `json:"owner_id" binding:"required"`
	TransactionID         string       `json:"transaction_id"`
	TransactionDate       string       `json:"transaction_date" binding:"required"`
	TransactionDatetime   string       `json:"transaction_datetime" binding:"required"`
	CollectedAmount       float64      `json:"collected_amount" binding:"required"`
	SettledAmount         float64      `json:"settled_amount" binding:"required"`
	CreatedAt             string       `json:"created_at"`
	UpdatedAt             string       `json:"updated_at"`
	DeletedAt             nulls.Time   `json:"deleted_at"`
}

func (server *Server) createProductTransaction(ctx *gin.Context) {
	var req CreateProductTrxRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	trxDate, err := stringToDate(req.TransactionDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIErrorResponse(http.StatusBadRequest, "ERROR", err))
		return
	}

	trxDatetime, err := stringToDatetime(req.TransactionDatetime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIErrorResponse(http.StatusBadRequest, "ERROR", err))
		return
	}

	args := db.CreateProductTransactionParams{
		TransactionStatusID:   req.TransactionStatusID,
		TransactionTypeID:     req.TransactionTypeID,
		ProductTransactionID:  req.ProductTransactionID,
		MerchantTransactionID: req.MerchantTransactionID,
		ChannelTransactionID:  req.ChannelTransactionID,
		OwnerID:               req.OwnerID,
		TransactionID:         req.TransactionID,
		TransactionDate:       trxDate,
		TransactionDatetime:   trxDatetime,
		CollectedAmount:       req.CollectedAmount,
		SettledAmount:         req.SettledAmount,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nulls.Time{},
	}

	data, err := server.store.CreateProductTransaction(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := APIResponse(http.StatusOK, "OK", data)
	ctx.JSON(http.StatusOK, response)
}

type FetchProductTrxRequest struct {
	PageID   int32 `form:"page" binding:"required,min=1,max=10"`
	PageSize int32 `form:"size" binding:"required,min=10,max=100"`
}

func (server *Server) allProductTransaction(ctx *gin.Context) {
	var req FetchProductTrxRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.AllProductTransactionParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	data, err := server.store.AllProductTransaction(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, APIErrorResponse(http.StatusNotFound, "ERROR", err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := APIResponse(http.StatusOK, "OK", data)
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) allDuplicateProductTransaction(ctx *gin.Context) {
	var req FetchProductTrxRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.AllDuplicateProductTransactionParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	data, err := server.store.AllDuplicateProductTransaction(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, APIErrorResponse(http.StatusNotFound, "ERROR", err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := APIResponse(http.StatusOK, "OK", data)
	ctx.JSON(http.StatusOK, response)
}

type CreateTrxFromCSV struct {
	Day        int                   `form:"day"`
	PlatformId string                `form:"platform_id" binding:"required"`
	File       *multipart.FileHeader `form:"file"`
}

type CreateTrxCSVResult struct {
	FileName      string           `json:"file_name"`
	FileSize      string           `json:"file_size"`
	FileExtension string           `json:"file_extension"`
	RowsReaded    int64            `json:"rows_readed"`
	Progress      db.ProgressEvent `json:"progress_event"`
}

func (server *Server) createTransactionFromCSV(ctx *gin.Context) {
	var req CreateTrxFromCSV

	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.ProductTrxCSVParams{
		Day:        req.Day,
		PlatformId: req.PlatformId,
		File:       req.File,
	}

	unixTime := time.Now().Unix()
	ext := filepath.Ext(req.File.Filename)
	fileName := fmt.Sprintf("upload/%s_%d_%s.%s", uuid.New().String(), unixTime, req.PlatformId, ext)

	if err := ctx.SaveUploadedFile(args.File, fileName); err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	// Read & Count CSV File
	var counter int64

	csvReader, csvFile, err := helper.ReadCSVFile(fileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}
	defer csvFile.Close()

	// Create progress data first
	progressArgs := db.CreateProgressEventParams{
		ProgressEventTypeID: 1,
		ProgressName:        req.PlatformId,
		Status:              "on process",
		Percentage:          0,
		File:                req.File.Filename,
	}

	progressResult, err := server.store.CreateProgressEvent(ctx, progressArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	isHeader := true

	for {
		_, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
				return
			}

			break
		}

		if isHeader {
			isHeader = false
			continue
		}
		counter++
	}

	res := CreateTrxCSVResult{
		FileName:      req.File.Filename,
		FileSize:      util.FormatFileSize(float64(req.File.Size)),
		FileExtension: req.File.Header.Get("Content-Type"),
		RowsReaded:    counter,
		Progress: db.ProgressEvent{
			ID:                  progressResult.ID,
			ProgressName:        progressResult.ProgressName,
			ProgressEventTypeID: progressResult.ProgressEventTypeID,
			Status:              progressResult.Status,
			File:                progressResult.File,
			CreatedAt:           progressResult.CreatedAt,
			UpdatedAt:           progressResult.UpdatedAt,
			DeletedAt:           progressResult.DeletedAt,
		},
	}

	mainArgs := db.SaveTrxCSVRequest{
		ProgressID: progressResult.ID,
		FileName:   fileName,
		PlatformID: req.PlatformId,
		Counter:    counter,
	}

	response := APIResponse(http.StatusOK, "OK", res)
	ctx.JSON(http.StatusOK, response)

	go func() {
		server.store.CreateProductTransactionCSV(ctx, mainArgs)
		helper.DestroyFile(fileName)
	}()
}

type DeleteDuplicateTrxRequest struct {
	StartDate  string `json:"start_date" binding:"required"`
	EndDate    string `json:"end_date" binding:"required"`
	PlatformID string `json:"platform_id" binding:"required"`
}

type DeleteDuplicateTrxResult struct {
	RowsAffected int64  `json:"rows_deleted"`
	PlatformID   string `json:"platform_id"`
}

func (server *Server) deleteDuplicateProductTransaction(ctx *gin.Context) {
	var req DeleteDuplicateTrxRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)

	args := db.DeleteDuplicateProductTrxParams{
		StartDate:  start,
		EndDate:    end,
		PlatformID: req.PlatformID,
	}

	data, err := server.store.DeleteDuplicateProductTrx(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, APIErrorResponse(http.StatusNotFound, "ERROR", err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := DeleteDuplicateTrxResult{
		RowsAffected: data,
		PlatformID:   req.PlatformID,
	}

	ctx.JSON(http.StatusOK, APIResponse(http.StatusOK, "OK", response))
}

type DeleteSingleTrxRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteProductTransactionByID(ctx *gin.Context) {
	var req DeleteSingleTrxRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	err := server.store.DeleteProductTrxByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, APIErrorResponse(http.StatusNotFound, "ERROR", err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	ctx.JSON(http.StatusOK, APIResponse(http.StatusOK, "OK", "Successfully delete"))
}
