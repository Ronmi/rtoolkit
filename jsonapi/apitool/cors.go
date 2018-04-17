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

// CORS is a middleware simply allow any host access your api by setting
// "Access-Control-Allow-Origin: *"
func CORS(h jsonapi.Handler) jsonapi.Handler {
	return ForceHeader(
		map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	)(h)
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

// NewCORS creates a middleware to set CORS headers
//
// Fields with zero value will not be set.
func NewCORS(opt CORSOption) jsonapi.Middleware {
	h := make(map[string]string)
	if opt.Origin != "" {
		h["Access-Control-Allow-Origin"] = opt.Origin
	}
	if len(opt.ExposeHeaders) > 0 {
		h["Access-Control-Expose-Headers"] = strings.Join(
			opt.ExposeHeaders, ", ")
	}
	if len(opt.Headers) > 0 {
		h["Access-Control-Allow-Headers"] = strings.Join(
			opt.Headers, ", ")
	}
	if opt.MaxAge > 0 {
		h["Access-Control-Max-Age"] = strconv.FormatUint(
			opt.MaxAge, 10)
	}
	if opt.Credential {
		h["Access-Control-Allow-Credentials"] = "true"
	}
	if len(opt.Methods) > 0 {
		h["Access-Control-Allow-Methods"] = strings.Join(
			opt.Methods, ", ")
	}
	return ForceHeader(h)
}
