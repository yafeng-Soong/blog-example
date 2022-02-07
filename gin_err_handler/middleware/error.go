package middleware

import (
	"gin_err_handler/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先调用c.Next()执行后面的中间件
		// 所有中间件及router处理完毕后从这里开始执行
		// 检查c.Errors中是否有错误
		if length := len(c.Errors); length > 0 {
			e := c.Errors[length-1]
			err := e.Err
			if err != nil {
				// 若是自定义的错误则将code、msg返回
				if myErr, ok := err.(*errors.MyError); ok {
					c.JSON(http.StatusOK, gin.H{
						"code": myErr.Code,
						"msg":  myErr.Msg,
					})
				} else {
					// 若非自定义错误则返回详细错误信息err.Error()
					// 比如保存session出错时设置的err
					c.JSON(http.StatusOK, gin.H{
						"code": 500,
						"msg":  "服务器异常",
						"data": err.Error(),
					})
				}
			}
		}
	}
}
