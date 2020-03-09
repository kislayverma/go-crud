// Go-CRUD Service
//
// Sample CRUD service built in Go, to learn Go
//
//     Schemes: http
//     Host: localhost:8080
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Kislay<kislay.nsit@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"database/sql"
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
	// Code Credit : https://medium.com/the-andela-way/build-a-restful-json-api-with-golang-85a83420c9da

	// swagger:operation GET / Hello Hello
	//
	// Returns a simple Hello message
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - text/plain
	// responses:
	//   '200':
	//     description: The hello message
	//     type: string
    fmt.Fprintf(w, "Welcome home!")
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /event/ Event createEvent
	// ---
	// summary: Creates an event with the give title and description and assign a new id to it
	// description: Creates an event with the give title and description and assign a new id to it
	// parameters:
	// - name: event
	//   in: body
	//   description: title and description of the event
	//   type: event
	//   required: true
	//   "$ref": "#/responses/event"
	// responses:
	//   "200":
	//     "$ref": "#/responses/event"
	//   "400":
	//     "$ref": "#/responses/event"
	//   "404":
	//     "$ref": "#/responses/notFoundReq"
	var newEvent Event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Kindly enter data with the event title and description only in order to update")
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.Unmarshal(reqBody, &newEvent)
	log.Printf("Incoming event %v", newEvent)
	db.insert(newEvent)

	w.WriteHeader(http.StatusCreated)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /event/{id} Event getEvent
	// ---
	// summary: Return the event identified by the id.
	// description: If the event is found, it will be returned else Error Not Found (404) will be returned.
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the event
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/event"
	//   "400":
	//     "$ref": "#/responses/event"
	//   "404":
	//     "$ref": "#/responses/notFoundReq"
	eventID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		log.Println("Failed to parse the id")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var event = db.findById(eventID)
		if event.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(event)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /event Event getAllEvents
	// ---
	// summary: Return all the events from the store
	// description: If event are found array of events will be returned, else an empty JSON array will be returned.
	// parameters:
	// responses:
	//   "200":
	//     "$ref": "#/responses/accountRes"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(db.getAll())
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

		var existingEvent = db.findById(eventId)
		if existingEvent.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			existingEvent.Title = updatedEvent.Title
			existingEvent.Description = updatedEvent.Description
			db.update(eventId, existingEvent)
			w.WriteHeader(http.StatusOK)
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
		log.Printf("Trying to delete event with id %d", eventId)
		var existingEvent = db.findById(eventId)
		if existingEvent.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			db.deleteById(eventId)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func swagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "swagger.json")
}

func main() {
	// Initialize the database access layer
	// Not the best : conn details leaking into main program - need to enable configs in Mysql dao
	// Swagger setup creadit :
	// https://medium.com/@supun.muthutantrige/lets-go-everything-you-need-to-know-about-creating-a-restful-api-in-go-part-iv-52666c5221d4
	dbConn, err := sql.Open("mysql", "root:@/gocrud")
	if err != nil {
		panic(err.Error())
	}
	db = EventMySqlDao{dbConn}
	db.initDb()
	defer dbConn.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/swagger.json", swagger).Methods("GET")
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/event", getAllEvents).Methods("GET")
	router.HandleFunc("/event/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/event/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/event/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}