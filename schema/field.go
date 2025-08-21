package schema

import (
	"context"
	"reflect"
)

type (
	DataType string
)

const (
	Bool   DataType = "bool"
	Int    DataType = "int"
	Uint   DataType = "uint"
	Float  DataType = "float"
	String DataType = "string"
	Time   DataType = "time"
	Bytes  DataType = "bytes"
)

type Field struct {
	Name                  string
	DBName                string
	ValueOf               func(context.Context, reflect.Value) (value interface{}, zero bool)
	DefaultValueInterface interface{}
	HasDefaultValue       bool
	Set                   func(context.Context, reflect.Value, interface{}) error
	FieldType             reflect.Type
	IndirectFieldType     reflect.Type
	DataType              DataType
	PrimaryKey            bool
}
