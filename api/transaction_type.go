package api

import (
	"database/sql"
	"net/http"

	db "github.com/arywr/od-reconciliation-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	TypeName        string `json:"type_name" binding:"required"`
	TypeDescription string `json:"type_description" binding:"required"`
}

func (server *Server) createTransactionType(ctx *gin.Context) {
	var request CreateRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateTransactionTypeParams{
		TypeName:        request.TypeName,
		TypeDescription: request.TypeDescription,
	}

	data, err := server.store.CreateTransactionType(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, data)
}

type ViewRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) viewTransactionType(ctx *gin.Context) {
	var request ViewRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data, err := server.store.ViewTransactionType(ctx, request.ID)
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

type FetchRequest struct {
	PageID   int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"size" binding:"required,min=10,max=20"`
}

func (server *Server) allTransactionType(ctx *gin.Context) {
	var request FetchRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.AllTransactionTypeParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}

	data, err := server.store.AllTransactionType(ctx, args)
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

type UpdateRequest struct {
	ID              int64  `uri:"id" binding:"required,min=1"`
	TypeName        string `json:"type_name"`
	TypeDescription string `json:"type_description"`
}

func (server *Server) updateTransactionType(ctx *gin.Context) {
	var request UpdateRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateTransactionTypeParams{
		ID:              request.ID,
		TypeName:        request.TypeName,
		TypeDescription: request.TypeDescription,
	}

	data, err := server.store.UpdateTransactionType(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, data)
}

type DeleteRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteTransactionType(ctx *gin.Context) {
	var request DeleteRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteTransactionType(ctx, request.ID)
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
