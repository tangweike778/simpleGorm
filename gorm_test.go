package gorm

import (
	"simpleGorm/mysql"
	"testing"
)

func TestOpen(t *testing.T) {
	mysqlDb := mysql.Open("root:twk123456@tcp(127.0.0.1:3306)/smallProgram?charset=utf8mb4&parseTime=True&loc=Local")
	if _, err := Open(mysqlDb); err != nil {
		t.Error(err)
	}

}
