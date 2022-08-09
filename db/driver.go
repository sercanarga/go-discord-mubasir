package db

import (
	"flag"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	BotToken = flag.String("token", "", "Bot access token")
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("mubasir.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = DB.AutoMigrate(&Users{})
	if err != nil {
		return
	}
}
