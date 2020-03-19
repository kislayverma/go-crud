// Code credits : https://www.gorillatoolkit.org/pkg/context
// Code credits : https://godoc.org/github.com/gorilla/mux
package http_utils

import (
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"log"
	"net/http"
)

const REQUEST_ID_HEADER_NAME string = "X-Req-Id"

type RequestContext struct {
	RequestId string
}

func GetRequestContext(r *http.Request, key string) RequestContext {
	if rv := context.Get(r, key); rv != nil {
		return rv.(RequestContext)
	} else {
		setRequestContext(r, uuid.New().String())
	}
	return RequestContext{""}
}

func CorrelationIdSettingMw (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Header[REQUEST_ID_HEADER_NAME]
		if !ok {
			newReqId := uuid.New().String()
			setRequestContext(r, newReqId)
			r.Header[REQUEST_ID_HEADER_NAME] = []string {newReqId}
		}

		next.ServeHTTP(w, r)
		context.Clear(r)
	})
}

func RequestLoggingMw(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    	// Log request
    	log.Printf("\n================Incoming request===============\n%s\n%s\n===============================================",
    		r.RequestURI, r.Header)

    	// Call the next handler
		next.ServeHTTP(w, r)

		// Log response headers
		log.Printf("\n================Outgoing Response===============\n%s\n================================================",
			w.Header())
	})
}

func setRequestContext(r *http.Request, requestId string) {
	reqCtx := RequestContext{requestId}
	context.Set(r, requestId, reqCtx)
}