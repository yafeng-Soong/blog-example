package main

import (
	"gorm_page/database"
	"gorm_page/model"
	"gorm_page/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	cityService    service.CityService
	countryService service.CountryService
)

func main() {
	if err := database.InitMysql(); err != nil {
		log.Fatalln("数据库连接出错")
	}
	defer database.Close()
	r := gin.Default()
	r.POST("/city/list", func(c *gin.Context) {
		var queryVo model.CityQueryInfo
		if e := c.ShouldBindJSON(&queryVo); e != nil {
			c.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
			return
		}
		p := &database.Page{}
		if e := cityService.SelectPageList(p, queryVo); e != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "操作失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": p})
	})

	r.POST("/country/list", func(c *gin.Context) {
		var queryVo model.CountryQueryInfo
		if e := c.ShouldBindJSON(&queryVo); e != nil {
			c.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
			return
		}
		p := &database.Page{}
		if e := countryService.SelectPageList(p, queryVo); e != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "操作失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功", "data": p})
	})
	r.Run(":8080")
}
