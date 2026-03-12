package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func TestRule4() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	password := "123"
	token := "abc"
	apiKey := "xyz"

	slog.Info("user password: " + password) // want "log message must not contain sensitive data"
	logger.Debug("api_key=" + apiKey)       // want "log message must not contain sensitive data"
	slog.Info("token: " + token)            // want "log message must not contain sensitive data"
	sugar.Error("secret: 404")              // want "log message must not contain sensitive data"

	slog.Info("user authenticated successfully")
	slog.Debug("api request completed")
	slog.Info("token validated")
	slog.Warn("checking password strength")
}
