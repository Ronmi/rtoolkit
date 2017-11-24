package apilog

import (
	"encoding/json"
	"net/http"

	"github.com/Ronmi/rtoolkit/jsonapi"
)

// LogProvider defines what info a jsonapi logger can use
type LogProvider func(r *http.Request, data interface{}, err error)

// For wraps handler for logging purpose
func Use(p LogProvider) jsonapi.Middleware {
	return jsonapi.Middleware(func(h jsonapi.Handler) jsonapi.Handler {
		return jsonapi.Handler(func(
			d *json.Decoder,
			r *http.Request,
			w http.ResponseWriter,
		) (interface{}, error) {
			data, err := h(d, r, w)
			p(r, data, err)

			return data, err
		})
	})
}
