package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	initDBOnce sync.Once
	db         *gorm.DB
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type User struct {
	gorm.Model
	Name     string
	Email    *string
	Age      uint8
	Birthday time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("创建前...")
	return
}

func main() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	db.Create(&Product{Code: "D42", Price: 100})

	var product Product
	db.First(&product, 1)
	db.First(&product, "code = ?", "D42")
	fmt.Println(product)

	db.Model(&product).Update("Price", 200)
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	fmt.Println(product)

	db.Delete(&product, 1)
}

func initDB() {
	d, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	db = d
	if err != nil {
		panic("failed to connect database")
	}
}
func GetDB() *gorm.DB {
	initDBOnce.Do(initDB)
	return db
}
