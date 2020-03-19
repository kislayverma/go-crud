// Code credit - http://gorm.io/docs/index.html
// Code credit - https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b
package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

type EventMySqlOrmDao struct {
	db *gorm.DB
}

func (dao EventMySqlOrmDao) initDb() {
    // I like table names to be singulars
    // http://gorm.io/docs/conventions.html#Pluralized-Table-Name
	dao.db.SingularTable(true)
}

func (dao EventMySqlOrmDao) findById(id int64) Event {
	log.Println("Searching DB for event id ", id)
	var event Event
	dao.db.First(&event, id)

	return event
}

func (dao EventMySqlOrmDao) insert(event Event) {
	log.Println("Inserting event into DB")
	dao.db.Create(&event)
}

func (dao EventMySqlOrmDao) update(id int64, event Event) {
	log.Println("Updating DB for event id ", id)

	var existingEvent = dao.findById(id)
	existingEvent.Title = event.Title
	existingEvent.Description = event.Description
	dao.db.Save(existingEvent)

	log.Printf("Updated event in DB")
}

func (dao EventMySqlOrmDao) getAll() []Event {
	log.Println("Searching all events in DB")

	var events []Event
	dao.db.Find(&events)

	return events
}

func (dao EventMySqlOrmDao) deleteById(id int64) {
	log.Println("Deleting from DB for event id ", id)
	dao.db.Delete(dao.findById(id))
	log.Printf("Deleted event from DB")
}
