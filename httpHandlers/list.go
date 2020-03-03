package httphandlers

import (
	"net/http"

	"cndf.order.was/httphandlers/httputils"
	"cndf.order.was/storage"
)

func List(w http.ResponseWriter, r *http.Request) {
	httputils.ResponseJSON(&w, storage.Get())
}
