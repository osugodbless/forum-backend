package auth

import (
	"net/http"
)

func (s *Service) serverError(w http.ResponseWriter, r *http.Request, status int, err error) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
	)

	s.logger.Error(err.Error(), "method", method, "url", url)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
