// Code credit - http://gorm.io/docs/index.html
// Code credit - https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b
package dao

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

type EventMySqlOrmDao struct {
	Db *gorm.DB
}

func (dao EventMySqlOrmDao) InitDb() {
    // I like table names to be singulars
    // http://gorm.io/docs/conventions.html#Pluralized-Table-Name
	dao.Db.SingularTable(true)
}

func (dao EventMySqlOrmDao) FindById(id int64) Event {
	log.Println("Searching DB for event id ", id)
	var event Event
	dao.Db.First(&event, id)

	return event
}

func (dao EventMySqlOrmDao) Insert(event Event) {
	log.Println("Inserting event into DB")
	dao.Db.Create(&event)
}

func (dao EventMySqlOrmDao) Update(id int64, event Event) {
	log.Println("Updating DB for event id ", id)

	var existingEvent = dao.FindById(id)
	existingEvent.Title = event.Title
	existingEvent.Description = event.Description
	dao.Db.Save(existingEvent)

	log.Printf("Updated event in DB")
}

func (dao EventMySqlOrmDao) GetAll() []Event {
	log.Println("Searching all events in DB")

	var events []Event
	dao.Db.Find(&events)

	return events
}

func (dao EventMySqlOrmDao) DeleteById(id int64) {
	log.Println("Deleting from DB for event id ", id)
	dao.Db.Delete(dao.FindById(id))
	log.Printf("Deleted event from DB")
}
