package httphandlers

import (
	"log"
	"net/http"

	"cndf.order.was/httphandlers/httputils"
)

// HandleRequest 공용 핸들러
// Get, Post, Delete 메소드를 수신받는다.
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming Request:", r.Method)
	switch r.Method {
	case http.MethodGet:
		List(w, r)
		break
	case http.MethodPost:
		Add(w, r)
		break
	case http.MethodDelete:
		Remove(w, r)
		break
	default:
		httputils.ResponseError(&w, 405, "Method not allowed", "Method not allowed", nil)
		break
	}
}
