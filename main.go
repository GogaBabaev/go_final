package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

func main() {

	webDir := "./web"

	dbFile := "scheduler.db"

	envDBFile := os.Getenv("TODO_DBFILE")

	if envDBFile != "" {
		dbFile = envDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	port := "7540"

	envPort := os.Getenv("TODO_PORT")

	if envPort != "" {
		port = envPort
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.ListenAndServe(":"+port, nil)

}
