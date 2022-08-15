package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/arywr/od-reconciliation-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateEventTypeRequest struct {
	EventName        string `json:"event_name" binding:"required,oneof=read recon others"`
	EventDescription string `json:"event_description" binding:"required"`
}

func (server *Server) createProgressEventType(ctx *gin.Context) {
	var request CreateEventTypeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.CreateProgressEventTypeParams{
		ProgressEventTypeName:        request.EventName,
		ProgressEventTypeDescription: request.EventDescription,
	}

	data, err := server.store.CreateProgressEventType(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := APIResponse(http.StatusOK, "OK", data)
	ctx.JSON(http.StatusOK, response)
}

type ViewEventTypeRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) viewProgressEventType(ctx *gin.Context) {
	var request ViewEventTypeRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	data, err := server.store.ViewProgressEventType(ctx, request.ID)
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

type FetchEventTypeRequest struct {
	PageID   int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"size" binding:"required,min=10,max=100"`
}

func (server *Server) allProgressEventTypeRequest(ctx *gin.Context) {
	var request FetchEventTypeRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.AllProgressEventTypeParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}

	data, err := server.store.AllProgressEventType(ctx, args)
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

type UpdateEventTypeRequest struct {
	ID               int64  `uri:"id" binding:"required,min=1"`
	EventName        string `json:"event_name"`
	EventDescription string `json:"event_description"`
}

func (server *Server) updateProgressEventType(ctx *gin.Context) {
	var request UpdateEventTypeRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	args := db.UpdateProgressEventTypeParams{
		ID:                           request.ID,
		ProgressEventTypeName:        request.EventName,
		ProgressEventTypeDescription: request.EventDescription,
	}

	data, err := server.store.UpdateProgressEventType(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, APIErrorResponse(http.StatusInternalServerError, "ERROR", err))
		return
	}

	response := APIResponse(http.StatusOK, "OK", data)
	ctx.JSON(http.StatusOK, response)
}

type DeleteEventTypeRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteProgressEventType(ctx *gin.Context) {
	var request DeleteEventTypeRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
		return
	}

	err := server.store.DeleteProgressEventType(ctx, request.ID)
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
