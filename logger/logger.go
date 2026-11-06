package logger

import (
	"Quote/config"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func GetInstance() *zap.SugaredLogger {
	return logger
}

func Init(cfg *config.RentalConfig) (*zap.SugaredLogger, error) {
	if logger == nil {
		var err error
		logger, err = NewLogger(cfg)
		return logger, err
	}

	return logger, nil
}

func NewLogger(cfg *config.RentalConfig) (*zap.SugaredLogger, error) {
	// Log level parsing
	logLevel, err := zap.ParseAtomicLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	output := cfg.LogOutput
	if output != "stderr" {
		output = fmt.Sprintf("%s/%s", cfg.RootPath, output)
	}

	zapCFG := zap.Config{
		Encoding:         "json",
		Level:            logLevel,
		OutputPaths:      []string{output},
		ErrorOutputPaths: []string{output},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := zapCFG.Build()
	if err != nil {
		return nil, err
	}

	sugaredLogger := logger.Sugar()

	return sugaredLogger, nil
}
