package apitool

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Ronmi/rtoolkit/jsonapi"
	"github.com/Ronmi/rtoolkit/session"
)

// Session creates a api middleware that handles session related functions
//
// If you are facing "Trailer Header" with original session middleware,
// this should be helpful.
func Session(m *session.Manager) jsonapi.Middleware {
	return func(h jsonapi.Handler) jsonapi.Handler {
		return func(
			d *json.Decoder,
			r *http.Request,
			w http.ResponseWriter,
		) (interface{}, error) {
			req := r
			sess, err := m.Start(w, r)
			if err != nil {
				return nil, jsonapi.E500.SetOrigin(err)
			}

			req = r.WithContext(context.WithValue(
				r.Context(),
				session.SessionObjectKey,
				sess,
			))
			return h(d, req, w)
		}
	}
}
