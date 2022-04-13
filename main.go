package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var logger log.Logger

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":4000", nil)
}

func hello(w http.ResponseWriter, _ *http.Request) {
	// send an invalid query param error without having to care about formatting, codes, messages...
	RespondError(w, ErrorInvalidQueryParam, nil)
}

// RespondError handles all API error responses
func RespondError(w http.ResponseWriter, e AppError, err error) {
	if err != nil {
		logger.Println(err)
	}

	jsonData, err := json.Marshal(e)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode())
	_, err = w.Write(jsonData)
	if err != nil {
		logger.Println(err)
	}
}
