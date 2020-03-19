// Go-CRUD Service
//
// Sample CRUD service built in Go, to learn Go
//
//     Schemes: http
//     Host: localhost:8080
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Kislay<kislay.nsit@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/kislayverma/go-crud/event"
	"github.com/kislayverma/go-crud/event/dao"
	"github.com/kislayverma/go-crud/http-utils"
	"log"
	"net/http"
)

var eventService event.EventService

// Swagger setup credit :
// https://medium.com/@supun.muthutantrige/lets-go-everything-you-need-to-know-about-creating-a-restful-api-in-go-part-iv-52666c5221d4
func main() {
	// Initialize the database access layer
	// Not the best : conn details leaking into main program - need to enable configs driven something something
	// TODO - Make this work
	// This is what I want, but the connection gets closed when I do it in the build*Dao methods below.
	// db = buildMySqlOrmDao()
	dbConn, err := gorm.Open("mysql", "root:@(localhost)/gocrud")
	if err != nil {
		log.Panic(err.Error())
		panic(err.Error())
	}

	db := dao.EventMySqlOrmDao{dbConn}
	db.InitDb()
	defer dbConn.Close()

	// Initialize the request router
	router := mux.NewRouter().StrictSlash(true)
	router.Use(http_utils.CorrelationIdSettingMw)
	router.Use(http_utils.RequestLoggingMw)

	// Set up top level routes
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/swagger.json", swagger).Methods("GET")
	// Create Event service and let it register its routes
	eventService = event.EventService{db, router}
	eventService.RegisterRoutes()

	// Launch the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	// Code Credit : https://medium.com/the-andela-way/build-a-restful-json-api-with-golang-85a83420c9da

	// swagger:operation GET / Hello Hello
	//
	// Returns a simple Hello message
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - text/plain
	// responses:
	//   '200':
	//     description: The hello message
	//     type: string
    fmt.Fprintf(w, "Welcome home!")
}

func swagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "swagger.json")
}
