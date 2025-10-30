package middleware

import (
	"net/http"
	"sync"

	"github.com/mohan7-code/url-shortener/database"
	context "github.com/mohan7-code/url-shortener/utils/context"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	logger   *zap.Logger
	mu       sync.Mutex
	limiters = make(map[string]*rate.Limiter)
)

func init() {
	var err error
	logConfig := zap.NewProductionConfig()

	logger, err = logConfig.Build()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}

}

func MiddleWare(next func(*context.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		//simple Rate limiting using Token Bucket method
		mu.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(1, 5)
			limiters[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests, Try after sometime"})
			c.Abort()
			return
		}

		appCtx := &context.Context{
			DB:      database.New(),
			Log:     logger,
			Context: c,
		}

		next(appCtx)
	}
}
