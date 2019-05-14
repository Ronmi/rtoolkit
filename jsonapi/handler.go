package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// these codes are inspired by http://go-talks.appspot.com/github.com/broady/talks/web-frameworks-gophercon.slide#1

// ErrObj defines how an error is exported to client
//
// For jsonapi.Error, Code will contains result of SetCode; Detail will be SetData
//
// For other error types, only Detail is set, as error.Error()
type ErrObj struct {
	Code   string `json:"code,omitempty"`
	Detail string `json:"detail,omitempty"`
}

// AsError creates an error object represents this error
//
// If Code is set, an Error instance will be returned. errors.New(Detail) otherwise.
func (o *ErrObj) AsError() error {
	if o.Code != "" {
		return Error{message: o.Detail, errCode: o.Code}
	}

	return errors.New(o.Detail)
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

func fromError(e *Error) *ErrObj {
	return &ErrObj{
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
	EUnknown = Error{Code: 0, message: "Unknown error"}
	E301     = Error{Code: 301, message: "Resource has been moved permanently"}
	E302     = Error{Code: 302, message: "Resource has been found at another location"}
	E303     = Error{Code: 303, message: "See other"}
	E304     = Error{Code: 304, message: "Not modified"}
	E307     = Error{Code: 307, message: "Resource has been moved to another location temporarily"}
	E400     = Error{Code: 400, message: "Error parsing request"}
	E401     = Error{Code: 401, message: "You have to be authorized before accessing this resource"}
	E403     = Error{Code: 403, message: "You have no right to access this resource"}
	E404     = Error{Code: 404, message: "Resource not found"}
	E408     = Error{Code: 408, message: "Request timeout"}
	E409     = Error{Code: 409, message: "Conflict"}
	E410     = Error{Code: 410, message: "Gone"}
	E413     = Error{Code: 413, message: "Request entity too large"}
	E415     = Error{Code: 415, message: "Unsupported media type"}
	E418     = Error{Code: 418, message: "I'm a teapot"}
	E426     = Error{Code: 426, message: "Upgrade required"}
	E429     = Error{Code: 429, message: "Too many requests"}
	E500     = Error{Code: 500, message: "Internal server error"}
	E501     = Error{Code: 501, message: "Not implemented"}
	E502     = Error{Code: 502, message: "Bad gateway"}
	E503     = Error{Code: 503, message: "Service unavailable"}
	E504     = Error{Code: 504, message: "Gateway timeout"}

	// application-defined error
	APPERR = Error{Code: 200}

	// special error, preventing ServeHTTP method to encode the returned data
	//
	// For string, []byte or anything implements fmt.Stringer returned, we will
	// write it to response as-is.
	//
	// For other type, we use fmt.FPrintf(responseWriter, "%v", returnedData).
	//
	// You will also have to:
	//    - Set HTTP status code manually.
	//    - Set necessary response headers manually.
	//    - Take care not to be overwritten by middleware.
	ASIS = Error{Code: -1}
)

// Failed wraps you error object and prepares suitable return type to be used in controller
//
// Here's a common usage:
//
//   if err := req.Decode(&param); err != nil {
//       return jsonapi.E(err, jsonapi.E400.SetData("invalid parameter"))
//   }
//   if err := param.IsValid(); err != nil {
//       return jsonapi.E(err, jsonapi.E400.SetData("invalid parameter"))
//   }
func Failed(e1 error, e2 Error) (data interface{}, err error) {
	return data, e2.SetOrigin(e1)
}

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
// Redirecting depends on http.Redirect(). The data returned from handler will never
// write to ResponseWriter.
//
// This basically obey the http://jsonapi.org rules:
//
//     - Return {"data": your_data} if error == nil
//     - Return {"errors": [{"code": application-defined-error-code, "detail": message}]} if error returned
type Handler func(r Request) (interface{}, error)

// ServeHTTP implements net/http.Handler
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	res, err := h(FromHTTP(w, r))
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
		if httperr.EqualTo(ASIS) {
			if res != nil {
				switch x := res.(type) {
				case string:
					w.Write([]byte(x))
				case []byte:
					w.Write(x)
				case fmt.Stringer:
					w.Write([]byte(x.String()))
				default:
					fmt.Fprintf(w, "%v", res)
				}
			}
			return
		}
		code = httperr.Code
		if code >= 301 && code <= 303 && httperr.location != "" {
			// 301~303 redirect
			http.Redirect(w, r, httperr.location, code)
			return
		}

		w.WriteHeader(code)
		resp["errors"] = []*ErrObj{fromError(&httperr)}
		enc.Encode(resp)
		return
	}

	w.WriteHeader(code)
	resp["errors"] = []*ErrObj{&ErrObj{Detail: err.Error()}}
	enc.Encode(resp)
}
