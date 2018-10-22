package set

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/romanyx/integral_db/internal/responses"
)

// NewHandler returns handler for set key requests.
func NewHandler(srv Setter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp response

		if err := srv.Set(r, &resp); err != nil {
			switch resp := errors.Cause(err).(type) {
			case validationErrorResponse:
				responses.BadRequest(w, resp)
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
}
