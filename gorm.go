package gorm

import (
	"context"
	"fmt"
	"reflect"
	"simpleGorm/clause"
	"simpleGorm/schema"
	"strings"
	"time"
)

type Statement struct {
	*DB
	Table        string
	Context      context.Context
	SQL          strings.Builder
	Clauses      map[string]clause.Clause
	Schema       *schema.Schema
	ReflectValue reflect.Value
	Selects      []string
	BuildClauses []string
	Vars         []interface{}
	Dest         interface{}
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

func (stmt *Statement) SelectAndOmitColumns(requireCreate, requireUpdate bool) (map[string]bool, bool) {
	results := map[string]bool{}
	notRestricted := false

	processColumn := func(column string, result bool) {
		if column == "*" {
			notRestricted = result
			for _, dbName := range stmt.Schema.DBNames {
				results[dbName] = result
			}
		} else if field := stmt.Schema.LookUpField(column); field != nil && field.DBName != "" {
			results[field.DBName] = result
		} else {
			results[column] = result
		}
	}

	for _, column := range stmt.Selects {
		processColumn(column, true)
	}

	return results, !notRestricted && len(stmt.Selects) > 0
}

func (stmt *Statement) Build(clauses ...string) {
	var firstClauseWritten bool
	for _, name := range clauses {
		if c, ok := stmt.Clauses[name]; ok {
			if firstClauseWritten {
				stmt.WriteByte(' ')
			}
			firstClauseWritten = true
			c.Build(stmt)
		}
	}
}

func (stmt *Statement) WriteQuoted(field interface{}) {
	//TODO implement me
	panic("implement me")
}

func (stmt *Statement) AddVar(writer clause.Writer, i ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (stmt *Statement) AddError(err error) error {
	//TODO implement me
	panic("implement me")
}

func (stmt *Statement) WriteByte(c byte) error {
	return stmt.SQL.WriteByte(c)
}

func (stmt *Statement) WriteString(s string) (int, error) {
	return stmt.SQL.WriteString(s)
}

func (stmt *Statement) QuetoTo(writer clause.Writer, field interface{}) {
	//write := func(raw bool, str string) {
	//	if raw {
	//		writer.WriteString(str)
	//	} else {
	//		stmt.DB.Dialector.QuoteTo(writer, str)
	//	}
	//}
}

func (stmt *Statement) Parse(value interface{}) error {
	return stmt.ParseWithSpecialTableName(value, "")
}

func (stmt *Statement) ParseWithSpecialTableName(value interface{}, specialTableName string) (err error) {
	if stmt.Schema, err = schema.ParseWithSpecialTableName(value, stmt.DB.NamingStrategy, specialTableName); err == nil && stmt.Table == "" {
		if tables := strings.Split(stmt.Schema.Table, "."); len(tables) == 2 {
			stmt.Table = tables[1]
		}
		stmt.Table = stmt.Schema.Table
	}
	return err
}

type Config struct {
	Dialector
	Conn           ConnPool
	callbacks      *callbacks
	NowFunc        func() time.Time
	NamingStrategy schema.Namer
}

type DB struct {
	*Config
	Error          error
	RowsAffected   int64
	Statement      *Statement
	NamingStrategy schema.Namer
}

func (db *DB) Callback() *callbacks {
	return db.callbacks
}

func (db *DB) AddError(err error) error {
	if err != nil {
		if db.Error == nil {
			db.Error = err
		} else {
			db.Error = fmt.Errorf("%v; %w", db.Error, err)
		}
	}
	return db.Error
}

func Open(dialector Dialector) (db *DB, err error) {
	config := &Config{}

	if config.Dialector == nil {
		config.Dialector = dialector
	}

	if config.NowFunc == nil {
		config.NowFunc = func() time.Time { return time.Now().Local() }
	}

	if config.NamingStrategy == nil {
		config.NamingStrategy = schema.NamingStrategy{IdentifierMaxLength: 64}
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
