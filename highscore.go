package main

import (
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
)

type HighScore struct {
	ID    uint      `gorm:"primary_key"`
	Name  string    `sql:"not null"`
	Score uint64    `sql:"not null"`
	Date  time.Time `sql:"DEFAULT:current_timestamp"`
}

type JsonHS struct {
	Name  string    `json:"name"`
	Rank  string    `json:"rank"`
	Score uint64    `json:"score,string"`
	Date  time.Time `json:"date,string"`
}

func reduceHighScores(hss []HighScore) []JsonHS {
	var jsons []JsonHS
	for rank, hs := range hss {
		var json JsonHS
		json = reduceHighScore(rank, hs)
		jsons = append(jsons, json)
	}
	return jsons
}

func reduceHighScore(rank int, hs HighScore) JsonHS {
	json := JsonHS{}
	json.Name = hs.Name
	json.Rank = number_to_rank(rank)
	json.Score = hs.Score
	json.Date = hs.Date

	return json
}

func (m *Model) GetHighScores(w rest.ResponseWriter, r *rest.Request) {
	hss := []HighScore{}
	m.DB.Find(&hss).Order("score DESC").Limit(2)
	jsons := reduceHighScores(hss)
	w.WriteJson(&jsons)
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

	json := reduceHighScore(0, hs)
	w.WriteJson(&json)
}
