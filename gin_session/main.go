package main

import (
	"encoding/gob"
	"gin_session/middleware"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var db = &User{Id: 10001, Email: "abc@gmail.cn", Username: "Alice", Password: "123456"}

func getCurrentUser(c *gin.Context) (userInfo User) {
	session := sessions.Default(c)
	userInfo = session.Get("currentUser").(User)
	return
}

func setCurrentUser(c *gin.Context, userInfo User) {
	session := sessions.Default(c)
	session.Set("currentUser", userInfo)
	session.Save()
}

func setupRouter(r *gin.Engine) {
	r.POST("/login", func(c *gin.Context) {
		var loginVo User
		if c.ShouldBindJSON(&loginVo) != nil {
			c.String(http.StatusOK, "参数错误")
			return
		}
		if loginVo.Email == db.Email && loginVo.Password == db.Password {
			setCurrentUser(c, *db)
			c.String(http.StatusOK, "登录成功")
		} else {
			c.String(http.StatusOK, "登录失败")
		}
	})

	r.GET("/sayHello", middleware.Cookie(), func(c *gin.Context) {
		userInfo := getCurrentUser(c)
		c.String(http.StatusOK, "Hello "+userInfo.Username)
	})
}

func main() {
	gob.Register(User{})
	r := gin.Default()
	store := cookie.NewStore([]byte("snaosnca"))
	r.Use(sessions.Sessions("SESSIONID", store))
	setupRouter(r)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
