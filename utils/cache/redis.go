package cache

import (
	"fmt"
	"sync"

	"github.com/mohan7-code/url-shortener/config"
	"github.com/redis/go-redis/v9"
)

var once sync.Once
var cl *redis.Client

func SetRedis() error {
	var err error
	once.Do(func() {
		option, err := redis.ParseURL(config.AppConfig.RedisURL)
		if err != nil {
			fmt.Println("Unable to parse Redis URL: %w", err)
			return
		}

		cl = redis.NewClient(option)
		fmt.Println("cl is", cl)
	})
	return err
}

type Pool struct {
	*redis.Client
}

func New() *Pool {
	return &Pool{
		Client: cl,
	}
}
