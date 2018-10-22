package get

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_NewHandler(t *testing.T) {
	tests := []struct {
		name    string
		getFunc func(*http.Request, *response) error
		code    int
	}{
		{
			name: "ok",
			getFunc: func(*http.Request, *response) error {
				return nil
			},
			code: http.StatusOK,
		},
		{
			name: "validation error",
			getFunc: func(*http.Request, *response) error {
				return validationErrorResponse{}
			},
			code: http.StatusBadRequest,
		},
		{
			name: "not found error",
			getFunc: func(*http.Request, *response) error {
				return notFoundResponse{}
			},
			code: http.StatusNotFound,
		},
		{
			name: "unexpected error",
			getFunc: func(*http.Request, *response) error {
				return errors.New("mock error")
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("POST", "http://any-host/auth", nil)
			res := httptest.NewRecorder()

			h := NewHandler(GetterFunc(tt.getFunc))
			h(res, req)

			assert.Equal(t, tt.code, res.Code)
		})
	}
}

type GetterFunc func(*http.Request, *response) error

func (f GetterFunc) Get(r *http.Request, resp *response) error {
	return f(r, resp)
}

func Test_NewService(t *testing.T) {}
