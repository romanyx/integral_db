package set

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/romanyx/integral_db/internal/storage"
	"github.com/stretchr/testify/assert"
)

type decoderFunc func(*http.Request, *request) error

func (f decoderFunc) Decode(r *http.Request, m *request) error {
	return f(r, m)
}

type validaterFunc func(request) error

func (f validaterFunc) Validate(r request) error {
	return f(r)
}

type setterFunc func(context.Context, string, interface{})

func (f setterFunc) Set(ctx context.Context, key string, value interface{}) {
	f(ctx, key, value)
}

func Test_muxMap_Set(t *testing.T) {
	tests := []struct {
		name         string
		decodeFunc   func(*http.Request, *request) error
		validateFunc func(request) error
		setFunc      func(ctx context.Context, key string, value interface{})
		wantErr      bool
		expect       response
	}{
		{
			name: "decoder error",
			decodeFunc: func(*http.Request, *request) error {
				return errors.New("mock error")
			},
			wantErr: true,
		},
		{
			name: "validater error",
			decodeFunc: func(*http.Request, *request) error {
				return nil
			},
			validateFunc: func(request) error {
				return errors.New("mock error")
			},
			wantErr: true,
		},
		{
			name: "ok",
			decodeFunc: func(*http.Request, *request) error {
				return nil
			},
			validateFunc: func(request) error {
				return nil
			},
			setFunc: func(ctx context.Context, key string, value interface{}) {},
			expect: response{
				Message: "key set",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := muxMap{
				decoder:   decoderFunc(tt.decodeFunc),
				validater: validaterFunc(tt.validateFunc),
				setter:    setterFunc(tt.setFunc),
			}

			var got response
			req := httptest.NewRequest("POST", "http://any", nil)
			err := s.Set(req, &got)

			if tt.wantErr {
				assert.Error(t, err)
			}

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.expect, got)
			}
		})
	}
}

func Test_ozzoValidater_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     request
		wantErr bool
		expect  validationErrorResponse
	}{
		{
			name: "valid",
			req: request{
				Key:   "key",
				Value: 0,
			},
		},
		{
			name: "invalid key",
			req: request{
				Key:   "",
				Value: 0,
			},
			wantErr: true,
			expect: validationErrorResponse{
				Message: "you have validation errors",
				Errors: []validationError{
					validationError{
						Field:   "key",
						Message: "cannot be blank",
					},
				},
			},
		},
	}

	validater := ozzoValidater{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validater.Validate(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				got, ok := err.(validationErrorResponse)
				assert.True(t, ok)
				assert.Equal(t, tt.expect, got)
			}

			if !tt.wantErr {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_sSetter_Set(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value interface{}
	}{
		{
			name:  "ok",
			key:   "key",
			value: 1,
		},
	}

	getter := &sSetter{
		storage: storage.New(),
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			getter.Set(context.Background(), tt.key, tt.value)
		})
	}
}
