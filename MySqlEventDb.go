package main

import (
	"log"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)

type EventMySqlDao struct {
	events map[int64]Event
	db sql.DB
	eventIdCount int64
}

func (dao EventMySqlDao) initDb() {
	db, err := sql.Open("mysql", "root:@/<database-name>")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
}

func (dao EventMySqlDao) findById(id int64) Event {
	log.Printf("returning event from DB : %v", dao.events[id])

	return dao.events[id]
}

func (dao EventMySqlDao) insert(event Event) Event {
	event.ID = dao.eventIdCount
	dao.eventIdCount++
	dao.events[event.ID] = event
	log.Printf("Inserted event in DB : %v", dao.events[event.ID])
	log.Printf("Record count : %d", len(dao.events));

	return event
}

func (dao EventMySqlDao) update(id int64, event Event) Event {
	dao.events[id] = event
	log.Printf("Updated event in DB : %v", dao.events[event.ID])

	return event
}

func (dao EventMySqlDao) getAll() []Event {
	var values = make([]Event, len(dao.events))
	idx := 0
	for _, value := range dao.events {
		values[idx] = value
		idx++
	}

	return values
}

func (dao EventMySqlDao) deleteById(id int64) Event {
	var deletedEvent = dao.events[id]
	delete(dao.events, id)
	return deletedEvent
}