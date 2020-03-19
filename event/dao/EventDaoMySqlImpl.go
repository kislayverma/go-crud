// Code credit - https://medium.com/@hugo.bjarred/mysql-and-golang-ea0d620574d2
// Code credit - https://medium.com/@hugo.bjarred/rest-api-with-golang-mux-mysql-c5915347fa5b
package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type EventMySqlDao struct {
	db *sql.DB
}

func (dao EventMySqlDao) InitDb() {
	// Do nothing
}

func (dao EventMySqlDao) FindById(id int64) Event {
	log.Println("Searching DB for event id ", id)
	data, err := dao.db.Query("SELECT id, title, description FROM event WHERE id = ?", id)
	if err != nil {
		log.Panic("Error running query : %v", err.Error())
	}

	var event Event
	for data.Next() {
		err := data.Scan(&event.ID, &event.Title, &event.Description)
		if err != nil {
			log.Panic(err.Error())
		}
	}

	return event
}

func (dao EventMySqlDao) Insert(event Event) {
	stmt, err := dao.db.Prepare("INSERT INTO event (title,description) VALUES (?, ?)")
	if err != nil {
		log.Panic("Error preparing query : %v", err.Error())
	}
	_, runErr := stmt.Exec(event.Title, event.Description)
	if runErr != nil {
		log.Panic("Error running query : %v", err.Error())
	}

	log.Printf("Inserted event in DB")
}

func (dao EventMySqlDao) Update(id int64, event Event) {
	stmt, err := dao.db.Prepare("UPDATE event SET title = ?, description=? WHERE id = ?")
	if err != nil {
		log.Panic("Error preparing query : %v", err.Error())
	}
	_, runErr := stmt.Exec(event.Title, event.Description, id)
	if runErr != nil {
		log.Panic("Error running query : %v", err.Error())
	}

	log.Printf("Updated event in DB")
}

func (dao EventMySqlDao) GetAll() []Event {
	data, err := dao.db.Query("SELECT id, title, description FROM event")
	if err != nil {
		log.Panic("Error running query : %v", err.Error())
	}

	var events []Event
	for data.Next() {
		var event Event
		err := data.Scan(&event.ID, &event.Title, &event.Description)
		if err != nil {
			panic(err.Error())
		}
		events = append(events, event)
	}

	return events
}

func (dao EventMySqlDao) DeleteById(id int64) {
	stmt, err := dao.db.Prepare("DELETE FROM event where id = ?")
	if err != nil {
		log.Panic("Error preparing query : %v", err.Error())
	}
	_, runErr := stmt.Exec(id)
	if runErr != nil {
		log.Panic("Error running query : %v", err.Error())
	}

	log.Printf("Deleted event from DB")
}
