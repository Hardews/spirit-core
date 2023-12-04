/**
 * @Author: Hardews
 * @Date: 2023/3/10 19:28
 * @Description:
**/

package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	db      *gorm.DB
	pwd     = os.Getenv("spirit_core_database_password")
	address = os.Getenv("spirit_core_database_address")
	dbName  = "spirit_core"
)

func InitDB() {
	dsn := "gmt_website:" + pwd + "@tcp(" + address + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	dB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database,err:", err)
	}

	db = dB

	err = db.AutoMigrate()
	if err != nil {
		log.Fatalln("failed to migrate table,err:", err)
	}
}
