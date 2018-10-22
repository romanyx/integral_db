package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/romanyx/integral_db/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

func Test_PostSet(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		code   int
		schema string
	}{
		{
			name:   "ok",
			body:   `{"key": "key", "value": "value"}`,
			code:   http.StatusOK,
			schema: `{"type":"object", "required": ["message"], "properties": {"message": {"type": "string"}}}`,
		},
		{
			name:   "validaion errors",
			schema: `{"type":"object", "required": ["message", "errors"], "properties": {"message": {"type": "string"}, "errors": {"type": "array", "items": {"type": "object", "required": ["field", "message"], "properties": {"field": {"type": "string"}, "message": {"type": "string"}}}}}}`,
			code:   http.StatusBadRequest,
		},
	}

	s := storage.New()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := httpMux(s, time.Second)
			s := httptest.NewServer(handler)
			defer s.Close()

			req := httptest.NewRequest("POST", fmt.Sprintf("%s/set", s.URL), strings.NewReader(tt.body))
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			assert.Equal(t, tt.code, res.Code)

			schema := gojsonschema.NewStringLoader(tt.schema)
			doc := gojsonschema.NewStringLoader(res.Body.String())

			result, err := gojsonschema.Validate(schema, doc)

			assert.Nil(t, err)
			assert.True(t, result.Valid())
			assert.Empty(t, result.Errors())
		})
	}
}
