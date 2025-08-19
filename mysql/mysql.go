package mysql

import (
	"database/sql"
	gorm "simpleGorm"
	"simpleGorm/callbacks"

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

var (
	// CreateClauses create clauses
	CreateClauses = []string{"INSERT", "VALUES", "ON CONFLICT"}
	// QueryClauses query clauses
	QueryClauses = []string{}
	// UpdateClauses update clauses
	UpdateClauses = []string{"UPDATE", "SET", "WHERE", "ORDER BY", "LIMIT"}
	// DeleteClauses delete clauses
	DeleteClauses = []string{"DELETE", "FROM", "WHERE", "ORDER BY", "LIMIT"}
)

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

	// register callbacks
	callbackConfig := &callbacks.Config{
		CreateClauses: CreateClauses,
		QueryClauses:  QueryClauses,
		UpdateClauses: UpdateClauses,
		DeleteClauses: DeleteClauses,
	}

	callbacks.RegisterDefaultCallbacks(db, callbackConfig)

	return nil
}

func Open(dsn string) gorm.Dialector {
	cfg, _ := mysql.ParseDSN(dsn)
	return &Dialector{Config: &Config{DSN: dsn, DSNConfig: cfg}}
}
