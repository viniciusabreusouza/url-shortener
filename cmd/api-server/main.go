package main

import (
	"github.com/viniciusabreusouza/url-shortener/internal/config"
	"github.com/viniciusabreusouza/url-shortener/internal/config/logger"
	"github.com/viniciusabreusouza/url-shortener/internal/router"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger(true)
	defer logger.Log.Sync()

	logger.Log.Info("Starting Application")

	err := config.Init()
	if err != nil {
		logger.Log.Error("Failed to initialize database", zap.Error(err))
		return
	}

	router.Initialize()
}
