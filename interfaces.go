package gorm

type Dialector interface {
	Initialize(*DB) error
}
