package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbConfig struct {
	Path string
	Db   *gorm.DB
}

func (d *DbConfig) InitDatabase() error {
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 5 * time.Second,
		},
	)
	dbConfig := gorm.Config{
		Logger:          dbLogger,
		CreateBatchSize: 500,
	}

	log.Printf("Opening SQLite db at: %s", d.Path)
	db, err := gorm.Open(sqlite.Open(d.Path), &dbConfig)
	if err != nil {
		return fmt.Errorf("Error opening database: %w", err)
	}
	if err = db.AutoMigrate(&User{}, &Place{}, &Event{}, &Ticket{}); err != nil {
		return fmt.Errorf("Error auto-migrating sqlite database: %w", err)
	}

	d.Db = db
	return nil
}
