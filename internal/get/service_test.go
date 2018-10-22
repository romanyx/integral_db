package get

import (
	"context"
	"net/http"
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

type getterFunc func(string) (interface{}, error)

func (f getterFunc) Get(key string) (value interface{}, err error) {
	return f(key)
}

func Test_muxMap_Get(t *testing.T) {
	tests := []struct {
		name         string
		decodeFunc   func(*http.Request, *request) error
		validateFunc func(request) error
		getFunc      func(string) (interface{}, error)
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
			name: "getter error",
			decodeFunc: func(*http.Request, *request) error {
				return nil
			},
			validateFunc: func(request) error {
				return nil
			},
			getFunc: func(key string) (value interface{}, err error) {
				return nil, errors.New("mock error")
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
			getFunc: func(key string) (value interface{}, err error) {
				return 0, nil
			},
			expect: response{
				Message: "key found",
				Data: data{
					Value: 0,
				},
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
				getter:    getterFunc(tt.getFunc),
			}

			var got response
			err := s.Get(nil, &got)

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
				Key: "key",
			},
		},
		{
			name: "invalid key",
			req: request{
				Key: "",
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

func Test_sGetter_Get(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
		err     notFoundResponse
		expect  interface{}
	}{
		{
			name:   "ok",
			key:    "key",
			expect: 0,
		},
		{
			name:    "key not found",
			key:     "not found",
			wantErr: true,
			err: notFoundResponse{
				Message: "key not found",
			},
			expect: nil,
		},
	}

	getter := &sGetter{
		storage: storage.New(),
	}

	getter.storage.Set(context.Background(), "key", 0)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := getter.Get(tt.key)

			if tt.wantErr {
				assert.Error(t, err)
				got, ok := err.(notFoundResponse)
				assert.True(t, ok)
				assert.Equal(t, tt.err, got)
			}

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.expect, got)
			}
		})
	}
}
