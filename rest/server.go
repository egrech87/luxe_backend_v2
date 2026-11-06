package rest

import (
	"Quote/config"
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(handler http.Handler, addr string) error {
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func NewRouter() *gin.Engine {
	cfg := config.GetInstance()

	gin.SetMode(cfg.GinMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		// AllowOrigins: []string{"https://*"},
		AllowMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "Selector"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	return router
}
