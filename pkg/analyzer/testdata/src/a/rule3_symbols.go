package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func TestRule3() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	slog.Info("server started! 🚀")                // want "log message must not contain special characters or emojis"
	logger.Error("connection failed!!!")          // want "log message must not contain special characters or emojis"
	slog.Warn("warning: something went wrong...") // want "log message must not contain special characters or emojis"
	sugar.Debug("user created? yes")              // want "log message must not contain special characters or emojis"

	slog.Info("server started")
	slog.Error("connection failed")
	slog.Warn("warning: something went wrong")
	slog.Info("user initialized, starting tasks.")
	slog.Debug("path-to-file is correct")
}
