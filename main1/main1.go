package main

import (
	"log"
	gorm "simpleGorm"
	"simpleGorm/mysql"
)

type User struct {
	ID   int64  `gorm:"primary_key;auto_increment;unique"`
	Name string `gorm:"type:varchar(255);not null;unique"`
	Age  int64  `gorm:"type:int"`
}

func main() {
	db, err := gorm.Open(mysql.Open("root:twk123456@tcp(127.0.0.1:3306)/smallProgram?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		log.Fatal(err)
	}
	user := User{
		Name: "twk",
		Age:  25,
	}
	err = db.Create(&user).Error
	if err != nil {
		log.Fatal(err)
	}
}
