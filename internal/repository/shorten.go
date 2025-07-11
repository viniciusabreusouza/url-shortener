package repository

import (
	"context"
	"errors"

	"github.com/viniciusabreusouza/url-shortener/internal/schemas"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("short url not found")

type ShortenRepository interface {
	ShortenUrl(ctx context.Context, shortId string, encryptedUrl string) error
	FindUrlByShortId(ctx context.Context, shortId string) (string, error)
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

func (r shortenRepository) FindUrlByShortId(ctx context.Context, shortId string) (string, error) {
	var shortUrl schemas.ShortedUrl

	err := r.db.
		WithContext(ctx).
		Where("short_id = ?", shortId).
		First(&shortUrl).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}

	return shortUrl.EncryptedUrl, nil
}
