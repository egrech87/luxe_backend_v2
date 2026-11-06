package internalApi

import (
	"Quote/rest"

	"github.com/gin-gonic/gin"
)

func NewInternalAPIRouter() *gin.Engine {
	router := rest.NewRouter()

	handler := NewInternalAPIHandler()

	// Routes
	router.GET("/health", handler.healthHandler)

	return router
}
