package repository

import (
	"context"

	"github.com/viniciusabreusouza/url-shortener/internal/schemas"
	"gorm.io/gorm"
)

type ShortenRepository interface {
	ShortenUrl(ctx context.Context, shortId string, encryptedUrl string) error
}

type shortenRepository struct {
	db *gorm.DB
}

func NewShortenRepository(db *gorm.DB) ShortenRepository {
	return &shortenRepository{
		db: db,
	}
}

func (r shortenRepository) ShortenUrl(ctx context.Context, shortId string, encryptedUrl string) error {
	return r.db.WithContext(ctx).Create(&schemas.ShortedUrl{
		ShortId:      shortId,
		EncryptedUrl: encryptedUrl,
	}).Error
}
