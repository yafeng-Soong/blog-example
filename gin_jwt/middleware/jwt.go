package middleware

import (
	"gin_jwt/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			// 无token直接拒绝
			c.Abort()
			c.String(http.StatusOK, "未登录无权限")
			return
		}
		// 校验token
		claims, err := utils.ParseToken(auth)
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				// 若过期调用续签函数
				newToken, _ := utils.RenewToken(claims)
				if newToken != "" {
					// 续签成功給返回头设置一个newtoken字段
					c.Header("newtoken", newToken)
					c.Request.Header.Set("Authorization", newToken)
					c.Next()
					return
				}
			}
			// Token验证失败或续签失败直接拒绝请求
			c.Abort()
			c.String(http.StatusOK, err.Error())
			return
		}
		// token未过期继续执行1其他中间件
		c.Next()
	}
}
