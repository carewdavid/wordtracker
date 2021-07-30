package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/carewdavid/wordtracker/record"
	_ "github.com/mattn/go-sqlite3"
)

const version = "1.0.0"

var schema = "CREATE TABLE IF NOT EXISTS wordcount(date INTEGER NOT NULL, words INTEGER NOT NULL, desc STRING)"
var db *sql.DB

//Current time, rounded down to the beginning of the day
func today() time.Time {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return today
}

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

func daily(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	stmt, err := db.Prepare("SELECT SUM(words) FROM wordcount where date >= ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	row := stmt.QueryRow(today().Unix())
	var total int
	err = row.Scan(&total)

	if err != nil {
		total = 0
	}

	data := struct {
		Words int `json:words`
	}{}
	data.Words = total
	buf, _ := json.Marshal(data)
	fmt.Fprintf(w, string(buf))

}

func since(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := r.URL.Query()
	t := params.Get("t")
	if t == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No time given")
		return
	}
	timestamp, err := strconv.ParseInt(t, 10, 64)

	stmt, err := db.Prepare("SELECT SUM(words) FROM wordcount WHERE date >= ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}

	row := stmt.QueryRow(t)

	var total int
	err = row.Scan(&total)

	if err != nil {
		total = 0
	}

	now := today()
	then := time.Unix(timestamp, 0)
	//Get elapsed time in days. Yucky magic number is the number of seconds in
	//one day sind Duration doesn't have a method to get the right number natively
	const SECONDS_PER_DAY = 86400
	elapsed := now.Sub(then).Round(SECONDS_PER_DAY).Seconds() / SECONDS_PER_DAY

	data := struct {
		Words   int     `json:words`
		Average float64 `json:average`
	}{}
	data.Words = total
	//No, we can't just use SQl AVERAGE(). The number of records is not the same
	//as the number of days in the time period (probably) so that will give the wrong answer
	data.Average = float64(total) / elapsed
	buf, _ := json.Marshal(data)
	fmt.Fprintf(w, string(buf))
}

func newRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not read request.")
		return
	}
	var rec record.Record
	err = json.Unmarshal(requestBody, &rec)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not read incoming record.")
		return
	}
	stmt, err := db.Prepare("INSERT INTO wordcount VALUES(?, ?, ?)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(rec.Date, rec.Words, rec.Desc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func serve() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/today", daily)
	http.HandleFunc("/update", newRecord)
	http.HandleFunc("/since", since)
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
