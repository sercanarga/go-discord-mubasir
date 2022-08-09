package db

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	DiscordId     string
	CustomMessage string
	IsAdmin       bool
}
