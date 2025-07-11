package config

import (
	"os"

	"github.com/viniciusabreusouza/url-shortener/internal/config/logger"
	"github.com/viniciusabreusouza/url-shortener/internal/schemas"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeSQLite() (*gorm.DB, error) {
	logger.Log.Info("Initializing SQLite database")

	dbPath := "./db/shortener.db"
	logger.Log.Info("Database path", zap.String("path", dbPath))

	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		logger.Log.Info("Database not found, creating new one")

		err = os.Mkdir("./db", os.ModePerm)
		if err != nil {
			logger.Log.Error("Failed to create database directory", zap.Error(err))
			return nil, err
		}

		file, err := os.Create(dbPath)
		if err != nil {
			logger.Log.Error("Failed to create database file", zap.Error(err))
			return nil, err
		}

		file.Close()
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		logger.Log.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	err = db.AutoMigrate(&schemas.ShortedUrl{})
	if err != nil {
		logger.Log.Error("Failed to migrate database", zap.Error(err))
		return nil, err
	}

	return db, nil
}
