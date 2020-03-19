package dao

import (
	"log"
)

type EventInMemoryDao struct {
	events map[int64] Event
	eventIdCount int64
}

func (dao EventInMemoryDao) InitDb() {
	// do nothing
}

func (dao EventInMemoryDao) FindById(id int64) Event {
	log.Printf("returning event from DB : %v", dao.events[id])
	return dao.events[id]
}

func (dao EventInMemoryDao) Insert(event Event) {
	var newEventForDb = Event{int64(len(dao.events)) + 1, event.Title, event.Description}
	dao.events[newEventForDb.ID] = newEventForDb
	log.Printf("Inserted event in DB : %v", dao.events[newEventForDb.ID])
	log.Printf("Record count : %d", len(dao.events));
}

func (dao EventInMemoryDao) Update(id int64, event Event) {
	dao.events[id] = event
	log.Printf("Updated event in DB : %v", dao.events[event.ID])
}

func (dao EventInMemoryDao) GetAll() []Event {
	var values = make([]Event, len(dao.events))
	idx := 0
	for _, value := range dao.events {
		values[idx] = value
		idx++
	}

	return values
}

func (dao EventInMemoryDao) DeleteById(id int64) {
	delete(dao.events, id)
}