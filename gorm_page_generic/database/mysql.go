package database

import (
	// 导入gorm工具包

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
