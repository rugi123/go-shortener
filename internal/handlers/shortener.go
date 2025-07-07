package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rugi123/go-shortener/internal/domain/service"
)

type ShortenerHandler struct {
	service    *service.ShortenerService
	baseDomain string
}

func NewShortenHandler(service *service.ShortenerService, baseDomain string) *ShortenerHandler {
	return &ShortenerHandler{
		service:    service,
		baseDomain: baseDomain,
	}
}

func (h *ShortenerHandler) Shorten(c *gin.Context) {
	var request struct {
		URL string `json:"url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	shortKey, err := h.service.ShortenURL(c.Request.Context(), request.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Could not shorten URL: %s", err),
		})
		return
	}

	shortURL := h.baseDomain + "/" + shortKey

	c.JSON(http.StatusOK, gin.H{
		"original_url": request.URL,
		"short_url":    shortURL,
		"short_key":    shortKey,
	})
}

func (h *ShortenerHandler) Redirect(c *gin.Context) {
	shortKey := c.Param("key")
	if shortKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing short key"})
		return
	}
	originalURL, err := h.service.ExpandURL(c.Request.Context(), shortKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "URL not found",
			"key":   shortKey,
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}
