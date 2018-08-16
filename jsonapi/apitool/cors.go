package apitool

import (
	"strconv"
	"strings"

	"github.com/Ronmi/rtoolkit/jsonapi"
)

// ForceHeader creates a middleware to enforce response header
func ForceHeader(headers map[string]string) jsonapi.Middleware {
	return func(h jsonapi.Handler) (ret jsonapi.Handler) {
		return func(r jsonapi.Request) (data interface{}, err error) {
			data, err = h(r)
			for k, v := range headers {
				r.W().Header().Set(k, v)
			}
			return
		}
	}
}

// CORSOption defines supported parameters used by NewCORS()
type CORSOption struct {
	Origin        string
	ExposeHeaders []string
	Headers       []string
	MaxAge        uint64
	Credential    bool
	Methods       []string
}

const (
	hOrigin      = "Access-Control-Allow-Origin"
	hCredentials = "Access-Control-Allow-Credentials"
	hMethods     = "Access-Control-Allow-Methods"
	hHeaders     = "Access-Control-Allow-Headers"
	hExpose      = "Access-Control-Expose-Headers"
	hMaxAge      = "Access-Control-Max-Age"
)

func doCORS(h jsonapi.Handler, opt CORSOption) (ret jsonapi.Handler) {
	// cache headers for simple requests
	simpleHeaders := map[string]string{
		hOrigin: opt.Origin,
	}
	if opt.Credential {
		simpleHeaders[hCredentials] = "true"
	}
	if opt.MaxAge > 0 {
		simpleHeaders[hMaxAge] = strconv.FormatUint(opt.MaxAge, 10)
	}
	simple := ForceHeader(simpleHeaders)(h)

	return func(r jsonapi.Request) (data interface{}, err error) {
		if r.R().Method != "OPTIONS" {
			// simple or normal request
			return simple(r)
		}

		hdr := map[string]string{}
		// parse Access-Control-Request-Method
		x := r.R().Header.Get("Access-Control-Request-Method")
		if x == "" {
			// not preflighted request
			return simple(r)
		}
		hdr[hMethods] = x
		if len(opt.Methods) > 0 {
			hdr[hMethods] = strings.Join(opt.Methods, ", ")
		}

		// parse Access-Control-Request-Headers
		x = r.R().Header.Get("Access-Control-Request-Headers")
		if x == "" {
			// not preflighted request
			return simple(r)
		}
		hdr[hHeaders] = x
		if len(opt.Methods) > 0 {
			hdr[hMethods] = strings.Join(opt.Methods, ", ")
		}

		// other headers
		if len(opt.ExposeHeaders) > 0 {
			hdr[hExpose] = strings.Join(opt.ExposeHeaders, ", ")
		}
		if opt.Credential {
			hdr[hCredentials] = "true"
		}
		if opt.MaxAge > 0 {
			hdr[hMaxAge] = strconv.FormatUint(opt.MaxAge, 10)
		}

		return ForceHeader(hdr)(h)(r)
	}
}

// NewCORS creates a middleware to set CORS headers
//
// Fields with zero value will not be set.
func NewCORS(opt CORSOption) jsonapi.Middleware {
	return func(h jsonapi.Handler) jsonapi.Handler {
		return doCORS(h, opt)
	}
}

// CORS is a middleware simply allow any host access your api by setting
// "Access-Control-Allow-Origin: *"
func CORS(h jsonapi.Handler) (ret jsonapi.Handler) {
	return doCORS(h, CORSOption{Origin: "*"})
}
