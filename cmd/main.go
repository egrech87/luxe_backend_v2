package main

import (
	"Quote/config"
	"Quote/db"
	"Quote/logger"
	"Quote/rest"
	"Quote/rest/api"
	"Quote/rest/internalApi"
	"Quote/utils/pathUtils"
	"Quote/utils/serverUtils"
	"Quote/utils/timeUtils"
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.InitConfigFromEnvYaml()
	if err != nil {
		log.Fatalf("failed to create config : %s", err.Error())
	}

	logger, err := logger.Init(cfg)
	if err != nil {
		log.Fatalf("failed to create logger : %s", err.Error())
	}
	defer logger.Sync()

	logger.Info("service initialising")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ContextTimeout)*time.Second)
	defer cancel()

	_, err = timeUtils.LoadLocation("Europe/Malta")
	if err != nil {
		logger.Warnf("Cron job - Failed to create location: %v", err)
		return
	}

	if err = pathUtils.Init(nil); err != nil {
		logger.Fatalf("failed to initialise the root path: %s", err.Error())
	}

	if err = db.InitRepo(); err != nil {
		logger.Fatalf("failed to connect to database %s : %s", cfg.PostgresDatabase, err.Error())
	}

	// Web server
	webRouter := api.NewPublicAPIRouter()
	webServer := new(rest.Server)
	go func() {
		if err := webServer.Run(webRouter, cfg.ServerAddress); err != nil && err != http.ErrServerClosed {
			logger.Infof("error occured while running web server : %s", err.Error())
		}
	}()

	// Internal Server
	internalRouter := internalApi.NewInternalAPIRouter()
	internalServer := new(rest.Server)
	go func() {
		if err := internalServer.Run(internalRouter, cfg.InternalServerAddress); err != nil &&
			err != http.ErrServerClosed {
			logger.Infof("error occured while running internal server : %s", err.Error())
		}
	}()

	serverUtils.KeepAliveWithSignals(logger, func() {
		webServer.Shutdown(ctx)
		internalServer.Shutdown(ctx)
	})

	logger.Info("successfully uninitialised")
}
