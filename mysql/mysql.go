package mysql

import (
	gorm "simpleGorm"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	DSN       string
	DSNConfig *mysql.Config
}

type Dialector struct {
	*Config
}

func (d Dialector) Initialize(db *gorm.DB) error {
	return nil
}

func Open(dsn string) gorm.Dialector {
	cfg, _ := mysql.ParseDSN(dsn)
	return &Dialector{Config: &Config{DSN: dsn, DSNConfig: cfg}}
}
