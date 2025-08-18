package gorm

type Dialector interface {
	Initialize(*DB) error
}

// ConnPool 连接池
type ConnPool interface {
}
