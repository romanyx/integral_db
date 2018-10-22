package set

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/romanyx/integral_db/internal/storage"
)

const (
	keySetMessage                  = "key set"
	validationErrorResponseMessage = "you have validation errors"
)

// Setter service for set key requests.
type Setter interface {
	Set(r *http.Request, resp *response) error
}

// NewService returns initialized service.
func NewService(storage storage.Storage, keyLiveTime time.Duration) Setter {
	srv := muxMap{
		decoder:   jsonDecoder{},
		validater: ozzoValidater{},
		setter: &sSetter{
			storage: storage,
		},
		keyLiveTime: keyLiveTime,
	}

	return &srv
}

type muxMap struct {
	decoder
	validater
	setter
	keyLiveTime time.Duration
}

type request struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type decoder interface {
	Decode(*http.Request, *request) error
}

type validater interface {
	Validate(request) error
}

type setter interface {
	Set(ctx context.Context, key string, value interface{})
}

func (s muxMap) Set(r *http.Request, resp *response) error {
	var req request

	if err := s.decoder.Decode(r, &req); err != nil {
		return errors.Wrap(err, "decode failed")
	}

	if err := s.validater.Validate(req); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	ctx, _ := context.WithTimeout(context.Background(), s.keyLiveTime)
	s.setter.Set(ctx, req.Key, req.Value)

	resp.Message = keySetMessage

	return nil
}

type jsonDecoder struct{}

func (d jsonDecoder) Decode(r *http.Request, req *request) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.Wrap(err, "unable to decode")
	}
	return nil
}

type ozzoValidater struct{}

func (v ozzoValidater) Validate(r request) error {
	validatationError := validationErrorResponse{Message: validationErrorResponseMessage}

	if err := validation.Validate(r.Key, validation.Required); err != nil {
		validatationError.Errors = append(validatationError.Errors,
			validationError{Field: "key", Message: err.Error()},
		)
	}

	if len(validatationError.Errors) > 0 {
		return validatationError
	}

	return nil
}

type sSetter struct {
	storage storage.Storage
}

func (s *sSetter) Set(ctx context.Context, key string, value interface{}) {
	s.storage.Set(ctx, key, value)
}

type notFoundResponse struct {
	Message string `json:"message"`
}

// Error implements the error interface.
func (r notFoundResponse) Error() string {
	return r.Message
}

type validationErrorResponse struct {
	Message string            `json:"message"`
	Errors  []validationError `json:"errors"`
}

type validationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (r validationErrorResponse) Error() string {
	return r.Message
}
