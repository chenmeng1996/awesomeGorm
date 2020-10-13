package main

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"testing"
	"time"
)

var testDB = GetDB()

func TestMigrate(t *testing.T) {
	testDB.AutoMigrate(&User{})
}

func TestCreate(t *testing.T) {
	user := User{
		Name:     "cm",
		Age:      18,
		Birthday: time.Now(),
	}
	result := testDB.Create(&user)

	fmt.Println(user.ID, result.Error, result.RowsAffected)
}

func TestCreateOmit(t *testing.T) {
	user := User{
		Name:     "cm",
		Age:      18,
		Birthday: time.Now(),
	}
	testDB.Omit("Age").Create(&user)
}

func TestBatchCreate(t *testing.T) {
	users := []User{{Name: "cm1"}, {Name: "cm2"}}
	testDB.Create(&users)

	for _, user := range users {
		fmt.Println(user.ID)
	}
}

func TestCreateFromMap(t *testing.T) {
	testDB.Model(&User{}).Create(map[string]interface{}{
		"Name": "cm", "Age": 18,
	})
}

func TestBatchCreateFromMap(t *testing.T) {
	testDB.Model(&User{}).Create([]map[string]interface{}{
		{"Name": "cm", "Age": 18},
		{"Name": "cm1", "Age": 18},
	})
}

func TestFirst(t *testing.T) {
	user := User{}
	// SELECT * FROM users ORDER BY id LIMIT 1;
	testDB.First(&user)
	fmt.Println(user)

	user1 := User{}
	// SELECT * FROM users WHERE id = 2 ORDER BY id LIMIT 1;
	testDB.First(&user1, 2) // 如果传入的结构体不为空，会报错
	fmt.Println(user1)

	user2 := User{}
	// SELECT * FROM users WHERE name = 'cm1' ORDER BY id LIMIT 1;
	testDB.First(&user2, "name = ?", "cm1")
	fmt.Println(user2)

	user3 := User{}
	// SELECT * FROM users WHERE name = 'cm1' ORDER BY id LIMIT 1;
	testDB.First(&user3, User{Name: "cm1"})
	fmt.Println(user3)

	user4 := User{}
	// SELECT * FROM users WHERE name = 'cm1' ORDER BY id LIMIT 1;
	testDB.First(&user4, map[string]interface{}{"Name": "cm1"})
	fmt.Println(user4)
}

func TestTake(t *testing.T) {
	user := User{}
	testDB.Take(&user) // SELECT * FROM users LIMIT 1;
	fmt.Println(user)
}

func TestLast(t *testing.T) {
	user := User{}
	testDB.Last(&user) // SELECT * FROM users ORDER BY id DESC LIMIT 1;
	fmt.Println(user)
}

func TestIn(t *testing.T) {
	var users []User
	testDB.Find(&users, []int{1, 2}) // SELECT * FROM users WHERE id IN (1,2)
	fmt.Println(users)
}

func TestRecordNotFound(t *testing.T) {
	user := User{}
	res := testDB.Where("id = ?", 999).Take(&user) // SELECT * FROM users WHERE id = 999 LIMIT 1;
	fmt.Println(res.RowsAffected)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("没找到记录")
	}
}

func TestNotEqual(t *testing.T) {
	var users []User
	testDB.Where("name <> ?", "cm").Find(&users) // SELECT * FROM users WHERE name <> 'cm'
	fmt.Println(users)
}

func TestLike(t *testing.T) {
	var users []User
	testDB.Where("name like ?", "%cm%").Find(&users) // SELECT * FROM users WHERE name like '%cm%'
	fmt.Println(users)
}

func TestAnd(t *testing.T) {
	var users []User
	testDB.Where("name = ? AND age >= ?", "cm", 18).Find(&users)
	fmt.Println(users)
}

func TestTime(t *testing.T) {
	var users []User
	//updateTime, _ := time.Parse("2006-01-02 15:04:05", "2020-10-13 16:00:00")
	testDB.Where("updated_at > ?", "2020-10-13 16:00:00").Find(&users) // SELECT * FROM users WHERE name <> 'cm'
	fmt.Println(users)
}

func TestBetween(t *testing.T) {
	var users []User
	testDB.Where("updated_at BETWEEN ? AND ?", "2020-10-13 16:00:00", "2020-10-13 16:01:00").Find(&users)
	fmt.Println(users)
}

func TestConditionFromStructOrMap(t *testing.T) {
	var users []User
	testDB.Where(&User{Name: "cm", Age: 18}).Find(&users)
	fmt.Println(users)

	var users1 []User
	testDB.Where(map[string]interface{}{"name": "cm", "age": 18}).Find(&users1)
	fmt.Println(users1)
}

func TestNot(t *testing.T) {
	user := User{}
	// SELECT * FROM users WHERE NOT name = "cm" ORDER BY id LIMIT 1;
	testDB.Not("name = ?", "cm").First(&user)
	fmt.Println(user)

	user1 := User{}
	// SELECT * FROM users WHERE name <> "cm" AND age <> 18 ORDER BY id LIMIT 1;
	testDB.Not(User{Name: "cm", Age: 18}).First(&user1)
	fmt.Println(user1)
}

func TestNotIn(t *testing.T) {
	var users []User
	// SELECT * FROM users WHERE name NOT IN ('cm', 'cm2');
	testDB.Not(map[string]interface{}{"name": []string{"cm", "cm2"}}).Find(&users)
	fmt.Println(users)
}

func TestOr(t *testing.T) {
	var users []User
	// SELECT * FROM users WHERE name = 'cm' OR (age = 18)
	testDB.Where("name = ?", "cm").Or("age = ?", 18).Find(&users)
	fmt.Println(users)
}

func TestGroupWhere(t *testing.T) {
	var users []User
	// SELECT * FROM users WHERE (name = 'cm' AND (age = 18 OR age = 19)) OR (name = 'cm1' AND age = 19)
	testDB.Debug().Where(
		testDB.Where("name = ?", "cm").Where(
			testDB.Where("age = ?", 18).Or("age = ?", 19)),
	).Or(
		testDB.Where("name = ?", "cm1").Where("age = ?", 18),
	).Find(&users)
	fmt.Println(users)
}

func TestSelect(t *testing.T) {
	var users []User
	testDB.Select("name", "age").Find(&users)
	fmt.Println(111, users)

	testDB.Select("name", "age").Find(&users)
	fmt.Println(222, users)

	rows, _ := testDB.Model(&User{}).Select("count(age)").Rows()
	for rows.Next() {
		rows.Scan()
	}
}

func TestOrder(t *testing.T) {
	var users []User
	testDB.Order("age desc, name").Find(&users)
	fmt.Println(users)
}

func TestOffsetLimit(t *testing.T) {
	var users []User
	// SELECT * FROM users limit 3, 2
	testDB.Offset(3).Limit(2).Find(&users)
	fmt.Println(users)
}

func TestGroupBy(t *testing.T) {
	// SELECT name, SUM(age) FROM users GROUP BY name
	var results []struct {
		Name  string
		Total int
	}
	testDB.Model(&User{}).Select("name, sum(age) as total").Group("name").Find(&results)
	fmt.Println(results)
}

func TestHaving(t *testing.T) {
	// SELECT name, SUM(age) FROM users GROUP BY name HAVING name = 'cm'
	var results []struct {
		Name  string
		Total int
	}
	testDB.Model(&User{}).Select("name, sum(age) as total").Group("name").Having("name = ?", "cm").Find(&results)
	fmt.Println(results)
}

func TestDistinct(t *testing.T) {
	var users []User
	testDB.Distinct("name", "age").Find(&users)
}
