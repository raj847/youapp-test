package utils

import (
	"os"
	"youapp/entity"

	// _ "github.com/jackc/pgx/v4/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB() error {
	// connect using gorm pgx
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "pgx",
		DSN:        os.Getenv("DATABASE_URL"),
	}), &gorm.Config{})
	if err != nil {
		return err
	}

	conn.AutoMigrate(
		entity.User{},
		entity.Profile{},
	)

	SetupDBConnection(conn)
	return nil
}

func SetupDBConnection(DB *gorm.DB) {
	db = DB
}

func GetDBConnection() *gorm.DB {
	return db
}
