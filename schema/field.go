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
	StructField           reflect.StructField
	ReflectValueOf        func(context.Context, reflect.Value) reflect.Value
}

func (f *Field) setupValuerAndSetter() {
	fieldIndex := f.StructField.Index[0]
	switch {
	case len(f.StructField.Index) == 1 && fieldIndex >= 0:
		f.ValueOf = func(ctx context.Context, value reflect.Value) (interface{}, bool) {
			fieldValue := reflect.Indirect(value).FieldByName(f.Name)
			return fieldValue.Interface(), fieldValue.IsZero()
		}
	}

	switch {
	case len(f.StructField.Index) == 1 && fieldIndex >= 0:
		f.ReflectValueOf = func(ctx context.Context, value reflect.Value) reflect.Value {
			return reflect.Indirect(value).Field(fieldIndex)
		}
	}

	switch f.FieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f.Set = func(ctx context.Context, value reflect.Value, v interface{}) (err error) {
			switch data := v.(type) {
			case int64:
				f.ReflectValueOf(ctx, value).SetInt(data)
			case int:
				f.ReflectValueOf(ctx, value).SetInt(int64(data))
			}
			return err
		}
	case reflect.String:
		f.Set = func(ctx context.Context, value reflect.Value, v interface{}) (err error) {
			switch data := v.(type) {
			case string:
				f.ReflectValueOf(ctx, value).SetString(data)
			}
			return err
		}
	}
}
