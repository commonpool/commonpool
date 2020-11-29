package store

import (
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
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

	sqlDB.SetMaxIdleConns(3)
	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&chat.Channel{},
		&chat.ChannelSubscription{},
		&group.Group{},
		&group.Membership{},
		&Message{},
		&trading.Offer{},
		&trading.OfferItem{},
		&trading.OfferDecision{},
		&resource.Resource{},
		&resource.Sharing{},
		&auth.User{},
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
