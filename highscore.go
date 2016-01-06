package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type AppServer struct {
	DBType     string
	DBUser     string
	DBPassword string
	DBIP       string
	Port       int
	DBPort     int
	Database   string
	DBOptions  []string
	DB         gorm.DB
}

// Struct that is stored in database
type HighScore struct {
	ID        uint      `gorm:"primary_key"`
	Name      string    `sql:"not null" json:"name"`
	Rank      string    `sql:"-"`
	Score     uint64    `sql:"not null" json:"score,string"`
	Round     uint      `sql:"not null" json:"round,string"`
	Seconds   uint      `sql:"not null" json:"time,string"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp" json:"date,string"`
}

func main() {
	server := AppServer{
		DBType:     "mysql",
		DBUser:     "highscores",
		DBPassword: "niko1niko",
		DBIP:       "192.168.0.20",
		Port:       8080,
		DBPort:     3306,
		Database:   "HighScores",
		DBOptions:  []string{"charset=utf8", "parseTime=true"},
	}
	server.Start()
}

func (app AppServer) Start() {
	app.InitDB()
	app.InitSchema()

	defer app.DB.Close()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/highscores/all", app.GetAllHighScores),
		rest.Get("/highscores/:count", app.GetHighScores),
		rest.Get("/highscores/", app.GetHighScores),
		rest.Post("/highscores/", app.PostHighScore),
	// rest.Get("/:name", m.GetPlayerScores),
	// rest.Delete("/:name/:id", m.Delete),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	log.Fatal(
		http.ListenAndServe(
			":"+strconv.Itoa(app.Port), api.MakeHandler()))
}

///////////////////////////////////////////////////////////////////////////////
// Database
///////////////////////////////////////////////////////////////////////////////
func (app *AppServer) InitDB() {
	var err error
	app.DB, err = gorm.Open(app.DBType, app.Source())
	if err != nil {
		log.Fatalf("Error connect to database, the error is '%v'", err)
	}
	app.DB.LogMode(true)
}

func (app *AppServer) InitSchema() {
	app.DB.AutoMigrate(&HighScore{})
}

func (app *AppServer) Source() string {
	switch app.DBType {
	case "mysql":
		dataString := fmt.Sprintf(
			"%s:%s@(%s:%s)/%s?%s",
			app.DBUser,
			app.DBPassword,
			app.DBIP,
			strconv.Itoa(app.DBPort),
			app.Database,
			strings.Join(app.DBOptions, "&"),
		)
		return dataString
	case "ps":
		return ""
	case "sqlite":
		return ""
	default:
		return ""
	}
}

///////////////////////////////////////////////////////////////////////////////
// Handlers
///////////////////////////////////////////////////////////////////////////////
func (app *AppServer) GetAllHighScores(w rest.ResponseWriter, r *rest.Request) {
	// var rank []int
	hss := []HighScore{}
	app.DB.Order("score DESC").Find(&hss)
	w.WriteJson(&hss)
}

func (app *AppServer) GetHighScores(w rest.ResponseWriter, r *rest.Request) {
	param := r.PathParam("count")
	count, err := strconv.Atoi(param)
	if err != nil {
		count = 50
	}
	// numbers, _ := regexp.Compile("[0-9]*")
	// if numbers.Match([]byte(param)) {
	// 	count, _ = strconv.Atoi(param)
	// } else {
	// 	count = 50
	// }

	hss := []HighScore{}
	app.DB.Order("score DESC").Limit(count).Find(&hss)
	w.WriteJson(&hss)
}

func (app *AppServer) PostHighScore(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var hs HighScore
	if err = r.DecodeJsonPayload(&hs); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// if !verifyScore(&hs) {
	// 	rest.Error(w, "invalid json", http.StatusBadRequest)
	// 	return
	// }

	if err := app.DB.Save(&hs).Error; err != nil {
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

func verifyScore(hs *HighScore) bool {
	return true
}
