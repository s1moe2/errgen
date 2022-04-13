//go:generate go run errorgen.go
package main

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	message    string
	code       string
	statusCode int
}

func (e AppError) StatusCode() int {
	return e.statusCode
}

func (a AppError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}{
		Message: a.message,
		Code:    a.code,
	})
}

var ErrorInternal = AppError{
	message:    "Unexpected server error",
	code:       "internal_error",
	statusCode: http.StatusInternalServerError,
}

var ErrorPayloadSize = AppError{
	message:    "Request body payload exceeds the maximum of bytes allowed",
	code:       "payload_size",
	statusCode: http.StatusRequestEntityTooLarge,
}

var ErrorPayloadParse = AppError{
	message:    "Request body payload could not be parsed",
	code:       "payload_parse",
	statusCode: http.StatusUnprocessableEntity,
}

var ErrorRouteNotFound = AppError{
	message:    "The request route does not exist",
	code:       "route_not_found",
	statusCode: http.StatusNotFound,
}

var ErrorMethodNotAllowed = AppError{
	message:    "The request HTTP method is not allowed in this server",
	code:       "method_not_allowed",
	statusCode: http.StatusMethodNotAllowed,
}

var ErrorInvalidQueryParam = AppError{
	message:    "One of the query parameters is not valid",
	code:       "invalid_query_param",
	statusCode: http.StatusBadRequest,
}

var ErrorInvalidPayload = AppError{
	message:    "Request payload could not be parsed correctly",
	code:       "invalid_payload",
	statusCode: http.StatusBadRequest,
}
