package schema

import (
	"errors"
	"sync"
)

type Schema struct {
	Table                    string
	DBNames                  []string
	FieldsByDBName           map[string]*Field
	FieldsWithDefaultDBValue map[string]*Field // 有默认值的字段
	PrioritizedPrimaryField  *Field
}

// ErrUnsupportedDataType unsupported data type
var ErrUnsupportedDataType = errors.New("unsupported data type")

func (s Schema) LookUpField(column string) *Field {
	if field, ok := s.FieldsByDBName[column]; ok {
		return field
	}
	return nil
}

func ParseWithSpecialTableName(dest interface{}, cacheStore *sync.Map, namer Namer, specialTableName string) (*Schema, error) {

}
