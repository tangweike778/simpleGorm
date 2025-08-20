package schema

import (
	"context"
	"reflect"
)

type Field struct {
	Name                  string
	DBName                string
	ValueOf               func(context.Context, reflect.Value) (value interface{}, zero bool)
	DefaultValueInterface interface{}
	HasDefaultValue       bool
	Set                   func(context.Context, reflect.Value, interface{}) error
}
