package main

type IEventDao interface {
	initDb()
	findById(id int64) Event
	insert(event Event) Event
	update(id int64, event Event) Event
	getAll() []Event
	deleteById(id int64) Event
}
