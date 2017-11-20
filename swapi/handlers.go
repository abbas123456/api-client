package swapi

import (
	"encoding/json"
	"errors"
	"net/http"

	h "github.com/kynrai/api-client/common/http"
)

var ErrMissingID = errors.New("missing ID")

func GetPerson(p PersonGetter) h.ErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// Can also use gorilla mux etc to get /***/1 style APIs
		id := r.URL.Query().Get("id")
		if id == "" {
			return h.HTTPError{Code: http.StatusBadRequest, Err: ErrMissingID}
		}
		person, err := p.GetPerson(r.Context(), id)
		if err != nil {
			return h.HTTPError{Code: http.StatusInternalServerError, Err: err}
		}
		if err := json.NewEncoder(w).Encode(person); err != nil {
			return h.HTTPError{Code: http.StatusInternalServerError, Err: err}
		}
		return nil
	}
}
