package main

import (
	"log"
)

type EventInMemoryDao struct {
	events map[int64]Event
	eventIdCount int64
}

func (dao EventInMemoryDao) initDb() {
	// do nothing
}

func (dao EventInMemoryDao) findById(id int64) Event {
	log.Printf("returning event from DB : %v", dao.events[id])

	return dao.events[id]
}

func (dao EventInMemoryDao) insert(event Event) Event {
	var newEventForDb = Event{int64(len(dao.events)) + 1, event.Title, event.Description}
	dao.events[newEventForDb.ID] = newEventForDb
	log.Printf("Inserted event in DB : %v", dao.events[newEventForDb.ID])
	log.Printf("Record count : %d", len(dao.events));

	return newEventForDb
}

func (dao EventInMemoryDao) update(id int64, event Event) Event {
	dao.events[id] = event
	log.Printf("Updated event in DB : %v", dao.events[event.ID])

	return event
}

func (dao EventInMemoryDao) getAll() []Event {
	var values = make([]Event, len(dao.events))
	idx := 0
	for _, value := range dao.events {
		values[idx] = value
		idx++
	}

	return values
}

func (dao EventInMemoryDao) deleteById(id int64) Event {
	var deletedEvent = dao.events[id]
	delete(dao.events, id)
	return deletedEvent
}