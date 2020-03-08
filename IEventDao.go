package main

type IEventDao interface {
	initDb()
	findById(id int64) Event
	insert(event Event)
	update(id int64, event Event)
	getAll() []Event
	deleteById(id int64)
}
