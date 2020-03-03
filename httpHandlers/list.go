package httpHandlers

import (
	"net/http"

	"cndf.order.was/httpHandlers/httpUtils"
	"cndf.order.was/storage"
)

func List(w http.ResponseWriter, r *http.Request) {
	httpUtils.HandleSuccess(&w, storage.Get())
}
