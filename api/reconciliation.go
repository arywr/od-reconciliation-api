package api

import (
	"errors"
	"net/http"

	db "github.com/arywr/od-reconciliation-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReconProductRequest struct {
	OwnerID       string `json:"platform_id" binding:"required"`
	DestinationID string `json:"destination_id" binding:"required"`
	StartDate     string `json:"start_date" binding:"required"`
	EndDate       string `json:"end_date" binding:"required"`
}

func (server *Server) reconciliationProduct(ctx *gin.Context) {
	var req ReconProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var valError validator.ValidationErrors
		if errors.As(err, &valError) {
			ctx.JSON(http.StatusBadRequest, APIValidationResponse(http.StatusBadRequest, "ERROR", valError))
			return
		}
	}

	start, _ := stringToDatetime(req.StartDate)
	end, _ := stringToDatetime(req.EndDate)

	arg1 := db.MatchReconciliationProductParams{
		PlatformID: req.OwnerID,
		StartDate:  start,
		EndDate:    end,
	}

	err := server.store.MatchReconciliationProduct(ctx, arg1)
	if err != nil {
		response := APIErrorResponse(http.StatusInternalServerError, "ERROR", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	arg2 := db.MatchReconciliationMerchantParams{
		DestinationID: req.DestinationID,
		StartDate:     start,
		EndDate:       end,
	}

	err2 := server.store.MatchReconciliationMerchant(ctx, arg2)
	if err2 != nil {
		response := APIErrorResponse(http.StatusInternalServerError, "ERROR", err)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	go func() {
		// Flagging Dispute Transactions
		arg3 := db.DisputeReconciliationProductParams{
			PlatformID: req.OwnerID,
			StartDate:  start,
			EndDate:    end,
		}

		err3 := server.store.DisputeReconciliationProduct(ctx, arg3)
		if err3 != nil {
			response := APIErrorResponse(http.StatusInternalServerError, "ERROR", err)
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		arg4 := db.DisputeReconciliationMerchantParams{
			DestinationID: req.DestinationID,
			StartDate:     start,
			EndDate:       end,
		}

		err4 := server.store.DisputeReconciliationMerchant(ctx, arg4)
		if err4 != nil {
			response := APIErrorResponse(http.StatusInternalServerError, "ERROR", err)
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}
	}()

	response := APIResponse(http.StatusOK, "OK", "Successfully reconcile transaction")
	ctx.JSON(http.StatusOK, response)
}
