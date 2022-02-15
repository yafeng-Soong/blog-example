package main

import (
	"encoding/gob"
	"gin_err_handler/errors"
	"gin_err_handler/middleware"
	"gin_err_handler/model"
	"gin_err_handler/service"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var userService service.UserService

func setPulicRouter(r *gin.Engine) {
	r.POST("/login", func(c *gin.Context) {
		var loginVo model.Login
		// 产生的一切错误都放到c.Error里
		// router里就可以不调用c.JSON()返回错误信息
		if e := c.ShouldBindJSON(&loginVo); e != nil {
			myErr := errors.VALID_ERROR
			myErr.Data = e.Error()
			c.Error(myErr)
			return
		}
		u, err := userService.Login(loginVo)
		if err != nil {
			c.Error(err)
			return
		}
		session := sessions.Default(c)
		session.Set("currentUser", u)
		if e := session.Save(); e != nil {
			// session保存出错也交给中间件处理，非自定义错误
			c.Error(e)
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功"})
	})
}

func setPrivateRouter(r *gin.Engine) {
	r.GET("/sayHello", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("currentUser").(model.User)
		c.String(http.StatusOK, "Hello "+user.Username)
	})
}

func main() {
	gob.Register(model.User{})
	r := gin.Default()
	r.Use(middleware.ErrorHandler()) // 错误处理中间件放最前面
	store := cookie.NewStore([]byte("yoursecret"))
	r.Use(sessions.Sessions("GSESSIONID", store))
	// 公共路由不需要cookie验证，所以放session中间件前注册
	setPulicRouter(r)
	r.Use(middleware.Cookie())
	setPrivateRouter(r)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
