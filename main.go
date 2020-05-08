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
	ID        string `json:"id"`
	TimeStamp string `json:"timestamp"`
	Action    string `json:"action"`
	Target    Target `json:"target"`
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	events := Events{}

	body, readErr := ioutil.ReadAll(r.Body)
	log.Printf("body is %v \n", string(body))
	if readErr != nil {
		log.Printf("error reading request body: %v", readErr)
		http.Error(w, http.StatusText(400), 400)
	}

	if err := json.Unmarshal(body, &events); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err := validate.Struct(&events)
	if err != nil {
		log.Printf("error validating struct: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, event := range events.Events {
		log.Printf(`docker registry event ===>
		ID: %v
		Time: %v
		Action: %v
		Size: %v
		Digest: %v
		Length: %v
		Repository: %v
		URL: %v
		Tag: %v`,
			event.ID,
			event.TimeStamp,
			event.Action,
			event.Target.Size,
			event.Target.Digest,
			event.Target.Length,
			event.Target.Repository,
			event.Target.URL,
			event.Target.Tag,
		)
	}

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
