// Code Credit : https://medium.com/the-andela-way/build-a-restful-json-api-with-golang-85a83420c9da
package main

import (
    "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent Event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)
	log.Printf("Incoming event %v", newEvent)
	newEvent = insert(newEvent)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	if err != nil {
		log.Println("Failed to parse the id")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var event = findById(eventID)
		if event.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(event)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(getAll())
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	if err != nil {
		log.Println("Failed to parse the id")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var updatedEvent Event

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Kindly enter data with the event title and description only in order to update")
		}
		json.Unmarshal(reqBody, &updatedEvent)

		var existingEvent = findById(eventId)
		if existingEvent.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			existingEvent.Title = updatedEvent.Title
			existingEvent.Description = updatedEvent.Description
			update(eventId, existingEvent)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(existingEvent)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	if err != nil {
		log.Println("Failed to parse the id")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var existingEvent = findById(eventId)
		if existingEvent.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			deleteById(eventId)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(existingEvent)
		}
	}
}

func main() {
	initDb()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/event", getAllEvents).Methods("GET")
	router.HandleFunc("/event/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/event/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/event/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}