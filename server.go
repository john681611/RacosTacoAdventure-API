package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Score ...
type Score struct {
	ID    int       `json:"id"`
	Score int       `json:"score"`
	Date  time.Time `json:"date"`
}

var scoreList []Score

func main() {
	//routing
	http.HandleFunc("/", ping)
	http.HandleFunc("/addScore", addScore)
	http.HandleFunc("/getBoard", getBoard)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "A Go Web Server")
	w.WriteHeader(200)
}

func addScore(w http.ResponseWriter, r *http.Request) {
	scoreList = append(scoreList, Score{1, 100, time.Now()})
	w.Header().Set("Server", "A Go Web Server")
	w.WriteHeader(200)
}

func getBoard(w http.ResponseWriter, r *http.Request) {
	js, _ := json.Marshal(scoreList)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ResetScore() {
	scoreList = scoreList[:0]
}
