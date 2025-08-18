package gorm

import (
	"context"
)

type Statement struct {
	*DB
	Context context.Context
}

type Config struct {
	Dialector
	Conn ConnPool
}

type DB struct {
	*Config
	Error        error
	RowsAffected int64
	Statement    *Statement
}

func Open(dialector Dialector) (db *DB, err error) {
	config := &Config{}

	if config.Dialector == nil {
		config.Dialector = dialector
	}

	if config.Dialector != nil {
		if err := config.Dialector.Initialize(db); err != nil {
			return
		}
	}

	db.Statement = &Statement{
		DB:      db,
		Context: context.Background(),
	}
	return
}
