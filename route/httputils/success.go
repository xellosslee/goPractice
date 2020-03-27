package httputils

import (
	"encoding/json"
	"net/http"
)

// ResponseJSON 공통 http 통신 성공 처리 함수. result 객체는 반드시 json 객체여야 한다.
func ResponseJSON(w *http.ResponseWriter, result interface{}) {
	writer := *w

	marshalled, err := json.Marshal(result)

	if err != nil {
		ResponseError(w, 500, "Internal Server Error", "Error marshalling response JSON", err)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(marshalled)
}
