package database

import (
	"fmt"
	"internal/configuration"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DATABASE *gorm.DB
)

func InitDatabase() {
	heavenlyDragonsDatabase()
}

func heavenlyDragonsDatabase() *gorm.DB {
	var err error

	if DATABASE == nil {
		DATABASE, err = initDatabaseConnection()

		if err != nil {
			slog.With(err).Error("Could not connect to Heavenly Dragons database")
			panic(err)
		}

		slog.Info("Connected to Heavenly Dragons database")
	}

	return DATABASE
}

func initDatabaseConnection() (db *gorm.DB, err error) {
	config := configuration.GetConfiguration()
	var (
		dsn string
	)

	// dbPortS := os.Getenv("DB_PORT")
	// dbPort, err := strconv.Atoi(dbPortS)
	// if err != nil {
	// 	panic("Missing DB_PORT")
	// }

	if !config.IsLocalEnv() {
		// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d application_name=",
		// 	os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), dbPort)

	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d application_name=HeavenlyDragons",
			"127.0.0.1", "local", "local", "local", 5432,
		)
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: slog.Logger,
	})
	if err != nil {
		return db, err
	}

	return db, nil
}
