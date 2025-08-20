package gorm

func (db *DB) Create(value interface{}) (tx *DB) {
	db.Statement.Dest = value
	return db.callbacks.Create().Execute(tx)
}
