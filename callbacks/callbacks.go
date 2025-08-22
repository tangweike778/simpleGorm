package callbacks

import (
	"log"
	"reflect"
	gorm "simpleGorm"
	"simpleGorm/clause"
	"simpleGorm/schema"
)

var (
	createClauses = []string{"INSERT", "VALUES", "ON CONFLICT"}
	queryClauses  = []string{"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT", "FOR"}
	updateClauses = []string{"UPDATE", "SET", "WHERE"}
	deleteClauses = []string{"DELETE", "FROM", "WHERE"}
)

func RegisterDefaultCallbacks(db *gorm.DB, config *Config) {
	if len(config.CreateClauses) == 0 {
		config.CreateClauses = createClauses
	}
	if len(config.QueryClauses) == 0 {
		config.QueryClauses = queryClauses
	}
	if len(config.DeleteClauses) == 0 {
		config.DeleteClauses = deleteClauses
	}
	if len(config.UpdateClauses) == 0 {
		config.UpdateClauses = updateClauses
	}

	createCallback := db.Callback().Create()
	createCallback.Register("gorm:create", Create(config))
	createCallback.Clauses = config.CreateClauses
}

// Create create hook
func Create(config *Config) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement.SQL.Len() == 0 {
			db.Statement.SQL.Grow(180)
			db.Statement.AddClauseIfNotExists(clause.Insert{})
			db.Statement.AddClause(ConvertToCreateValues(db.Statement))
			db.Statement.Build(db.Statement.BuildClauses...)
		}

		log.Printf("sql: %v", db.Statement.SQL.String())
		result, err := db.Statement.Conn.ExecContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
		if err != nil {
			db.AddError(err)
			return
		}
		db.RowsAffected, _ = result.RowsAffected()

		var (
			pkField *schema.Field
		)
		insertID, err := result.LastInsertId()
		insertOk := err == nil && insertID > 0
		if !insertOk {
			return
		}

		if db.Statement.Schema != nil {
			if db.Statement.Schema.PrioritizedPrimaryField == nil || !db.Statement.Schema.PrioritizedPrimaryField.HasDefaultValue {
				return
			}
			pkField = db.Statement.Schema.PrioritizedPrimaryField
		}

		if pkField == nil {
			return
		}

		switch db.Statement.ReflectValue.Kind() {
		case reflect.Struct:
			_, isZero := pkField.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
			if isZero {
				db.AddError(pkField.Set(db.Statement.Context, db.Statement.ReflectValue, insertID))
			}
		}
	}
}

func ConvertToCreateValues(stmt *gorm.Statement) (values clause.Values) {
	selectColumns, restricted := stmt.SelectAndOmitColumns(true, false)
	var isZero bool
	values = clause.Values{Columns: make([]clause.Column, 0, len(stmt.Schema.DBNames))}
	for _, dbName := range stmt.Schema.DBNames {
		if field := stmt.Schema.FieldsByDBName[dbName]; !field.HasDefaultValue {
			if v, ok := selectColumns[dbName]; (ok && v) || (!ok && !restricted) {
				values.Columns = append(values.Columns, clause.Column{Name: dbName})
			}
		}
	}

	values.Values = [][]interface{}{make([]interface{}, len(values.Columns))}
	for idx, column := range values.Columns {
		field := stmt.Schema.FieldsByDBName[column.Name]
		if values.Values[0][idx], isZero = field.ValueOf(stmt.Context, stmt.ReflectValue); isZero {
			if field.DefaultValueInterface != nil {
				values.Values[0][idx] = field.DefaultValueInterface
			}
		} else {
			values.Values[0][idx], _ = field.ValueOf(stmt.Context, stmt.ReflectValue)
		}
	}

	//for _, field := range stmt.Schema.FieldsWithDefaultDBValue {
	//	if v, ok := selectColumns[field.DBName]; (ok && v) || field.DefaultValueInterface == nil {
	//		if rv
	//	}
	//}
	return values
}

type Config struct {
	LastInsertIDReversed bool
	CreateClauses        []string
	QueryClauses         []string
	UpdateClauses        []string
	DeleteClauses        []string
}
