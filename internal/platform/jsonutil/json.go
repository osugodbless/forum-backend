package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// envelopeJSON envelopes response JSON data in a parent string
type EnvelopeJSON map[string]any

type Responder struct {
	logger *slog.Logger
}

func NewResponder(logger *slog.Logger) *Responder {
	return &Responder{logger: logger}
}

func (res *Responder) WriteJSON(w http.ResponseWriter, status int, data any, headers map[string]string) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	payload = append(payload, '\n')

	for key, value := range headers {
		w.Header().Set(key, value)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(payload)

	return nil
}

func (res *Responder) ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// maxBytes is the max payload size allowed
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize a new json.Decoder for decoding JSON values
	// json.Decoder is only allowed to read up to a maximum of 1MB
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode request body to dst(destination)
	err := dec.Decode(dst)
	if err != nil {
		var (
			syntaxError           *json.SyntaxError
			unmarshalTypeError    *json.UnmarshalTypeError
			invalidUnmarshalError *json.InvalidUnmarshalError
			maxBytesError         *http.MaxBytesError
		)

		switch {

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field '%s'", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (res *Responder) logError(r *http.Request, err error) {
	res.logger.Error("http request failed", "error", err.Error(), "method", r.Method, "url", r.URL.RequestURI())
}

func (res *Responder) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	err := res.WriteJSON(w, status, EnvelopeJSON{"error": message}, nil)
	if err != nil {
		res.logError(r, err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

func (res *Responder) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	res.logError(r, err)
	res.ErrorResponse(w, r, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func (res *Responder) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	res.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (res *Responder) NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	res.ErrorResponse(w, r, http.StatusNotFound, err.Error())
}

// func (res *Responder) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
// 	res.ErrorResponse(w, r, http.StatusMethodNotAllowed, fmt.Sprintf("the %s method is not allowed for the requested resource", r.Method))
// }

func (res *Responder) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	res.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (res *Responder) RequestConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	res.ErrorResponse(w, r, http.StatusConflict, err.Error())
}
