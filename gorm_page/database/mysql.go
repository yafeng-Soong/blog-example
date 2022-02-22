package database

import (
	// 导入gorm工具包

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
	dsn := "root:yourPassword@tcp(127.0.0.1:3306)/world?charset=utf8mb4&parseTime=True&loc=Local"
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
	DB.Model(&model).Where(wrapper).Count(&page.Total)
	if page.Total == 0 {
		page.Data = []interface{}{}
		return
	}
	// 反射获得类型
	t := reflect.TypeOf(model)
	// 再通过反射创建创建对应类型的数组
	list := reflect.Zero(reflect.SliceOf(t)).Interface()
	e = DB.Model(&model).Where(wrapper).Scopes(Paginate(page)).Find(&list).Error
	if e != nil {
		return
	}
	page.Data = list
	return
}

func Paginate(page *Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.CurrentPage <= 0 {
			page.CurrentPage = 0
		}
		switch {
		case page.PageSize > 100:
			page.PageSize = 100
		case page.PageSize <= 0:
			page.PageSize = 10
		}
		page.Pages = page.Total / page.PageSize
		if page.Total%page.PageSize != 0 {
			page.Pages++
		}
		p := page.CurrentPage
		if page.CurrentPage > page.Pages {
			p = page.Pages
		}
		size := page.PageSize
		offset := int((p - 1) * size)
		return db.Offset(offset).Limit(int(size))
	}
}
