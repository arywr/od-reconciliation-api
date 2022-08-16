package api

import (
	db "github.com/arywr/od-reconciliation-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Serve HTTP Request to execute db or transactions
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// Creating a new HTTP server instance
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	apiTypes := router.Group("/api/v1/go-reconciliation/transaction-types")
	apiStatus := router.Group("/api/v1/go-reconciliation/transaction-statuses")
	apiEventTypes := router.Group("/api/v1/go-reconciliation/progress-event-types")
	apiProgressEvents := router.Group("/api/v1/go-reconciliation/progress-events")
	apiProductTrx := router.Group("/api/v1/go-reconciliation/product-transactions")
	apiMerchantTrx := router.Group("/api/v1/go-reconciliation/merchant-transactions")
	apiReconciliation := router.Group("/api/v1/go-reconciliation/reconcile")

	apiTypes.POST("/", server.createTransactionType)
	apiTypes.GET("/:id", server.viewTransactionType)
	apiTypes.GET("/", server.allTransactionType)
	apiTypes.PUT("/:id", server.updateTransactionType)
	apiTypes.DELETE("/:id", server.deleteTransactionType)

	apiStatus.POST("/", server.createTransactionStatus)
	apiStatus.GET("/:id", server.viewTransactionStatus)
	apiStatus.GET("/", server.allTransactionStatus)
	apiStatus.PUT("/:id", server.updateTransactionStatus)
	apiStatus.DELETE("/:id", server.deleteTransactionStatus)

	apiEventTypes.POST("/", server.createProgressEventType)
	apiEventTypes.GET("/:id", server.viewProgressEventType)
	apiEventTypes.GET("/", server.allProgressEventTypeRequest)
	apiEventTypes.PUT("/:id", server.updateProgressEventType)
	apiEventTypes.DELETE("/:id", server.deleteProgressEventType)

	apiProgressEvents.POST("/", server.createProgressEvent)
	apiProgressEvents.PUT("/:id", server.updateProgressEvent)
	apiProgressEvents.GET("/:id", server.viewProgressEvent)
	apiProgressEvents.GET("/", server.allProgressEvent)
	apiProgressEvents.DELETE("/:id", server.deleteProgressEvent)

	apiProductTrx.POST("/", server.createProductTransaction)
	apiProductTrx.POST("/csv", server.createTransactionFromCSV)
	apiProductTrx.GET("/", server.allProductTransaction)
	apiProductTrx.GET("/duplicates", server.allDuplicateProductTransaction)
	apiProductTrx.DELETE("/duplicates", server.deleteDuplicateProductTransaction)
	apiProductTrx.DELETE("/:id", server.deleteProductTransactionByID)

	apiMerchantTrx.POST("/", server.createMerchantTransaction)
	apiMerchantTrx.POST("/csv", server.createMerchantTrxFromCSV)
	apiMerchantTrx.GET("/", server.allMerchantTransaction)
	apiMerchantTrx.GET("/duplicates", server.allDuplicateMerchantTransaction)
	apiMerchantTrx.DELETE("/duplicates", server.deleteDuplicateMerhcantTransaction)
	apiMerchantTrx.DELETE("/:id", server.deleteMerchantTransactionByID)

	apiReconciliation.POST("/", server.reconciliationProduct)

	server.router = router
	return server
}

// Starting HTTP server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
