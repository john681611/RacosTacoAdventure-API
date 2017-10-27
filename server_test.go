package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func serverTest(method string, route string, h http.HandlerFunc, reader io.Reader, expected string, expectedStatus int, t *testing.T) {

	req, err := http.NewRequest(method, route, reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h)

	handler.ServeHTTP(rr, req)

	if rr.Code != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, expectedStatus)
	}

	if rr.Body.String() != expected && expected != "" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAddScoreHandlerGoodPost(t *testing.T) {
	body, _ := json.Marshal(Score{100, time.Now()})
	serverTest("POST", "/addScore", addScore, bytes.NewReader(body), "", http.StatusOK, t)
}
func TestAddScoreHandlerBadPost(t *testing.T) {
	serverTest("POST", "/addScore", addScore, bytes.NewReader([]byte("badString")), "", http.StatusInternalServerError, t)
}

func TestAddScoreHandlerGet(t *testing.T) {
	serverTest("GET", "/addScore", addScore, nil, "", http.StatusMethodNotAllowed, t)
}

func TestGetBoardHandler(t *testing.T) {
	resetScore()
	serverTest("GET", "/getBoard", getBoard, nil, "[]", http.StatusOK, t)
}

func TestResetScoreHandler(t *testing.T) {
	scoreList = append(scoreList, Score{100, time.Now()})
	resetScore()
	if len(scoreList) != 0 {
		t.Errorf("resetScore did not clear scores: got %v want %v",
			len(scoreList), 0)
	}
}
