package utils

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.WriteHeader(statusCode)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"error": "unknown error"})
	}
}

func Response(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
