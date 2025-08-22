package gorm

import (
	"context"
	"database/sql"
	"simpleGorm/clause"
)

type Dialector interface {
	Initialize(*DB) error
	QuoteTo(writer clause.Writer, str string)
	BindVarTo(writer clause.Writer, stmt *Statement, v interface{})
}

// ConnPool 连接池
type ConnPool interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
