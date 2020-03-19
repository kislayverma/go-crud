package event

import (
	"github.com/kislayverma/go-crud/http-utils"
	"log"
)

func logResponse(reqCtx http_utils.RequestContext, response string) {
	log.Print(reqCtx.RequestId, " : Response in ResponseLogger ", response)
}
