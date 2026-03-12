package main

import "log/slog"

func main() {
	password := "123"
	apiKey := "xyz"

	slog.Info("Starting server on port 8080")
	slog.Error("ошибка подключения к бд")
	slog.Info("server started!!")
	slog.Info("user password:" + password)
	slog.Debug("api_key:" + apiKey)
}
