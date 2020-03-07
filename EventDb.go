package main

import (
	"log"
)

type Event struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
var events map[int64]Event

var eventIdCount int64

func initDb() {
	events = make(map [int64]Event)
	eventIdCount = 1
}

func findById(id int64) Event {
	log.Printf("returning event from DB : %v", events[id])

	return events[id]
}

func insert(event Event) Event {
	event.ID = eventIdCount
	eventIdCount++
	events[event.ID] = event
	log.Printf("Inserted event in DB : %v", events[event.ID])
	log.Printf("Record count : %d", len(events));

	return event
}

func update(id int64, event Event) Event {
	events[id] = event
	log.Printf("Updated event in DB : %v", events[event.ID])

	return event
}

func getAll() []Event {
	var values = make([]Event, len(events))
	idx := 0
	for _, value := range events {
		values[idx] = value
		idx++
	}

	return values
}

func deleteById(id int64) Event {
	var deletedEvent = events[id]
	delete(events, id)
	return deletedEvent
}