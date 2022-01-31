package main

import (
	"gin_jwt/middleware"
	"gin_jwt/model"
	"gin_jwt/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = &model.User{Id: 10001, Email: "abc@gmail.xyz", UserName: "Alice", Password: "123456"}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("login", func(c *gin.Context) {
		var userVo model.User
		if c.ShouldBindJSON(&userVo) != nil {
			c.String(http.StatusOK, "参数错误")
			return
		}
		if userVo.Email == db.Email && userVo.Password == db.Password {
			info := model.NewInfo(*db)
			tokenString, _ := utils.GenerateToken(*info)
			c.JSON(http.StatusOK, gin.H{
				"code":  201,
				"token": tokenString,
				"msg":   "登录成功",
			})
			return
		}
		c.String(http.StatusOK, "登录失败")
		return
	})

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	authorized := r.Group("/", middleware.JWTAuth())

	authorized.GET("/sayHello", func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		claims, _ := utils.ParseToken(auth)
		log.Println(claims)
		c.String(http.StatusOK, "hello "+claims.User.UserName)
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
