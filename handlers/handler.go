package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohan7-code/url-shortener/config"
	"github.com/mohan7-code/url-shortener/dtos"
	service "github.com/mohan7-code/url-shortener/services"
	context "github.com/mohan7-code/url-shortener/utils/context"
)

func CreateShortURL(c *context.Context) {
	var req dtos.URLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	s := service.NewURLService()
	url, err := s.ShortenURL(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	shortURL := fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, url.ShortCode)

	c.JSON(http.StatusCreated, gin.H{
		"original_url": url.OriginalURL,
		"short_url":    shortURL,
	})
}

func RedirectURL(c *context.Context) {
	shortCode := c.Param("shortCode")

	s := service.NewURLService()
	url, err := s.GetOriginalURL(c, shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, url.OriginalURL)
}

func ListURLs(c *context.Context) {
	page, _ := strconv.Atoi(c.Query("page"))

	limit, _ := strconv.Atoi(c.Query("limit"))

	s := service.NewURLService()
	resp, err := s.ListURLs(c, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetAnalytics(ctx *context.Context) {

	code := ctx.Param("code")

	data, err := service.NewURLService().GetAnalytics(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
