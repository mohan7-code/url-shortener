package routes

import (
	"github.com/gin-gonic/gin"
	handler "github.com/mohan7-code/url-shortener/handlers"
	mw "github.com/mohan7-code/url-shortener/middleware"
)

func UrlRoutes(router *gin.RouterGroup) {

	router.POST("/shorten", mw.MiddleWare(handler.CreateShortURL))
	router.GET("/:shortCode", mw.MiddleWare(handler.RedirectURL))
	router.GET("/urls", mw.MiddleWare(handler.ListURLs))
}
