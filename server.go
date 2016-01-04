package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)
import "github.com/ant0ine/go-json-rest/rest"

func main() {
	m := Model{}
	m.InitDB()
	m.InitSchema()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/", m.GetHighScores),
		rest.Post("/", m.PostHighScore),
	// rest.Get("/:playerId", i.GetReminder),
	// rest.Delete("/:id", i.DeleteReminder),
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
