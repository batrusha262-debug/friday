package reply_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.appkode.ru/pub/go/failure"
	"github.com/stretchr/testify/assert"

	"friday/pkg/httpx/reply"
)

func TestError_statusCodes(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		status int
	}{
		{"invalid argument", failure.NewInvalidArgumentError("x"), http.StatusBadRequest},
		{"not found", failure.NewNotFoundError("x"), http.StatusNotFound},
		{"conflict", failure.NewConflictError("x"), http.StatusConflict},
		{"unauthorized", failure.NewUnauthorizedError("x"), http.StatusUnauthorized},
		{"forbidden", failure.NewForbiddenError("x"), http.StatusForbidden},
		{"unprocessable entity", failure.NewUnprocessableEntityError("x"), http.StatusUnprocessableEntity},
		{"internal server error", failure.NewInternalServerError("x"), http.StatusInternalServerError},
		{"unknown error", errors.New("something went wrong"), http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			reply.Error(context.Background(), w, tc.err)

			assert.Equal(t, tc.status, w.Code)
		})
	}
}
