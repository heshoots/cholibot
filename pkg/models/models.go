package models

import (
//"github.com/jinzhu/gorm"
//_ "github.com/jinzhu/gorm/dialects/sqlite"
)

/*
type DiscordGuild struct {
	gorm.Model
	GuildID string `gorm:"unique;not null"`
}

type DiscordRole struct {
	gorm.Model
	RoleID  string `gorm:"unique;not null"`
	GuildID string `gorm:"not null"`
	Name    string
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
}

func Create() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
}

func Migrate() {
	// Migrate the schema
	db.AutoMigrate(&DiscordGuild{})
	db.AutoMigrate(&DiscordRole{})
}

func GetDB() *gorm.DB {
	return db
}*/
