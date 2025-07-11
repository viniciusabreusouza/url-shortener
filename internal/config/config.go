package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/viniciusabreusouza/url-shortener/internal/config/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() error {
	if err := godotenv.Load(); err != nil {
		logger.Log.Error("Failed to load environment variables", zap.Error(err))
		return err
	}

	var error error
	db, error = InitializeSQLite()
	if error != nil {
		return fmt.Errorf("failed to initialize database: %w", error)
	}
	return nil
}

func GetDB() *gorm.DB {
	return db
}
