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
	Driver     string
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
	Rank      int       `sql:"-"`
	Score     uint64    `sql:"not null" json:"score,string"`
	Round     uint      `sql:"not null" json:"round,string"`
	Seconds   uint      `sql:"not null" json:"time,string"`
	CreatedAt time.Time `sql:"DEFAULT:current_timestamp" json:"date,string"`
}

func main() {
	server := AppServer{
		Driver:     "mysql",
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
		rest.Get("/api/highscores", app.GetUrls),
		rest.Get("/api/highscores/all", app.GetAllHighScores),
		rest.Get("/api/highscores/:count", app.GetHighScores),
		rest.Get("/api/highscores/:count/:offset", app.GetHighScoresRange),
		rest.Post("/api/highscores", app.PostHighScore),
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
	app.DB, err = gorm.Open(app.Driver, app.Conn())
	if err != nil {
		log.Fatalf("Error connect to database, the error is '%v'", err)
	}
	app.DB.LogMode(true)
}

func (app *AppServer) InitSchema() {
	app.DB.AutoMigrate(&HighScore{})
}

func (app *AppServer) Conn() string {
	switch app.Driver {
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
func (app *AppServer) GetUrls(w rest.ResponseWriter, r *rest.Request) {
	var host string = "https://raviko.dlinkddns.com/api/highscores"
	w.WriteJson(struct {
		AllHighscores string `json:"get_all_highscores_url"`
		GetTopScores  string `json:"get_n_scores_url"`
		GetScoreRange string `json:"get_n_score_with_offset_url"`
		PostHighScore string `json:"post_highscore_url"`
	}{
		host + "/all",
		host + "/{count}",
		host + "/{count}/{offset}",
		host + "/",
	})
}

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

	hss := []HighScore{}
	app.DB.Order("score DESC").Limit(count).Find(&hss)
	w.WriteJson(&hss)
}

func (app *AppServer) GetHighScoresRange(w rest.ResponseWriter, r *rest.Request) {
	countParam := r.PathParam("count")
	offsetParam := r.PathParam("offset")
	count, err := strconv.Atoi(countParam)
	if err != nil {
		count = 50
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		offset = 0
	}

	hss := []HighScore{}
	app.DB.Order("score DESC").Offset(offset).Limit(count).Find(&hss)
	w.WriteJson(&hss)
}

func (app *AppServer) PostHighScore(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var hs HighScore

	if err = r.DecodeJsonPayload(&hs); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = hs.VerifyScore(); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

func (hs HighScore) VerifyScore() error {
	if err := checkMissingFields(hs); err != nil {
		return err
	}

	return nil
}

func checkMissingFields(hs HighScore) error {
	var missingFields []string

	if hs.Name == "" {
		missingFields = append(missingFields, "name")
	}
	if hs.Round == 0 {
		missingFields = append(missingFields, "round")
	}
	if hs.Score == 0 {
		missingFields = append(missingFields, "score")
	}
	if hs.Seconds == 0 {
		missingFields = append(missingFields, "seconds")
	}

	count := len(missingFields)
	switch count {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("%s must be present",
			strings.Title(missingFields[0]))
	case 2:
		return fmt.Errorf("%s and %s must be present",
			strings.Title(missingFields[0]),
			missingFields[1],
		)
	default:
		return fmt.Errorf("%s, %s, and %s must be present",
			strings.Title(missingFields[0]),
			strings.Join(missingFields[1:count-1], ", "),
			missingFields[count-1],
		)
	}
}
