package httpx

import (
	"net/http"

	"friday/pkg/httpx/reply"
)

func Handler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			reply.Error(r.Context(), w, err)
		}
	}
}
