package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/arywr/od-reconciliation-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateStatusRequest struct {
	StatusName        string `json:"status_name" binding:"required,oneof=unknown partially-match fully-match dispute"`
	StatusDescription string `json:"status_description" binding:"required"`
}

func (server *Server) createTransactionStatus(ctx *gin.Context) {
	var req CreateStatusRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.CreateTransactionStatusParams{
		StatusName:        req.StatusName,
		StatusDescription: req.StatusDescription,
	}

	data, err := server.store.CreateTransactionStatus(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := APIResponse(http.StatusOK, "OK", data)
	ctx.JSON(http.StatusOK, response)
}

type ViewStatusRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) viewTransactionStatus(ctx *gin.Context) {
	var req ViewStatusRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	data, err := server.store.ViewTransactionStatus(ctx, req.ID)
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

type FetchStatusRequest struct {
	PageID   int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"size" binding:"required,min=10,max=100"`
}

func (server *Server) allTransactionStatus(ctx *gin.Context) {
	var req FetchStatusRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.AllTransactionStatusParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	data, err := server.store.AllTransactionStatus(ctx, args)
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

type UpdateStatusRequest struct {
	ID                int64  `uri:"id" binding:"required,min=1"`
	StatusName        string `json:"status_name"`
	StatusDescription string `json:"status_description"`
}

func (server *Server) updateTransactionStatus(ctx *gin.Context) {
	var req UpdateStatusRequest

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

	args := db.UpdateTransactionStatusParams{
		ID:                req.ID,
		StatusName:        req.StatusName,
		StatusDescription: req.StatusDescription,
	}

	data, err := server.store.UpdateTransactionStatus(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := APIResponse(http.StatusOK, "OK", data)
	ctx.JSON(http.StatusOK, response)
}

type DeleteStatusRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteTransactionStatus(ctx *gin.Context) {
	var request DeleteStatusRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	err := server.store.DeleteTransactionStatus(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, APIErrorResponse(http.StatusNotFound, "ERROR", err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	ctx.JSON(http.StatusOK, "Succesfully delete")
}
