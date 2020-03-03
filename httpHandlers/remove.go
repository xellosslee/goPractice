package httphandlers

import (
	"io/ioutil"
	"net/http"

	"encoding/json"

	"cndf.order.was/httphandlers/httputils"
	"cndf.order.was/model"
	"cndf.order.was/storage"
)

func Remove(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		httputils.ResponseError(&w, 500, "Internal Server Error", "Error reading data from body", err)
		return
	}

	var users model.Users
	// http를 통해 body의 문자열을 json으로 Convert 시도
	err = json.Unmarshal(requestBody, &users)

	if err != nil {
		httputils.ResponseError(&w, 400, "Bad Request", "Error unmarshalling", err)
		return
	}

	if users.ID == 0 {
		httputils.ResponseError(&w, 500, "Bad Request", "ID not provided", nil)
		return
	}

	// Remove후 결과를 리턴
	if storage.Remove(users.ID) {
		httputils.ResponseJSON(&w, model.Users{ID: users.ID})
	} else {
		httputils.ResponseError(&w, 400, "Bad Request", "ID not found", nil)
	}
}
