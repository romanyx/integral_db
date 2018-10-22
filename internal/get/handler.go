package get

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/romanyx/integral_db/internal/responses"
)

// NewHandler returns handler for get key requests.
func NewHandler(srv Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp response

		if err := srv.Get(r, &resp); err != nil {
			switch resp := errors.Cause(err).(type) {
			case validationErrorResponse:
				responses.BadRequest(w, resp)
			case notFoundResponse:
				responses.NotFound(w, resp)
			default:
				responses.InternalServerError(w)
			}
			return
		}

		responses.OK(w, resp)
	}
}

type response struct {
	Message string `json:"message"`
	Data    data   `json:"data"`
}

type data struct {
	Value interface{} `json:"value"`
}
