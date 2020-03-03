package httpHandlers

import (
	"net/http"

	"cndf.order/httpHandlers/httpUtils"
	"cndf.order/storage"
)

func List(w http.ResponseWriter, r *http.Request) {
	httpUtils.HandleSuccess(&w, storage.Get())
}
