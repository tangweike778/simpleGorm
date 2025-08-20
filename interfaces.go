package gorm

import (
	"context"
	"database/sql"
)

type Dialector interface {
	Initialize(*DB) error
}

// ConnPool 连接池
type ConnPool interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
