package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/arywr/od-reconciliation-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateProgressEventRequest struct {
	ProgressEventTypeID int16   `json:"progress_event_type_id" binding:"required,min=1"`
	ProgressName        string  `json:"progress_name" binding:"required"`
	Status              string  `json:"status" binding:"required"`
	Percentage          float64 `json:"percentage" binding:"required,min=0"`
	File                string  `json:"file"`
}

func (server *Server) createProgressEvent(ctx *gin.Context) {
	var req CreateProgressEventRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.CreateProgressEventParams{
		ProgressEventTypeID: req.ProgressEventTypeID,
		ProgressName:        req.ProgressName,
		Status:              req.Status,
		Percentage:          req.Percentage,
		File:                req.File,
	}

	data, err := server.store.CreateProgressEvent(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := APIResponse(http.StatusOK, "OK", data)
	ctx.JSON(http.StatusOK, response)
}

type UpdateProgressEventRequest struct {
	Status     string  `json:"status"`
	Percentage float64 `json:"percentage"`
	ID         int64   `uri:"id" binding:"required,min=1"`
}

func (server *Server) updateProgressEvent(ctx *gin.Context) {
	var req UpdateProgressEventRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.UpdateProgressTxParams{
		ID:         req.ID,
		Percentage: req.Percentage,
		Status:     req.Status,
	}

	data, err := server.store.UpdateProgressTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := APIResponse(http.StatusOK, "OK", data)
	ctx.JSON(http.StatusOK, response)
}

type ViewProgressEventRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) viewProgressEvent(ctx *gin.Context) {
	var req ViewProgressEventRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
		return
	}

	data, err := server.store.ViewProgressEvent(ctx, req.ID)
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

type FetchEventProgressRequest struct {
	PageID   int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"size" binding:"required,min=10,max=100"`
}

func (server *Server) allProgressEvent(ctx *gin.Context) {
	var req FetchEventProgressRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.AllProgressEventParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	data, err := server.store.AllProgressEvent(ctx, args)
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

func (server *Server) deleteProgressEvent(ctx *gin.Context) {
	var req ViewProgressEventRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	err := server.store.DeleteProgressEvent(ctx, req.ID)
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
