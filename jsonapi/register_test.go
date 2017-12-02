package jsonapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCamelConverting(t *testing.T) {
	// [method name, expect]
	cases := [][2]string{
		{"CamelCase", "camel_case"},
		{"URL123", "url123"},
		{"TestIdGetter", "test_id_getter"},
		{"TestIDGetter", "test_id_getter"},
		{"TestID3Getter", "test_id3_getter"},
	}

	for _, c := range cases {
		t.Run(c[0], func(t *testing.T) {
			if actual := convertCamelTo_(c[0]); c[1] != actual {
				t.Fatalf("expected %s, got %s", c[1], actual)
			}
		})
	}
}

func TestRegisterOrder(t *testing.T) {
	buf := &bytes.Buffer{}

	m1 := func(h Handler) Handler {
		return Handler(func(
			d *json.Decoder,
			r *http.Request,
			w http.ResponseWriter,
		) (interface{}, error) {
			buf.WriteByte('1')
			data, err := h(d, r, w)
			buf.WriteByte('1')
			return data, err
		})
	}
	m2 := func(h Handler) Handler {
		return Handler(func(
			d *json.Decoder,
			r *http.Request,
			w http.ResponseWriter,
		) (interface{}, error) {
			buf.WriteByte('2')
			data, err := h(d, r, w)
			buf.WriteByte('2')
			return data, err
		})
	}

	h := func(
		d *json.Decoder,
		r *http.Request,
		w http.ResponseWriter,
	) (interface{}, error) {
		buf.WriteByte('3')
		return nil, nil
	}

	apis := []API{
		{Pattern: "/api", Handler: h},
	}

	mux := http.NewServeMux()
	With(m1).With(m2).Register(mux, apis)
	req := httptest.NewRequest("GET", "http://localhost/api", nil)
	handler, _ := mux.Handler(req)
	handler.ServeHTTP(httptest.NewRecorder(), req)

	if actual := buf.String(); actual != "12321" {
		t.Fatalf("expected 12321, got %s", actual)
	}
}
