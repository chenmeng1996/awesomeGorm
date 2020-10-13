package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMysql() {
	dsn := "user:root@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
