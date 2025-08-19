package callbacks

import gorm "simpleGorm"

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
	createCallback.Register()
}

type Config struct {
	LastInsertIDReversed bool
	CreateClauses        []string
	QueryClauses         []string
	UpdateClauses        []string
	DeleteClauses        []string
}
