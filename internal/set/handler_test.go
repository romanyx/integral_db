package set

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
		setFunc func(*http.Request, *response) error
		code    int
	}{
		{
			name: "ok",
			setFunc: func(*http.Request, *response) error {
				return nil
			},
			code: http.StatusOK,
		},
		{
			name: "validation error",
			setFunc: func(*http.Request, *response) error {
				return validationErrorResponse{}
			},
			code: http.StatusBadRequest,
		},
		{
			name: "unexpected error",
			setFunc: func(*http.Request, *response) error {
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

			h := NewHandler(SetterFunc(tt.setFunc))
			h(res, req)

			assert.Equal(t, tt.code, res.Code)
		})
	}
}

type SetterFunc func(*http.Request, *response) error

func (f SetterFunc) Set(r *http.Request, resp *response) error {
	return f(r, resp)
}
