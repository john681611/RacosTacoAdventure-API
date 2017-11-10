package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"sort"
	mgo "gopkg.in/mgo.v2"
)

// Score ...
type Score struct {
	Name  string    `json:"name"`
	Score int       `json:"score"`
	Date  time.Time `json:"date,omitempty"`
}

type ByScore []Score

func (c ByScore) Len() int           { return len(c) }
func (c ByScore) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByScore) Less(i, j int) bool { return c[i].Score > c[j].Score }

var scoreList []Score
var ignoreDB = false

func main() {
	//routing
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
			fmt.Println("Error: Addscore: ", err)
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		newScore =  Score{newScore.Name, newScore.Score, time.Now()}
		if ignoreDB {
			scoreList = append(scoreList, newScore)
		} else {
			addScoretoDB(newScore)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, "POST done")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func getBoard(w http.ResponseWriter, r *http.Request) {
	 var response []Score
	if ignoreDB {
		response = scoreList
	} else {
		response = getDBLeaderBoard()
	}
	sort.Sort(ByScore(response))
	js, _ := json.Marshal(response)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func resetScore() {
	scoreList = scoreList[:0]
}

func setIgnoreDB(x bool) {
	ignoreDB = x
}

func getDBLeaderBoard() []Score {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("RacosDB").C("leaderBoard")
	result := []Score{}
	err = c.Find(nil).All(&result)
	if err != nil {
		log.Fatal(err)
	} else {
		return result
	}
	return []Score{}
}

func addScoretoDB(newScore Score) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("RacosDB").C("leaderBoard")

	err = c.Insert(newScore)
	if err != nil {
		log.Fatal(err)
	}
}
