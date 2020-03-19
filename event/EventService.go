package event

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kislayverma/go-crud/event/dao"
	"github.com/kislayverma/go-crud/http-utils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type EventService struct {
	Db     dao.IEventDao
	Router *mux.Router
}

func (svc EventService) RegisterRoutes() {
	subRouter := svc.Router.PathPrefix("/event").Subrouter()
	subRouter.HandleFunc("", svc.createEvent).Methods("POST")
	subRouter.HandleFunc("", svc.getAllEvents).Methods("GET")
	subRouter.HandleFunc("/{id}", svc.getOneEvent).Methods("GET")
	subRouter.HandleFunc("/{id}", svc.updateEvent).Methods("PATCH")
	subRouter.HandleFunc("/{id}", svc.deleteEvent).Methods("DELETE")
	subRouter.HandleFunc("/validate/{id}", svc.validateEventExists).Methods("GET")
}

func (svc EventService) getOneEvent(w http.ResponseWriter, r *http.Request) {
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
		var event dao.Event = svc.Db.FindById(eventID)
		if event.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(event)
		}
	}
}

func (svc EventService) getAllEvents(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /event Event getAllEvents
	// ---
	// summary: Return all the events from the store
	// description: If event are found array of events will be returned, else an empty JSON array will be returned.
	// parameters:
	// responses:
	//   "200":
	//     "$ref": "#/responses/accountRes"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(svc.Db.GetAll())
}

func (svc EventService) createEvent(w http.ResponseWriter, r *http.Request) {
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
	var newEvent dao.Event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Kindly enter data with the event title and description only in order to update")
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.Unmarshal(reqBody, &newEvent)
	log.Printf("Incoming event %v", newEvent)
	svc.Db.Insert(newEvent)

	w.WriteHeader(http.StatusCreated)
}

func (svc EventService) updateEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		log.Println("Failed to parse the id")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		var updatedEvent dao.Event

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Kindly enter data with the event title and description only in order to update")
		}
		json.Unmarshal(reqBody, &updatedEvent)

		var existingEvent dao.Event = svc.Db.FindById(eventId)
		if existingEvent.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			existingEvent.Title = updatedEvent.Title
			existingEvent.Description = updatedEvent.Description
			svc.Db.Update(eventId, existingEvent)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (svc EventService) deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		log.Println("Failed to parse the id")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		log.Printf("Trying to delete event with id %d", eventId)
		var existingEvent dao.Event = svc.Db.FindById(eventId)
		if existingEvent.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			svc.Db.DeleteById(eventId)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (svc EventService) validateEventExists(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	reqContext := http_utils.GetRequestContext(r, r.Header[http_utils.REQUEST_ID_HEADER_NAME][0])

	if err != nil {
		log.Println(reqContext.RequestId, "Failed to parse the id")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		responseChannel := make(chan bool)
		go getEvent(reqContext, eventId, responseChannel)
		eventFound := <- responseChannel
		w.WriteHeader(http.StatusOK)
		if eventFound {
			go logResponse(reqContext, "Event With id " + strconv.FormatInt(eventId, 10) + " was found")
		} else {
			go logResponse(reqContext,"Event With id " + strconv.FormatInt(eventId, 10) + " was not found")
		}
	}
}
