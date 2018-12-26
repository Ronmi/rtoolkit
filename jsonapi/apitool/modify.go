package apitool

import (
	"net/http"
	"time"

	"github.com/Ronmi/rtoolkit/jsonapi"
)

func setHeaderIfEmpty(r jsonapi.Request, key, val string) {
	if v := r.W().Header().Get(key); v == "" {
		r.W().Header().Set(key, val)
	}
}

// LastModify handles If-Modified-Since and Last-Modified
//
// It detects if there's "If-Modified-Since" header in request and
// "Last-Modified" in response. If both exist and modified time <= request,
// the response is discarded and 304 is replied.
//
// It's no-op if handler returns any error.
//
//     // If browser send If-Modified-Since >= dateA, 304 is returned, "hello" otherwise.
//     //
//     // If browser send If-Modified-Since < dateA, "hello" is always replied.
//     func h1(r jsonapi.Request) (interface{}, error) {
//         r.W().Header().Set("Last-Mdified", dateA.UTC().Format(http.TimeFormat))
//         return "hello", nil
//     }
//
//     // Always return 404 to browser.
//     func h2(r jsonapi.Request) (interface{}, error) {
//         r.W().Header().Set("Last-Mdified", dateA.UTC().Format(http.TimeFormat))
//         return 1, jsonapi.E404
//     }
func LastModify(h jsonapi.Handler) (ret jsonapi.Handler) {
	return func(r jsonapi.Request) (data interface{}, err error) {
		if data, err = h(r); err != nil {
			return
		}

		// check if response header is set
		txt := r.W().Header().Get("Last-Modified")
		if txt == "" {
			return
		}
		has, e := http.ParseTime(txt)
		if e != nil {
			return
		}

		setHeaderIfEmpty(r, "Date", time.Now().UTC().Format(http.TimeFormat))

		// check if request header is set
		txt = r.R().Header.Get("If-Modified-Since")
		if txt != "" {
			want, e := http.ParseTime(txt)
			if e == nil {
				if !has.After(want) {
					return nil, jsonapi.E304
				}
			}
		}

		return
	}
}
