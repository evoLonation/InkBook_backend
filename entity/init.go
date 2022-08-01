package entity

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open(mysql.Open("diamond@tcp(43.138.71.108:3306)/InkBook?charset=utf8mb4&parseTime=True&loc=Local&serverTimezone=Asia/Shanghai"), &gorm.Config{})
	if err != nil {
		print(err.Error())
	}
}
