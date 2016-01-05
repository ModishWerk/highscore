package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const Port = 8080

type Model struct {
	DB gorm.DB
}

// Struct that is stored in database
type HighScore struct {
	ID        uint      `gorm:"primary_key"`
	Name      string    `sql:"not null" json:"name"`
	Score     uint64    `sql:"not null" json:"score,string"`
	Round     uint      `sql:"not null" json:"round,string"`
	Seconds   uint      `sql:"not null" json:"seconds,string"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp" json:"date,string"`
}

func main() {
	m := Model{}
	m.InitDB()
	m.InitSchema()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/", m.GetHighScores),
		rest.Post("/", m.PostHighScore),
	// rest.Get("/:name", m.GetPlayerScores),
	// rest.Delete("/:name/:id", m.Delete),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Port), api.MakeHandler()))
}

///////////////////////////////////////////////////////////////////////////////
// Database
///////////////////////////////////////////////////////////////////////////////
func (m *Model) InitDB() {
	var err error
	dataSource := "highscores:niko1niko@(192.168.0.20:3306)/HighScores"
	m.DB, err = gorm.Open("mysql", dataSource)
	if err != nil {
		log.Fatalf("Error connect to database, the error is '%v'", err)
	}
	m.DB.LogMode(true)
}

func (m *Model) InitSchema() {
	m.DB.AutoMigrate(&HighScore{})
}

///////////////////////////////////////////////////////////////////////////////
// Handlers
///////////////////////////////////////////////////////////////////////////////
func (m *Model) GetHighScores(w rest.ResponseWriter, r *rest.Request) {
	hss := []HighScore{}
	m.DB.Find(&hss).Order("score DESC").Limit(2)
	w.WriteJson(&hss)
}

func (m *Model) PostHighScore(w rest.ResponseWriter, r *rest.Request) {
	hs := HighScore{}
	if err := r.DecodeJsonPayload(&hs); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := m.DB.Save(&hs).Error; err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteJson(&hs)
}

///////////////////////////////////////////////////////////////////////////////
// Helpers
///////////////////////////////////////////////////////////////////////////////
func NumberToRank(n int) string {
	switch {
	case n <= 0:
		return ""
	case n == 1:
		return "1st"
	case n == 2:
		return "2nd"
	case n == 3:
		return "3rd"
	default:
		return fmt.Sprintf("%sth", strconv.Itoa(n))
	}
}

func verifyScore(hs HighScore) {
}
