package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

//	curl localhost:8080/event \
//	 -d '{
//	   "timestamp": 1710000000,
//	   "session_id":"abc",
//	   "channel_id":"sports",
//	   "event_type":"impression",
//	   "ad_id":"42"
//	}'

type Event struct {
	Timestamp int64  `json:"timestamp"`
	SessionID string `json:"session_id"`
	ChannelID string `json:"channel_id"`
	EventType string `json:"event_type"`
	AdID      string `json:"ad_id"`
}

func main() {
	http.HandleFunc("/event", handleEvent)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	var ev Event
	if err := json.Unmarshal(body, &ev); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	// fire and forget
	go exportToGCP(ev)

	w.WriteHeader(http.StatusNoContent)
}
