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
//
// Created middleware will try to save update cookie ttl value if possible. It
// fails silently.
func Session(m *session.Manager) jsonapi.Middleware {
	return func(h jsonapi.Handler) jsonapi.Handler {
		return func(req jsonapi.Request) (i interface{}, e error) {
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
			i, e = h(jsonapi.WrapRequest(req, r))

			_ = sess.Save(req.W())
			return
		}
	}
}
