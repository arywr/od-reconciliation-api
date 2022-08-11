package api

import (
	"database/sql"
	"net/http"

	db "github.com/arywr/od-reconciliation-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateStatusRequest struct {
	StatusName        string `json:"status_name" binding:"required"`
	StatusDescription string `json:"status_description" binding:"required"`
}

func (server *Server) createTransactionStatus(ctx *gin.Context) {
	var req CreateStatusRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateTransactionStatusParams{
		StatusName:        req.StatusName,
		StatusDescription: req.StatusDescription,
	}

	data, err := server.store.CreateTransactionStatus(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, data)
}

type ViewStatusRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) viewTransactionStatus(ctx *gin.Context) {
	var req ViewStatusRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	data, err := server.store.ViewTransactionStatus(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, data)
}

type FetchStatusRequest struct {
	PageID   int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"size" binding:"required,min=10,max=100"`
}

func (server *Server) allTransactionStatus(ctx *gin.Context) {
	var req FetchStatusRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.AllTransactionStatusParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	data, err := server.store.AllTransactionStatus(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, data)
}

type UpdateStatusRequest struct {
	ID                int64  `uri:"id" binding:"required,min=1"`
	StatusName        string `json:"status_name"`
	StatusDescription string `json:"status_description"`
}

func (server *Server) updateTransactionStatus(ctx *gin.Context) {
	var req UpdateStatusRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateTransactionStatusParams{
		ID:                req.ID,
		StatusName:        req.StatusName,
		StatusDescription: req.StatusDescription,
	}

	data, err := server.store.UpdateTransactionStatus(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, data)
}

type DeleteStatusRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteTransactionStatus(ctx *gin.Context) {
	var request DeleteStatusRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteTransactionStatus(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Succesfully delete")
}
