package jsonapi

import (
	"encoding/json"
	"net/http"
)

// Request represents most used data a handler need
type Request interface {
	Decode(interface{}) error
	R() *http.Request
	W() http.ResponseWriter
}

// FakeRequest implements a Request and let you do some magic in it
type FakeRequest struct {
	Decoder *json.Decoder
	Req     *http.Request
	Resp    http.ResponseWriter
}

func (r *FakeRequest) Decode(data interface{}) error {
	return r.Decoder.Decode(data)
}

func (r *FakeRequest) R() *http.Request {
	return r.Req
}

func (r *FakeRequest) W() http.ResponseWriter {
	return r.Resp
}

// FromHTTP creates a Request instance from http request and response
func FromHTTP(w http.ResponseWriter, r *http.Request) Request {
	dec := json.NewDecoder(r.Body)
	return &FakeRequest{
		Decoder: dec,
		Req:     r,
		Resp:    w,
	}
}

type reqWrapper struct {
	Request
	req *http.Request
}

func (r *reqWrapper) R() *http.Request {
	return r.req
}

// WrapRequest creates a new Request, with http request replaced
func WrapRequest(q Request, r *http.Request) Request {
	return &reqWrapper{
		Request: q,
		req:     r,
	}
}

type respWrapper struct {
	Request
	resp http.ResponseWriter
}

func (r *respWrapper) W() http.ResponseWriter {
	return r.resp
}

// WrapResponse creates a new Request, with http response replaced
func WrapResponse(q Request, w http.ResponseWriter) Request {
	return &respWrapper{
		Request: q,
		resp:    w,
	}
}
