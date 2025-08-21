package schema

import (
	"errors"
	"fmt"
	"go/ast"
	"reflect"
	"strings"
	"sync"

	"gorm.io/gorm/utils"
)

type Schema struct {
	Table                    string
	DBNames                  []string
	FieldsByDBName           map[string]*Field
	FieldsWithDefaultDBValue map[string]*Field // 有默认值的字段
	PrioritizedPrimaryField  *Field
	Name                     string
	ModelType                reflect.Type
	FieldsByName             map[string]*Field
	Namer                    Namer
	Fields                   []*Field
	PrimaryFields            []*Field
	PrimaryFieldDBNames      []string
}

// ErrUnsupportedDataType unsupported data type
var ErrUnsupportedDataType = errors.New("unsupported data type")

func (s Schema) LookUpField(column string) *Field {
	if field, ok := s.FieldsByDBName[column]; ok {
		return field
	}
	return nil
}

func (s Schema) ParseField(fieldStruct reflect.StructField) *Field {
	var (
		tagSetting = ParseTagSetting(fieldStruct.Tag.Get("gorm"), ";")
	)

	field := &Field{
		Name:              fieldStruct.Name,
		DBName:            tagSetting["COLUMN"],
		HasDefaultValue:   false,
		PrimaryKey:        utils.CheckTruth(tagSetting["PRIMARYKEY"], tagSetting["PRIMARY_KEY"]),
		FieldType:         fieldStruct.Type,
		IndirectFieldType: fieldStruct.Type,
		StructField:       fieldStruct,
	}

	fieldValue := reflect.New(field.IndirectFieldType)
	switch reflect.Indirect(fieldValue).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.DataType = Int
	case reflect.String:
		field.DataType = String
	default:
	}
	return field
}

func ParseTagSetting(str string, sep string) map[string]string {
	settings := map[string]string{}
	names := strings.Split(str, sep)

	for i := 0; i < len(names); i++ {
		j := i
		if len(names[j]) > 0 {
			for {
				if names[j][len(names[j])-1] == '\\' {
					i++
					names[j] = names[j][0:len(names)-1] + sep + names[i]
					names[i] = ""
				} else {
					break
				}
			}
		}

		values := strings.Split(names[j], ":")
		k := strings.TrimSpace(strings.ToUpper(values[0]))
		if len(values) > 2 {
			settings[k] = strings.Join(values[1:], ":")
		} else if k != "" {
			settings[k] = k
		}
	}
	return settings
}

type Namer interface {
	TableName(table string) string
	SchemaName(table string) string
	ColumnName(table, column string) string
}

func ParseWithSpecialTableName(dest interface{}, cacheStore *sync.Map, namer Namer, specialTableName string) (*Schema, error) {
	if dest == nil {
		return nil, fmt.Errorf("%w: %+v", ErrUnsupportedDataType, dest)
	}

	value := reflect.ValueOf(dest)
	modelType := reflect.Indirect(value).Type()
	tableName := namer.TableName(modelType.Name())
	if specialTableName != "" && specialTableName != tableName {
		tableName = specialTableName
	}
	schema := &Schema{
		Name:           modelType.Name(),
		ModelType:      modelType,
		Table:          tableName,
		FieldsByName:   map[string]*Field{},
		FieldsByDBName: map[string]*Field{},
		Namer:          namer,
	}

	for i := 0; i < modelType.NumField(); i++ {
		if fieldStruct := modelType.Field(i); ast.IsExported(fieldStruct.Name) {
			field := schema.ParseField(fieldStruct)
			schema.Fields = append(schema.Fields, field)
		}
	}

	for _, field := range schema.Fields {
		if field.DBName == "" && field.DataType != "" {
			field.DBName = namer.ColumnName(schema.Table, field.Name)
		}

		if field.DBName != "" {
			if v, ok := schema.FieldsByDBName[field.DBName]; !ok {
				schema.DBNames = append(schema.DBNames, field.DBName)
				schema.FieldsByName[field.Name] = field
				schema.FieldsByDBName[field.DBName] = field

				if v != nil && v.PrimaryKey {
					schema.PrimaryFields = append(schema.PrimaryFields, field)
				}
			}
		}

		field.setupValuerAndSetter()
	}

	prioritizedPrimaryField := schema.LookUpField("id")
	if prioritizedPrimaryField == nil {
		prioritizedPrimaryField = schema.LookUpField("ID")
	}

	if prioritizedPrimaryField != nil {
		if prioritizedPrimaryField.PrimaryKey {
			schema.PrioritizedPrimaryField = prioritizedPrimaryField
		} else if len(schema.PrimaryFields) == 0 {
			prioritizedPrimaryField.PrimaryKey = true
			schema.PrioritizedPrimaryField = prioritizedPrimaryField
			schema.PrimaryFields = append(schema.PrimaryFields, prioritizedPrimaryField)
		}
	}

	if schema.PrioritizedPrimaryField == nil {
		if len(schema.PrimaryFields) == 1 {
			schema.PrioritizedPrimaryField = schema.PrimaryFields[0]
		}
	}

	for _, field := range schema.PrimaryFields {
		schema.PrimaryFieldDBNames = append(schema.PrimaryFieldDBNames, field.DBName)
	}

	return schema, nil
}
