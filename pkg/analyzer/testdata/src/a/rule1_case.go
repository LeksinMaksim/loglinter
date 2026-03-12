package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func TestRule1() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	slog.Info("Starting server on port 8080") // want "log message must start with a lowercase letter"
	logger.Error("Failed to connect")         // want "log message must start with a lowercase letter"

	slog.Info("starting server on port 8080")
	sugar.Error("failed to connect")

	slog.Info("123 workers started")
	slog.Info(" http server run")
}
