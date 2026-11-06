package serverUtils

import (
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"go.uber.org/zap"
)

const (
	// Localhost
	localHostMode   = "LOCAL"
	localPathPrefix = "../"

	// Google Cloud Platform
	gcpHostMode   = "GCP"
	gcpPathPrefix = ""
)

func IsInLocalhost() (bool, error) {
	ex, err := os.Executable()
	if err != nil {
		return false, err
	}
	exPath := filepath.Dir(ex)

	return strings.HasPrefix(exPath, "c:\\source\\"), nil
}

// Returns the real path within the server
func GetServerPath(hostMode, path string) string {
	isLocalHost, _ := IsInLocalhost()

	pathPrefix := localPathPrefix
	if !isLocalHost {
		switch hostMode {
		case gcpHostMode:
			pathPrefix = gcpPathPrefix
		}
	}
	return pathPrefix + path
}

func KeepAliveWithSignals(logger *zap.SugaredLogger, serversShutdown func()) {
	quit := make(chan os.Signal, 1)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var osSignalMsg string

	switch <-quit {
	case syscall.SIGINT:
		osSignalMsg = "signal interrupt triggered."

	case syscall.SIGTERM:
		osSignalMsg = "signal termination triggered"
	}

	logger.Debugw("service exited due to os signal", "os_signal", osSignalMsg)
	logger.Info("shutting down servers...")
	serversShutdown()
	logger.Info("servers stopped")
}
