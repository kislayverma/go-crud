package event

import (
	"log"
	"net/http"
	"strconv"
)

func getEvent(id int64, responseChannel chan bool) {
	var url string = "http://localhost:8080/event/" + strconv.FormatInt(id, 10)
	log.Print("Invoking url ", url)
	resp, err := http.Get(url)
	if err != nil {
		responseChannel <- false
	}
	if resp.StatusCode == 200 {
		responseChannel <- true
	} else {
		responseChannel <- false
	}
	resp.Body.Close()
}
