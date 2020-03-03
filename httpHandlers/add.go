package httphandlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"log"

	"cndf.order.was/httphandlers/httputils"
	"cndf.order.was/model"
	"cndf.order.was/storage"
)

// Add 함수 model 객체에 insert 하는 용도
// model.Message 객체로 된 json 을 받아 model.Message-list 객체에 insert 한다
// {
//   "sender":"a",
//   "message":"m"
// }
func Add(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	byteData, err := ioutil.ReadAll(r.Body)

	if err != nil {
		httputils.ResponseError(&w, 500, "Internal Server Error", "Error reading data from body", err)
		return
	}

	var message model.Message

	err = json.Unmarshal(byteData, &message)

	if err != nil {
		httputils.ResponseError(&w, 500, "Internal Server Error", "Error unmarhsalling JSON", err)
		return
	}

	if message.Message == "" || message.Sender == "" {
		httputils.ResponseError(&w, 400, "Bad Request", "Unmarshalled JSON didn't have required fields", nil)
		return
	}

	id := storage.Add(message)

	log.Println("Added message:", message)

	httputils.ResponseJSON(&w, model.Users{ID: id})
}
