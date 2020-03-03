package httphandlers

import (
	"net/http"

	"cndf.order.was/httphandlers/httputils"
	"cndf.order.was/storage"
)

func List(w http.ResponseWriter, r *http.Request) {
	// 데이터의 List를 JSON형태로 리턴
	httputils.ResponseJSON(&w, storage.Get())
}
