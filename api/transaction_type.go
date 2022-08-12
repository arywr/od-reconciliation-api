package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/arywr/od-reconciliation-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateTypeRequest struct {
	TypeName        string `json:"type_name" binding:"required"`
	TypeDescription string `json:"type_description" binding:"required"`
}

func (server *Server) createTransactionType(ctx *gin.Context) {
	var request CreateTypeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.CreateTransactionTypeParams{
		TypeName:        request.TypeName,
		TypeDescription: request.TypeDescription,
	}

	data, err := server.store.CreateTransactionType(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := APIResponse(http.StatusOK, "OK", data)
	ctx.JSON(http.StatusOK, response)
}

type ViewTypeRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) viewTransactionType(ctx *gin.Context) {
	var request ViewTypeRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	data, err := server.store.ViewTransactionType(ctx, request.ID)
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

type FetchTypeRequest struct {
	PageID   int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"size" binding:"required,min=10,max=100"`
}

func (server *Server) allTransactionType(ctx *gin.Context) {
	var request FetchTypeRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.AllTransactionTypeParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}

	data, err := server.store.AllTransactionType(ctx, args)
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

type UpdateTypeRequest struct {
	ID              int64  `uri:"id" binding:"required,min=1"`
	TypeName        string `json:"type_name"`
	TypeDescription string `json:"type_description"`
}

func (server *Server) updateTransactionType(ctx *gin.Context) {
	var request UpdateTypeRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, APIErrorResponse(http.StatusBadRequest, "ERROR", err))
		return
	}

	args := db.UpdateTransactionTypeParams{
		ID:              request.ID,
		TypeName:        request.TypeName,
		TypeDescription: request.TypeDescription,
	}

	data, err := server.store.UpdateTransactionType(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := APIResponse(http.StatusOK, "OK", data)
	ctx.JSON(http.StatusOK, response)
}

type DeleteTypeRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteTransactionType(ctx *gin.Context) {
	var request DeleteTypeRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	err := server.store.DeleteTransactionType(ctx, request.ID)
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
