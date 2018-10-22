package get

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/romanyx/integral_db/internal/storage"
)

const (
	keyFoundMessage                = "key found"
	validationErrorResponseMessage = "you have validation errors"
	notFoundMessage                = "key not found"
)

// Getter service for get key requests.
type Getter interface {
	Get(r *http.Request, resp *response) error
}

// NewService returns initialized service.
func NewService(storage storage.Storage) Getter {
	srv := muxMap{
		decoder:   jsonDecoder{},
		validater: ozzoValidater{},
		getter: &sGetter{
			storage: storage,
		},
	}

	return &srv
}

type muxMap struct {
	decoder
	validater
	getter
}

type request struct {
	Key string `json:"key"`
}

type decoder interface {
	Decode(*http.Request, *request) error
}

type validater interface {
	Validate(request) error
}

type getter interface {
	Get(key string) (value interface{}, err error)
}

func (s muxMap) Get(r *http.Request, resp *response) error {
	var req request

	if err := s.decoder.Decode(r, &req); err != nil {
		return errors.Wrap(err, "decode failed")
	}

	if err := s.validater.Validate(req); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	value, err := s.getter.Get(req.Key)
	if err != nil {
		return errors.Wrap(err, "get value failed")
	}

	resp.Message = keyFoundMessage
	resp.Data = data{
		Value: value,
	}

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

type sGetter struct {
	storage storage.Storage
}

func (g *sGetter) Get(key string) (interface{}, error) {
	value, err := g.storage.Get(key)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, notFoundResponse{
				Message: notFoundMessage,
			}
		}
	}

	return value, nil
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
