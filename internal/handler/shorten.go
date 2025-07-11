package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/viniciusabreusouza/url-shortener/internal/config/logger"
	"github.com/viniciusabreusouza/url-shortener/internal/repository"
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

	if !strings.HasPrefix(shortenRequest.URL, "http://") || !strings.HasPrefix(shortenRequest.URL, "https://") {
		c.JSON(http.StatusBadRequest, gin.H{"invalid request": "URL must be a valid URL containing http:// or https://"})
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

func (s ShortenHandler) RedirectUrl(c *gin.Context) {
	logger.Log.Info("Redirecting URL")

	shortId := c.Param("shortId")

	logger.Log.Info("Redirecting URL", zap.String("shortId", shortId))

	r, err := s.shortenService.FindUrlByShortId(c, shortId)

	logger.Log.Info("URL found", zap.String("url", r))

	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if r == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, r)
}
