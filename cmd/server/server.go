package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/carewdavid/wordtracker/record"
	_ "github.com/mattn/go-sqlite3"
)

const version = "1.0.0"

var schema = "CREATE TABLE IF NOT EXISTS wordcount(date INTEGER NOT NULL, words INTEGER NOT NULL, desc STRING)"
var db *sql.DB

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

func newRecord(w http.ResponseWriter, r *http.Request) {
	//Just assume it's a POST. Sure it's lazy but it's not like anyone else is going to be writing clients
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Could not read request.")
	}
	var rec record.Record
	err = json.Unmarshal(requestBody, &rec)
	if err != nil {
		log.Fatal("Could not read incoming record.")
	}
	stmt, err := db.Prepare("INSERT INTO wordcount VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(rec.Date, rec.Words, rec.Desc)
	if err != nil {
		log.Fatal("Error inserting into database.")
	}
}

func serve() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/update", newRecord)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	dbPath := os.Getenv("DBPATH")
	if dbPath == "" {
		log.Fatal("No path to database. Make sure DBPATH is set.")
	}
	db = initDatabase(dbPath)
	serve()

}