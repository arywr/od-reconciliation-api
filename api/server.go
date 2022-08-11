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

	apiTypes := router.Group("/api/transaction-types")
	apiStatus := router.Group("/api/transaction-statuses")

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

	server.router = router
	return server
}

// Starting HTTP server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
