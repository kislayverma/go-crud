package main

import (
"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Code Credit : https://medium.com/the-andela-way/build-a-restful-json-api-with-golang-85a83420c9da
type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	
	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(newEvent)
}
