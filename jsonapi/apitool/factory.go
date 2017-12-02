package apitool

import (
	"encoding/json"
	"net/http"

	"github.com/Ronmi/rtoolkit/jsonapi"
)

// LogProvider defines what info a jsonapi logger can use
type LogProvider func(r *http.Request, data interface{}, err error)

// LogIn wraps handler, uses LogProvider p for logging purpose
//
//      jsonapi.With(
//          apitools.LogIn(apitool.JSONFormat(myLogger))
//      ).RegisterAll(myHandlerClass)
func LogIn(p LogProvider) jsonapi.Middleware {
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

// LogErrIn wraps handler, uses LogProvider p for logging purpose, but only for errors
//
//      jsonapi.With(
//          apitools.LogErrIn(apitool.JSONFormat(myLogger))
//      ).RegisterAll(myHandlerClass)
func LogErrIn(p LogProvider) jsonapi.Middleware {
	return jsonapi.Middleware(func(h jsonapi.Handler) jsonapi.Handler {
		return jsonapi.Handler(func(
			d *json.Decoder,
			r *http.Request,
			w http.ResponseWriter,
		) (interface{}, error) {
			data, err := h(d, r, w)
			if err != nil {
				p(r, data, err)
			}

			return data, err
		})
	})
}
