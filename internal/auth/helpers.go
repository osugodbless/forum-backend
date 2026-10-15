package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (s *Service) writeJSON(w http.ResponseWriter, status int, data any, headers map[string]string) error {
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

func (s *Service) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
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

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	return nil
}

// func (s *Service) serverError(w http.ResponseWriter, r *http.Request, status int, err error) {
// 	var (
// 		method = r.Method
// 		url    = r.URL.RequestURI()
// 	)

// 	s.logger.Error(err.Error(), "method", method, "url", url)
// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// }
