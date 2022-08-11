package api

import (
	"database/sql"
	"net/http"

	db "github.com/arywr/od-reconciliation-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type CreateEventTypeRequest struct {
	EventName        string `json:"event_name" binding:"required"`
	EventDescription string `json:"event_description" binding:"required"`
}

func (server *Server) createProgressEventType(ctx *gin.Context) {
	var request CreateEventTypeRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateProgressEventTypeParams{
		ProgressEventTypeName:        request.EventName,
		ProgressEventTypeDescription: request.EventDescription,
	}

	data, err := server.store.CreateProgressEventType(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, data)
}

type ViewEventTypeRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) viewProgressEventType(ctx *gin.Context) {
	var request ViewEventTypeRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	data, err := server.store.ViewProgressEventType(ctx, request.ID)
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

type FetchEventTypeRequest struct {
	PageID   int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"size" binding:"required,min=10,max=100"`
}

func (server *Server) allProgressEventTypeRequest(ctx *gin.Context) {
	var request FetchEventTypeRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.AllProgressEventTypeParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}

	data, err := server.store.AllProgressEventType(ctx, args)
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

type UpdateEventTypeRequest struct {
	ID               int64  `uri:"id" binding:"required,min=1"`
	EventName        string `json:"event_name"`
	EventDescription string `json:"event_description"`
}

func (server *Server) updateProgressEventType(ctx *gin.Context) {
	var request UpdateEventTypeRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.UpdateProgressEventTypeParams{
		ID:                           request.ID,
		ProgressEventTypeName:        request.EventName,
		ProgressEventTypeDescription: request.EventDescription,
	}

	data, err := server.store.UpdateProgressEventType(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, data)
}

type DeleteEventTypeRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteProgressEventType(ctx *gin.Context) {
	var request DeleteEventTypeRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteProgressEventType(ctx, request.ID)
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
