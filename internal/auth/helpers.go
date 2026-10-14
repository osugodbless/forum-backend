package auth

import (
	"encoding/json"
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

// func (s *Service) serverError(w http.ResponseWriter, r *http.Request, status int, err error) {
// 	var (
// 		method = r.Method
// 		url    = r.URL.RequestURI()
// 	)

// 	s.logger.Error(err.Error(), "method", method, "url", url)
// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// }
