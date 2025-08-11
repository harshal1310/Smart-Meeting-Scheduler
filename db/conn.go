package db

import (
	"log"
	"smart-scheduler/model"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var once sync.Once

func InitDB(dburl string, dbName string) {
	once.Do(func() {
		// First connect to default postgres database to create our target database
		defaultURL := "host=localhost user=myuser dbname=postgres port=5432 password=mypassword sslmode=disable"
		defaultDB, err := gorm.Open(postgres.Open(defaultURL), &gorm.Config{})
		if err != nil {
			log.Printf("failed to connect to default database: %v", err)
		}

		var exists int64
		err = defaultDB.Raw("SELECT 1 FROM pg_database WHERE datname = ?", dbName).Scan(&exists).Error
		if err != nil {
			log.Fatalf("failed to check if database exists: %v", err)
		}
		if exists == 0 {
			log.Printf("Database '%s' does not exist. Creating...\n", dbName)
			if err := defaultDB.Exec("CREATE DATABASE " + dbName).Error; err != nil {
				log.Fatalf("failed to create database: %v", err)
			}
		} else {
			log.Printf("Database '%s' already exists. Skipping creation.\n", dbName)
		}

		// Now connect to our target database
		DB, err = gorm.Open(postgres.Open(dburl), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to target database: %v", err)
		}

		// Auto-migrate tables
		DB.AutoMigrate(&model.User{}, &model.Event{})
	})
}
