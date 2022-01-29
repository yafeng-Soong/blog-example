package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Cookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		// log.Println("进入了Session中间件")
		session := sessions.Default(c)
		// log.Println(session)
		if session.Get("currentUser") == nil {
			c.Abort()
			c.String(http.StatusOK, "未登录无权限")
		} else {
			c.Next()
		}
	}
}
