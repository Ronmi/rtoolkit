package apitest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/Ronmi/rtoolkit/jsonapi"
)

// NewRequest wraps httptest.NewRequest, use your data (encoded to JSON) as request body
//
// It also sets "Content-Type" to "application/json".
func NewRequest(method, target string, data interface{}) *http.Request {
	buf, _ := json.Marshal(data)
	return httptest.NewRequest(method, target, bytes.NewReader(buf))
}

// Test wraps your handler for test purpose
type Test jsonapi.Handler

// With creates new Test instance by wrapping the handler with the middleware
func (t Test) With(m jsonapi.Middleware) Test {
	return Test(m(jsonapi.Handler(t)))
}

// UseRequest executes handler with specified request
func (t Test) UseRequest(req *http.Request) (interface{}, error) {
	defer req.Body.Close()

	w := httptest.NewRecorder()
	dec := json.NewDecoder(req.Body)
	return t(dec, req, w)
}

// Use executes handler with your data
//
// The request address will be "/" and using POST method.
func (t Test) Use(data interface{}) (interface{}, error) {
	return t.UseRequest(
		NewRequest("POST", "/", data),
	)
}
