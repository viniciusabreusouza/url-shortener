package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"os"

	"github.com/viniciusabreusouza/url-shortener/internal/config/logger"
	"github.com/viniciusabreusouza/url-shortener/internal/repository"
	"go.uber.org/zap"
)

var lettersRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type ShortenService interface {
	ShortenUrl(ctx context.Context, url string) (string, error)
	FindUrlByShortId(ctx context.Context, shortId string) (string, error)
}

type shortenService struct {
	repository repository.ShortenRepository
}

func NewShortenService(repository repository.ShortenRepository) ShortenService {
	return &shortenService{
		repository: repository,
	}
}

func (s shortenService) ShortenUrl(ctx context.Context, url string) (string, error) {
	encryptedUrl, err := s.encrypt(url)

	if err != nil {
		return "", err
	}

	shortId := s.generateShortId()

	if err := s.repository.ShortenUrl(ctx, shortId, encryptedUrl); err != nil {
		return "", err
	}

	shortUrl := "http://localhost:8080/api/v1/" + shortId

	return shortUrl, err
}

func (s shortenService) FindUrlByShortId(ctx context.Context, shortId string) (string, error) {
	encryptedUrl, err := s.repository.FindUrlByShortId(ctx, shortId)

	if err != nil {
		return "", err
	}

	if encryptedUrl == "" {
		return "", nil
	}

	return s.decrypt(encryptedUrl)
}

func (s shortenService) encrypt(url string) (shortUrl string, err error) {
	secret := os.Getenv("ENCRYPTION_KEY")

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		logger.Log.Error("Failed to create cipher block", zap.Error(err))
		return "", err
	}

	plainText := []byte(url)
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]

	if _, err = rand.Read(iv); err != nil {
		logger.Log.Error("Failed to generate random iv", zap.Error(err))
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return hex.EncodeToString(cipherText), nil
}

func (s shortenService) generateShortId() string {
	b := make([]rune, 6)

	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(lettersRunes))))
		if err != nil {
			logger.Log.Error("Failed to generate random number", zap.Error(err))
			return ""
		}

		b[i] = lettersRunes[num.Int64()]
	}

	return string(b)
}

func (s shortenService) decrypt(encryptedUrl string) (string, error) {
	secret := os.Getenv("ENCRYPTION_KEY")

	block, err := aes.NewCipher([]byte(secret))

	if err != nil {
		logger.Log.Error("Failed to create cipher block", zap.Error(err))
		return "", err
	}

	cipherText, err := hex.DecodeString(encryptedUrl)

	if err != nil {
		logger.Log.Error("Failed to decode hex string", zap.Error(err))
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
