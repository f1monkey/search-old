package node

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/invopop/validation"
)

type errorResponse struct {
	Message string               `json:"message"`
	Errors  []errorResponseError `json:"errors"`
}

type errorResponseError struct {
	Field string `json:"field"`
	Error string `json:"error"`
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

func writeError(w http.ResponseWriter, statusCode int, msg string, errors []errorResponseError) {
	data, _ := json.Marshal(errorResponse{Message: msg, Errors: errors})

	setContentType(w)
	w.WriteHeader(statusCode)
	w.Write(data)
}

func writeValidationErr(w http.ResponseWriter, errorList validation.Errors) {
	responseErrors := make([]errorResponseError, 0, len(errorList))
	for field, errors := range errorList {
		responseErrors = append(responseErrors, errorResponseError{Field: field, Error: errors.Error()})
	}

	writeError(w, http.StatusUnprocessableEntity, "Validation error", responseErrors)
}

func handleErr(w http.ResponseWriter, err error) {
	if e, ok := err.(*requestValidationErr); ok {
		writeValidationErr(w, e.errors)
		return
	}

	// @todo handle invalid json request errors

	log.Println(err) // @todo

	writeSimpleError(
		w,
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError),
	)
}

type requestValidationErr struct {
	errors validation.Errors
}

func (e *requestValidationErr) Error() string {
	return e.errors.Error()
}

func newRequestValidationErr(errors validation.Errors) *requestValidationErr {
	return &requestValidationErr{
		errors: errors,
	}
}
