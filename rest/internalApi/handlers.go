package internalApi

import (
	"Quote/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type InternalAPIHandler struct {
}

func NewInternalAPIHandler() *InternalAPIHandler {
	return &InternalAPIHandler{}
}

type health struct {
	Status      string `json:"status"`
	ServiceName string `json:"service_name"`
	BootTime    string `json:"boot_time"`
}

func (h *InternalAPIHandler) healthHandler(c *gin.Context) {
	cfg := config.GetInstance()

	healthRes := health{
		Status:      "running",
		ServiceName: cfg.ServiceName,
		BootTime:    cfg.ServiceBootTime.Format(time.RFC1123),
	}

	c.JSON(http.StatusOK, healthRes)
}
