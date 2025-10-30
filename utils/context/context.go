package context

import (
	"github.com/gin-gonic/gin"
	"github.com/mohan7-code/url-shortener/database"
	"go.uber.org/zap"
)

type Context struct {
	DB  *database.DBConn
	Log *zap.Logger
	*gin.Context
}

func (a *Context) Copy() *Context {

	return &Context{
		DB:      a.DB,
		Log:     a.Log,
		Context: a.Context.Copy(),
	}
}
