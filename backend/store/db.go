package store

import (
	"fmt"
	"github.com/commonpool/backend/auth"
	chatstore "github.com/commonpool/backend/pkg/chat/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

func NewTestDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./realworld_test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&chatstore.Channel{},
		&chatstore.ChannelSubscription{},
		&chatstore.Message{},
		&auth.User{},
		&TransactionEntry{},
	)
	if err != nil {
		panic(err)
	}
}

func DropTestDB() error {
	if err := os.Remove("./realworld_test.db"); err != nil {
		return err
	}
	return nil
}
