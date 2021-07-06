package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Record struct {
	Date  int64  `json:date`
	Words int    `json:date`
	Desc  string `json:desc`
}

func add(words int, desc string) {
	now := time.Now().Local()
	timestamp := now.Unix()
	record := Record{timestamp, words, desc}
	buf, _ := json.Marshal(record)
	data := bytes.NewBuffer(buf)
	http.Post("http://localhost:10000/update", "application/json", data)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: $ client words [desc]")
		os.Exit(1)
	}
	words, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "First argument must be integer")
		os.Exit(1)
	}
	var desc string = ""
	if len(os.Args) > 2 {
		desc = os.Args[2]
	}
	add(words, desc)

}
