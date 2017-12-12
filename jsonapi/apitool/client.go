package apitool

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Ronmi/rtoolkit/jsonapi"
)

// EClient indicates something goes wrong at client side while calling remote jsonapi
var EClient = jsonapi.Error{Code: -1}.SetData("Client error")

// Client is a helper to simplify the process of calling jsonapi
//
// Any error during the call process will immediately return as
// jsonapi.InternalError.SetOrigin(the_error)
type Client interface {
	// Synchronized call
	//
	// If any io error or json decoding error occured, an
	// EClient.SetOrigin(the_error) returns.
	Exec(param, result interface{}) error
	// Asynchronized call, Client take response to close the channel
	// result is guaranteed to be filled when error returns.
	//
	// If any io error or json decoding error occured, an
	// EClient.SetOrigin(the_error) returns.
	Do(param, result interface{}) chan error
}

type clientFunc func(interface{}, interface{}) error

// Exec calls to specified
func (c clientFunc) Exec(param, result interface{}) error {
	return c(param, result)
}

func (c clientFunc) Do(param, result interface{}) chan error {
	ret := make(chan error, 1)

	go func() {
		defer close(ret)

		ret <- c.Exec(param, result)
	}()

	return ret
}

type callResp struct {
	Data   *json.RawMessage `json:"data"`
	Errors []jsonapi.ErrObj `json:"errors"`
}

// ParseResponse parses response of a jsonapi
//
// It's caller's response to close response body.
//
// If any io error or json decoding error occured, an
// EClient.SetOrigin(the_error) returns.
func ParseResponse(resp *http.Response, result interface{}) error {
	var res callResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return EClient.SetOrigin(err)
	}

	if d := res.Data; d != nil {
		if err := json.Unmarshal([]byte(*d), result); err != nil {
			return EClient.SetOrigin(err)
		}
	}

	if len(res.Errors) == 0 {
		return nil
	}

	return res.Errors[0].AsError()
}

// Call creates an Client to a jsonapi entry
//
// It will use http.DefaultClient if c == nil, but it's not recommended.
func Call(method, uri string, client *http.Client) Client {
	c := client
	if c == nil {
		c = http.DefaultClient
	}

	return clientFunc(func(param, result interface{}) error {
		data, err := json.Marshal(param)
		if err != nil {
			return EClient.SetOrigin(err)
		}

		req, err := http.NewRequest(method, uri, bytes.NewReader(data))
		if err != nil {
			return EClient.SetOrigin(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		req = req.WithContext(ctx)

		resp, err := c.Do(req)
		if err != nil {
			cancel()
			return EClient.SetOrigin(err)
		}
		defer resp.Body.Close()
		defer io.Copy(ioutil.Discard, resp.Body)
		defer cancel()

		return ParseResponse(resp, result)
	})
}
