package store

import (
	"fmt"
	"github.com/commonpool/backend/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

func NewTestDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./realworld_test.db"), &gorm.Config{
	})
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(3)
	return db
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Resource{},
		&model.User{},
		&model.Thread{},
		&model.Message{},
		&model.Topic{},
		&model.ResourceTopic{},
	)
}

func DropTestDB() error {
	if err := os.Remove("./realworld_test.db"); err != nil {
		return err
	}
	return nil
}
