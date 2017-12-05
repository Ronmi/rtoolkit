package jsonapi

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// these codes are inspired by http://go-talks.appspot.com/github.com/broady/talks/web-frameworks-gophercon.slide#1

type errObj struct {
	Code   string `json:"code,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// Error represents an error status of the HTTP request. Used with APIHandler.
type Error struct {
	Code     int
	Origin   error // prepared for application errors
	message  string
	location string // url for 3xx redirect
	errCode  string
}

// Data retrieves user defined error message
func (h Error) Data() string {
	return h.message
}

// ErrCode retrieves user defined error code
func (h Error) ErrCode() string {
	return h.errCode
}

// SetOrigin creates a new Error instance to preserve original error
func (h Error) SetOrigin(err error) Error {
	h.Origin = err
	return h
}

// EqualTo tells if two Error instances represents same kind of error
//
// It compares all fields no matter exported or not, excepts Origin
func (h Error) EqualTo(e Error) bool {
	switch {
	case e.errCode != h.errCode:
		return false
	case e.message != h.message:
		return false
	case e.location != h.location:
		return false
	case e.Code != h.Code:
		return false
	}

	return true
}

// SetData creates a new Error instance and set the error message or url according to the error code
func (h Error) SetData(data string) Error {
	if h.Code >= 301 && h.Code <= 303 {
		h.location = data
		return h
	}

	h.message = data
	return h
}

// SetCode forks a new instance with application-defined error code
func (h Error) SetCode(code string) Error {
	h.errCode = code
	return h
}

func (h Error) Error() string {
	ret := strconv.Itoa(h.Code)
	if h.message != "" {
		ret += ": " + h.message
	}

	if h.location != "" {
		ret += ": " + h.location
	}

	return ret
}

func (h Error) String() string {
	ret := h.Error()
	if h.Origin != nil {
		ret += ": " + h.Origin.Error()
	}

	return ret
}

func fromError(e *Error) *errObj {
	return &errObj{
		Code:   e.errCode,
		Detail: e.message,
	}
}

// here are predefined error instances, you should call SetData before use it like
//
//     return nil, E404.SetData("User not found")
//
// You might noticed that here's no 500 error. You should just return a normal error
// instance instead.
//
//     return nil, errors.New("internal server error")
var (
	E301 = Error{Code: 301, message: "Resource has been moved permanently"}
	E302 = Error{Code: 302, message: "Resource has been found at another location"}
	E303 = Error{Code: 303, message: "See other"}
	E304 = Error{Code: 304, message: "Not modified"}
	E307 = Error{Code: 307, message: "Resource has been moved to another location temporarily"}
	E400 = Error{Code: 400, message: "Error parsing request"}
	E401 = Error{Code: 401, message: "You have to be authorized before accessing this resource"}
	E403 = Error{Code: 403, message: "You have no right to access this resource"}
	E404 = Error{Code: 404, message: "Resource not found"}
	E408 = Error{Code: 408, message: "Request timeout"}
	E409 = Error{Code: 409, message: "Conflict"}
	E410 = Error{Code: 410, message: "Gone"}
	E413 = Error{Code: 413, message: "Request entity too large"}
	E415 = Error{Code: 415, message: "Unsupported media type"}
	E418 = Error{Code: 418, message: "I'm a teapot"}
	E426 = Error{Code: 426, message: "Upgrade required"}
	E429 = Error{Code: 429, message: "Too many requests"}
	E500 = Error{Code: 500, message: "Internal server error"}
	E501 = Error{Code: 501, message: "Not implemented"}
	E502 = Error{Code: 502, message: "Bad gateway"}
	E503 = Error{Code: 503, message: "Service unavailable"}
	E504 = Error{Code: 504, message: "Gateway timeout"}

	// application-defined error
	APPERR = Error{Code: 200}
)

// Handler is easy to use entry for API developer.
//
// Just return something, and it will be encoded to JSON format and send to client.
// Or return an Error to specify http status code and error string.
//
//     func myHandler(dec *json.Decoder, httpData *HTTP) (interface{}, error) {
//         var param paramType
//         if err := dec.Decode(&param); err != nil {
//             return nil, jsonapi.E400.SetData("You must send parameters in JSON format.")
//         }
//         return doSomething(param), nil
//     }
//
// To redirect clients, return 301~303 status code and set Data property
//
//     return nil, jsonapi.E301.SetData("http://google.com")
//
// This basically obey the http://jsonapi.org rules:
//
//     - Return {"data": your_data} if error == nil
//     - Return {"errors": [{"code": application-defined-error-code, "detail": message}]} if error returned
type Handler func(dec *json.Decoder, r *http.Request, w http.ResponseWriter) (interface{}, error)

// ServeHTTP implements net/http.Handler
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)
	res, err := h(dec, r, w)
	resp := make(map[string]interface{})
	if err == nil {
		resp["data"] = res
		e := enc.Encode(resp)
		if e == nil {
			return
		}
		delete(resp, "data")

		err = E500.SetOrigin(e).SetData(
			`Failed to marshal data`,
		)
	}

	code := http.StatusInternalServerError
	if httperr, ok := err.(Error); ok {
		code = httperr.Code
		if code >= 301 && code <= 303 && httperr.location != "" {
			// 301~303 redirect
			http.Redirect(w, r, httperr.location, code)
			return
		}

		w.WriteHeader(code)
		resp["errors"] = []*errObj{fromError(&httperr)}
		enc.Encode(resp)
		return
	}

	w.WriteHeader(code)
	resp["errors"] = []*errObj{&errObj{Detail: err.Error()}}
	enc.Encode(resp)
}
