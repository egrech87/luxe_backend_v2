package main

import (
	"Quote/config"
	"Quote/db"
	"Quote/logger"
	"Quote/rest"
	"Quote/rest/api"
	"Quote/rest/internalApi"
	"Quote/utils/pathUtils"
	"Quote/utils/timeUtils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

type LuxeService struct {
	name           string
	cfg            *config.RentalConfig
	ctx            context.Context
	logger         *zap.SugaredLogger
	webServer      *rest.Server
	internalServer *rest.Server
}

func (m *LuxeService) terminate() {
	m.webServer.Shutdown(m.ctx)
	m.internalServer.Shutdown(m.ctx)
	log.Println("Service shut down")
}

func (m *LuxeService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	status <- svc.Status{State: svc.StartPending}

	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		for c := range r {
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				log.Print("Shutting service...!")
				m.terminate()
				break loop
			default:
				log.Printf("Unexpected service control request #%d", c)
			}
		}
	}

	status <- svc.Status{State: svc.StopPending}
	return false, 0
}

func runService(service *LuxeService, isDebug bool) {
	if isDebug {
		err := debug.Run(service.name, service)
		if err != nil {
			log.Fatalln("Error running service in debug mode.")
		}
	} else {
		err := svc.Run(service.name, service)
		if err != nil {
			log.Fatalln("Error running service in Service Control mode.")
		}
	}
}

func mainService() {
	cfg, err := config.InitConfigFromEnvYaml()
	if err != nil {
		log.Fatalf("failed to create config : %s", err.Error())
	}

	if err := pathUtils.Init(&cfg.RootPath); err != nil {
		log.Fatalf("failed to initialise the root path: %s", err.Error())
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/debug.log", cfg.RootPath), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(fmt.Errorf("error opening file: %v", err))
	}
	defer f.Close()

	log.SetOutput(f)

	loger, err := logger.Init(cfg)
	if err != nil {
		log.Fatalf("failed to create logger : %s", err.Error())
	}
	defer loger.Sync()

	loger.Info("service initialising")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ContextTimeout)*time.Second)
	defer cancel()

	_, err = timeUtils.LoadLocation("Europe/Malta")
	if err != nil {
		loger.Warnf("Cron job - Failed to create location: %v", err)
		return
	}

	if err = db.InitRepo(); err != nil {
		loger.Fatalf("failed to connect to database %s : %s", cfg.PostgresDatabase, err.Error())
	}

	// Web server
	webRouter := api.NewPublicAPIRouter()
	webServer := new(rest.Server)
	go func() {
		if err := webServer.Run(webRouter, cfg.ServerAddress); err != nil && err != http.ErrServerClosed {
			loger.Infof("error occured while running web server : %s", err.Error())
		}
	}()

	// Internal Server
	internalRouter := internalApi.NewInternalAPIRouter()
	internalServer := new(rest.Server)
	go func() {
		if err := internalServer.Run(internalRouter, cfg.InternalServerAddress); err != nil &&
			err != http.ErrServerClosed {
			loger.Infof("error occured while running internal server : %s", err.Error())
		}
	}()

	service := &LuxeService{
		name:           "luxe_service",
		cfg:            cfg,
		ctx:            ctx,
		logger:         loger,
		webServer:      webServer,
		internalServer: internalServer,
	}

	runService(service, false)
}
