package entity

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	var err error
	//Db, err = gorm.Open(mysql.Open("diamond@tcp(43.138.71.108:3306)/InkBook?charset=utf8mb4&parseTime=True&loc=Local&loc=Asia%2FShanghai"), &gorm.Config{})
	Db, err = gorm.Open(mysql.Open("root:Longyizhou2001@tcp(127.0.0.1:3306)/summer_project?charset=utf8mb4&parseTime=True&loc=Local&loc=Asia%2FShanghai"), &gorm.Config{})

	if err != nil {
		print(err.Error())
	}
}
