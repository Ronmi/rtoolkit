package apitest

import (
	"bytes"
	"testing"

	"github.com/Ronmi/rtoolkit/jsonapi"
)

func TestWithOrder(t *testing.T) {
	buf := &bytes.Buffer{}

	m1 := func(h jsonapi.Handler) jsonapi.Handler {
		return func(r jsonapi.Request) (interface{}, error) {
			buf.WriteByte('1')
			data, err := h(r)
			buf.WriteByte('1')
			return data, err
		}
	}
	m2 := func(h jsonapi.Handler) jsonapi.Handler {
		return func(r jsonapi.Request) (interface{}, error) {
			buf.WriteByte('2')
			data, err := h(r)
			buf.WriteByte('2')
			return data, err
		}
	}

	h := func(r jsonapi.Request) (interface{}, error) {
		buf.WriteByte('3')
		return nil, nil
	}

	Test(h).With(m1).With(m2).Use(nil)

	if actual := buf.String(); actual != "21312" {
		t.Fatalf("expected 21312, got %s", actual)
	}
}
