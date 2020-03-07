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

var db IEventDao

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent Event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Kindly enter data with the event title and description only in order to update")
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.Unmarshal(reqBody, &newEvent)
	log.Printf("Incoming event %v", newEvent)
	newEvent = getDao().insert(newEvent)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		log.Println("Failed to parse the id")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var event = getDao().findById(eventID)
		if event.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(event)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(getDao().getAll())
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

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

		var existingEvent = getDao().findById(eventId)
		if existingEvent.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			existingEvent.Title = updatedEvent.Title
			existingEvent.Description = updatedEvent.Description
			getDao().update(eventId, existingEvent)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(existingEvent)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		log.Println("Failed to parse the id")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var existingEvent = getDao().findById(eventId)
		if existingEvent.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			getDao().deleteById(eventId)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(existingEvent)
		}
	}
}

func getDao() IEventDao {
	return db
}

func main() {
	// Initialize the database access layer
	db = EventInMemoryDao{make(map [int64]Event), 1}
	db.initDb()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/event", getAllEvents).Methods("GET")
	router.HandleFunc("/event/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/event/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/event/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}