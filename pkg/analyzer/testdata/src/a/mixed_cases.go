package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func TestMixed() {
	logger, _ := zap.NewProduction()

	slog.Info("Запуск сервера!")  // want "log message must start with a lowercase letter" "log message must be in English only" "log message must not contain special characters or emojis"
	logger.Error("Password: err") // want "log message must start with a lowercase letter" "log message must not contain sensitive data"
}
