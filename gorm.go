package gorm

import (
	"context"
	"simpleGorm/clause"
	"strings"
)

type Statement struct {
	*DB
	Context context.Context
	SQL     strings.Builder
	Clauses map[string]clause.Clause
}

func (stmt *Statement) AddClauseIfNotExists(v clause.Interface) {
	if c, ok := stmt.Clauses[v.Name()]; !ok || c.Expression == nil {
		stmt.AddClause(v)
	}
}

func (stmt *Statement) AddClause(v clause.Interface) {
	name := v.Name()
	c := stmt.Clauses[name]
	c.Name = name
	v.MergeClause(&c)
	stmt.Clauses[name] = c
}

type Config struct {
	Dialector
	Conn      ConnPool
	callbacks *callbacks
}

type DB struct {
	*Config
	Error        error
	RowsAffected int64
	Statement    *Statement
}

func (db *DB) Callback() *callbacks {
	return db.callbacks
}

func Open(dialector Dialector) (db *DB, err error) {
	config := &Config{}

	if config.Dialector == nil {
		config.Dialector = dialector
	}

	initializeCallbacks(db)

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
