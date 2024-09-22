package database

import (
	"database/sql"
	"internal/database"

	"gorm.io/gorm"
)

var MIGRATOR_INTERFACE gorm.Migrator

func ConnectToDatabase() *sql.DB {
	database.InitDatabase()

	migrator, db := database.GetMigratorAndDbInstance()

	MIGRATOR_INTERFACE = migrator

	return db
}
