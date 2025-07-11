package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniciusabreusouza/url-shortener/internal/config/logger"
	"github.com/viniciusabreusouza/url-shortener/internal/service"
	"go.uber.org/zap"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortResponse struct {
	URL string `json:"url"`
}

type ShortenHandler struct {
	shortenService service.ShortenService
}

func NewShortenHandler(shortenService service.ShortenService) ShortenHandler {
	return ShortenHandler{
		shortenService: shortenService,
	}
}

func (s ShortenHandler) ShortUrl(c *gin.Context) {
	logger.Log.Info("Shortening URL")

	var shortenRequest ShortenRequest

	if err := c.ShouldBindJSON(&shortenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid request": err.Error()})
		return
	}

	if shortenRequest.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"invalid request": "URL is required"})
		return
	}

	logger.Log.Info("Shortening URL", zap.String("url", shortenRequest.URL))

	r, err := s.shortenService.ShortenUrl(c, shortenRequest.URL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ShortResponse{URL: r})
}
