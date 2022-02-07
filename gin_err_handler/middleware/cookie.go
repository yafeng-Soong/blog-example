package middleware

import (
	"gin_err_handler/errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Cookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get("currentUser") == nil {
			c.Abort()
			c.Error(errors.UNAUTHORIZED)
			// 缩写成c.AbortWithError(errors.UNAUTHORIZED)会有问题
			// 导致response的Content-type不是JSON
			return
		}
		c.Next()
	}
}
