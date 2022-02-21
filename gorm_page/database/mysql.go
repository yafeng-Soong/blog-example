package database

import (
	// 导入gorm工具包

	"log"
	"reflect"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Close() {
	db, _ := DB.DB()
	db.Close()
}

func InitMysql() error {
	dsn := "root:fuck@you@tcp(127.0.0.1:3306)/world?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return err
	}
	DB = db
	return nil
}

func SelectPage(page *Page, wrapper map[string]interface{}, model interface{}) (e error) {
	e = nil
	tag := false
	if page.PageSize <= 0 {
		page.PageSize = 0
		tag = true
	}
	if page.CurrentPage <= 0 {
		page.CurrentPage = 0
		tag = true
	}
	if tag {
		page.Data = []interface{}{}
		return
	}
	t := reflect.TypeOf(model)
	list := reflect.Zero(reflect.SliceOf(t)).Interface()
	DB.Model(&model).Where(wrapper).Count(&page.Total)
	if page.Total == 0 {
		page.Data = []interface{}{}
	}
	page.Pages = page.Total / page.PageSize
	if page.Total%page.PageSize != 0 {
		page.Pages++
	}
	size := page.PageSize
	offset := int((page.CurrentPage - 1) * size)
	log.Println(size, offset)
	e = DB.Model(&model).Where(wrapper).Limit(int(size)).Offset(offset).Find(&list).Error
	if e != nil {
		return
	}
	// log.Println(list)
	page.Data = list
	return
}
