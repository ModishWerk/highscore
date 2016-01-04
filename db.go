package main

import (
	"log"

	"github.com/jinzhu/gorm"
)

type Model struct {
	DB gorm.DB
}

func (m *Model) InitDB() {
	var err error
	var dataSource string = "highscores:niko1niko@(192.168.0.20:3306)/HighScores"
	m.DB, err = gorm.Open("mysql", dataSource)
	if err != nil {
		log.Fatalf("Error connect to database, the error is '%v'", err)
	}
	m.DB.LogMode(true)
}

func (m *Model) InitSchema() {
	m.DB.AutoMigrate(&HighScore{})
}
