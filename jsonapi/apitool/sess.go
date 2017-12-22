package apitool

import (
	"context"

	"github.com/Ronmi/rtoolkit/jsonapi"
	"github.com/Ronmi/rtoolkit/session"
)

// Session creates a api middleware that handles session related functions
//
// If you are facing "Trailer Header" problem with original session middleware,
// this should be helpful.
//
//     jsonapi.With(
//         apitool.Session(mySessMgr),
//     ).RegisterAll(myHandlerClass)
func Session(m *session.Manager) jsonapi.Middleware {
	return func(h jsonapi.Handler) jsonapi.Handler {
		return func(req jsonapi.Request) (interface{}, error) {
			r := req.R()
			sess, err := m.Start(req.W(), r)
			if err != nil {
				return nil, jsonapi.E500.SetOrigin(err)
			}

			r = r.WithContext(context.WithValue(
				r.Context(),
				session.SessionObjectKey,
				sess,
			))
			return h(jsonapi.WrapRequest(req, r))
		}
	}
}
