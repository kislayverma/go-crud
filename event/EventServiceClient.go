package event

import (
	"github.com/kislayverma/go-crud/http-utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getEvent(reqCtx http_utils.RequestContext, id int64, responseChannel chan bool) {
	var url string = "http://localhost:8080/event/" + strconv.FormatInt(id, 10)
	log.Print("Invoking url ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Unable to create request", err.Error())
	}
	req.Header.Set(http_utils.REQUEST_ID_HEADER_NAME, reqCtx.RequestId)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{Timeout: time.Second * 1}
	resp, err := client.Do(req)
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
