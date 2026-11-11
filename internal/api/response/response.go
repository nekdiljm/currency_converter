package response

import (
	"currency-converter/pkg/serializer"
	"net/http"
)

func JSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := serializer.Encode(w, v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
