package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Score ...
type Score struct {
	Name  string    `json:"name"`
	Score int       `json:"score"`
	Date  time.Time `json:"date,omitempty"`
}

var scoreList []Score

func main() {
	//routing
	resetScore()
	http.HandleFunc("/addScore", addScore)
	http.HandleFunc("/getBoard", getBoard)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func addScore(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var newScore Score

		err := decoder.Decode(&newScore)

		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		scoreList = append(scoreList, Score{newScore.Name, newScore.Score, time.Now()})
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, "POST done")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func getBoard(w http.ResponseWriter, r *http.Request) {
	js, _ := json.Marshal(scoreList)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func resetScore() {
	scoreList = scoreList[:0]
}
