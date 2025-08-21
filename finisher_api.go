package gorm

func (db *DB) Create(value interface{}) *DB {
	db.Statement.Dest = value
	return db.callbacks.Create().Execute(db)
}
