package jsonapi

import (
	"encoding/json"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

// HTTPMux abstracts http.ServeHTTPMux, so it will be easier to write tests
//
// Only needed methods are added here.
type HTTPMux interface {
	Handle(pattern string, handler http.Handler)
}

// API denotes how a json api handler registers to a servemux
type API struct {
	Pattern string
	Handler func(dec *json.Decoder, r *http.Request, w http.ResponseWriter) (interface{}, error)
}

// Register helps you to register many APIHandlers to a http.ServeHTTPMux
func Register(mux HTTPMux, apis []API) {
	reg := http.Handle
	if mux != nil {
		reg = mux.Handle
	}

	for _, api := range apis {
		reg(api.Pattern, Handler(api.Handler))
	}
}

var reCamelTo_ *regexp.Regexp
var reCamelTo_Excepts *regexp.Regexp

func init() {
	reCamelTo_ = regexp.MustCompile(
		`([^A-Z])([A-Z])|([A-Z0-9]+)([A-Z])`,
	)
	reCamelTo_Excepts = regexp.MustCompile(
		`^[A-Z0-9]*$`,
	)
}

func findMatchedMethods(prefix string, handlers interface{}) []API {
	v := reflect.ValueOf(handlers)

	ret := make([]API, 0, v.NumMethod())

	for x, t := 0, v.Type(); x < v.NumMethod(); x++ {
		h, ok := v.Method(x).Interface().(func(dec *json.Decoder, r *http.Request, w http.ResponseWriter) (interface{}, error))
		if !ok {
			// incorrect signature, skip
			continue
		}

		ret = append(ret, API{
			Pattern: prefix + "/" + convertCamelTo_(t.Method(x).Name),
			Handler: h,
		})
	}

	return ret
}

// RegisterAll helps you to register all handler methods
//
// As using reflection to do the job, only exported methods with correct
// signature are registered.
//
// The pattern are generated by converting CamelCase to
// underscore_pattern then add prefix and "/" before it. Take a look at
// the test cases as example.
func RegisterAll(mux HTTPMux, prefix string, handlers interface{}) {
	Register(mux, findMatchedMethods(prefix, handlers))
}

func convertCamelTo_(name string) string {
	if reCamelTo_Excepts.MatchString(name) {
		return strings.ToLower(name)
	}

	return strings.ToLower(
		reCamelTo_.ReplaceAllString(
			name,
			"${1}${3}_${2}${4}",
		),
	)
}
