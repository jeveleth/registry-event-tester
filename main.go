package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator"
	"github.com/hashicorp/go-retryablehttp"
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

var slackWebhookURL = os.Getenv("SLACK_WEBHOOK_URL")

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

	var fmtEvent string
	for _, event := range events.Events {
		fmtEvent = fmt.Sprintf(`
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

	notifySlack(fmtEvent)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(events)
}

func notifySlack(message string) {
	body := ` { "text" : "New image pushed to registry` + message + `."}`

	req, err := retryablehttp.NewRequest("POST", slackWebhookURL, strings.NewReader(body))
	checkErr(err)
	req.Header.Add("content-type", "application/json")

	res, err := retryablehttp.NewClient().Do(req)
	checkErr(err)

	defer res.Body.Close()

	checkErr(err)

	if res.StatusCode != 200 {
		log.Printf("non 200 status Code from Slack: %v\n", res.StatusCode)
		os.Exit(1)
	}

	log.Printf("response from slack is is %v", body)
}

func checkErr(err error) {
	if err != nil {
		log.Println(fmt.Sprintf("an error occurred:  %v", err.Error()))
		os.Exit(1)
	}
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	http.HandleFunc("/", GetEvent)

	http.ListenAndServe(":"+port, nil)
}
