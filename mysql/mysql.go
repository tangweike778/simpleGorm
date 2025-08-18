package mysql

import (
	"database/sql"
	gorm "simpleGorm"

	"github.com/go-sql-driver/mysql"
)

const defaultDriver = "mysql"

type Config struct {
	DSN        string
	DSNConfig  *mysql.Config
	DriverName string
	Conn       gorm.ConnPool
}

type Dialector struct {
	*Config
}

func (d Dialector) Initialize(db *gorm.DB) error {
	if d.DriverName == "" {
		d.DriverName = defaultDriver
	}

	if d.Conn != nil {
		db.Conn = d.Conn
	} else {
		var err error
		db.Conn, err = sql.Open(d.DriverName, d.DSN)
		if err != nil {
			return err
		}
	}

	//

	return nil
}

func Open(dsn string) gorm.Dialector {
	cfg, _ := mysql.ParseDSN(dsn)
	return &Dialector{Config: &Config{DSN: dsn, DSNConfig: cfg}}
}
