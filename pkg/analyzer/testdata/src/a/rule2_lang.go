package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func TestRule2() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	slog.Info("запуск сервера")         // want "log message must be in English only"
	sugar.Error("ошибка подключения")   // want "log message must be in English only"
	slog.Warn("server status: падение") // want "log message must be in English only"

	sugar.Info("server starting")
	slog.Error("connection error")
}
