package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Record struct {
	Date  int    `json:"date"`
	Words int    `json:"words"`
	Desc  string `json:"desc"`
}

const version = "1.0.0"

var schema = "CREATE TABLE IF NOT EXISTS wordcount(date INTEGER NOT NULL, words INTEGER NOT NULL, desc STRING)"

func initDatabase(path string) *sql.DB {
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("Could not open database.")
	}
	stmt, _ := database.Prepare(schema)
	stmt.Exec()

	return database
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wordtracker %s", version)
}

func serve() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	dbPath := os.Getenv("DBPATH")
	if dbPath == "" {
		log.Fatal("No path to database. Make sure DBPATH is set.")
	}
	initDatabase(dbPath)
	serve()

}
