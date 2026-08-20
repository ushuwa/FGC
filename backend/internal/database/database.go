package database

import (
	"fmt"
	"log"

	"github.com/jj.jobo/FGC/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.App.DBHost,
		config.App.DBUser,
		config.App.DBPassword,
		config.App.DBName,
		config.App.DBPort,
		config.App.DBSSLMode,
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = db

	sqlDB, err := DB.DB()

	if err != nil {
		log.Fatal("Failed to get database connection: ", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	log.Println("Database Connected")
}

func CloseDatabase() {

	if DB == nil {
		return
	}

	sqlDB, err := DB.DB()

	if err != nil {
		log.Println("Failed to get database connection for closing:", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Println("Failed to close database connection:", err)
		return
	}

	log.Println("Database connection closed")
}
