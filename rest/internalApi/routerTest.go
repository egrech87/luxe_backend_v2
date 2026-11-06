package internalApi

import (
	"Quote/config"
	"Quote/logger"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	// Initialise Config with defaults
	cfg, err := config.InitConfigFromEnv()
	if err != nil {
		t.Fatalf("failed to read config : %s", err.Error())
	}
	cfg.LogLevel = "fatal"

	// Inititialse logger
	logger, err := logger.NewLogger(cfg)
	if err != nil {
		return
	}
	defer logger.Sync()

	server := httptest.NewServer(NewInternalAPIRouter())
	defer server.Close()

	res, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("failed to request API: %v", err)
	}

	statusCode := res.StatusCode
	if statusCode != 200 {
		t.Fatalf("status code = %v", statusCode)
	}

	var body map[string]any
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read body bytes : %v", err)
	}

	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		t.Fatalf("failed to unmarshal body bytes : %v", err)
	}

	status := body["status"]
	if status != "running" {
		t.Errorf("wrong status obtained, got %s", status.(string))
	}
}
