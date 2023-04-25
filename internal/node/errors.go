package node

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func setContentType(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

func writeSimpleError(w http.ResponseWriter, statusCode int, msg string) {
	data, _ := json.Marshal(errorResponse{Message: msg})

	setContentType(w)
	w.WriteHeader(statusCode)
	w.Write(data)
}

func handleErr(w http.ResponseWriter, err error) {
	log.Println(err) // @todo

	w.WriteHeader(http.StatusInternalServerError) // @todo
}
