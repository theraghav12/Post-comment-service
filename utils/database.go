package utils

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"post-comments-api/config"
)

var (
	db   *gorm.DB
	dbMu sync.Mutex
)

func InitDB() {
	dbMu.Lock()
	defer dbMu.Unlock()
	if db != nil {
		return
	}
	cfg := config.AppConfig
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db = database
}

func GetDB() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}
