package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniciusabreusouza/url-shortener/internal/config"
	"github.com/viniciusabreusouza/url-shortener/internal/handler"
	"github.com/viniciusabreusouza/url-shortener/internal/repository"
	"github.com/viniciusabreusouza/url-shortener/internal/service"
)

func initializeRoutes(router *gin.Engine) {
	r := router.Group("/api/v1")

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
	var shortenRepo repository.ShortenRepository = repository.NewShortenRepository(config.GetDB())
	var shortenService service.ShortenService = service.NewShortenService(shortenRepo)

	shortenHandler := handler.NewShortenHandler(shortenService)

	r.POST("/shorten", shortenHandler.ShortUrl)
	r.GET("/:shortId", shortenHandler.RedirectUrl)

}
