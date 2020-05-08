package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator"
)

type Target struct {
	MediaType  string `json:"mediaType"`
	Size       int    `json:"size"`
	Digest     string `json:"digest"`
	Length     int    `json:"length"`
	Repository string `json:"repository"`
	URL        string `json:"url"`
	Tag        string `json:"tag"`
}

type Events struct {
	Events []Event `json:"events"`
}
type Event struct {
	ID        string   `json:"id"`
	TimeStamp string   `json:"timestamp"`
	Action    string   `json:"action"`
	Target    []Target `json:"target"`
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	events := Events{}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Printf("error reading request body: %v", readErr)
		http.Error(w, http.StatusText(400), 400)
	}

	if err := json.Unmarshal(body, &events); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	validate := validator.New()
	err := validate.Struct(&events)
	if err != nil {
		log.Printf("error validating struct: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("event is %v\n", events)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(events)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	http.HandleFunc("/", GetEvent)

	http.ListenAndServe(":"+port, nil)
}
