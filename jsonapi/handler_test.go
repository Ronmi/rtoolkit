package jsonapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	fac := func(data interface{}, err error) Handler {
		return func(
			d *json.Decoder, r *http.Request, w http.ResponseWriter,
		) (interface{}, error) {
			return data, err
		}
	}
	run := func(h Handler) *httptest.ResponseRecorder {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}

	basicCases := []struct {
		name   string
		expect string
		h      Handler
		status int
	}{
		{
			name:   "nil-nil",
			expect: `{"data":null}`,
			h:      fac(nil, nil),
			status: 200,
		},
		{
			name:   "int-nil",
			expect: `{"data":1}`,
			h:      fac(1, nil),
			status: 200,
		},
		{
			name:   "string-nil",
			expect: `{"data":"asd"}`,
			h:      fac("asd", nil),
			status: 200,
		},
		{
			name:   "struct-nil",
			expect: `{"data":{"A":1}}`,
			h:      fac(struct{ A int }{A: 1}, nil),
			status: 200,
		},
		{
			name:   "chan-nil",
			expect: `{"errors":[{"detail":"Failed to marshal data"}]}`,
			h:      fac(make(chan int), nil),
			status: 500,
		},
		{
			name:   "nil-error",
			expect: `{"errors":[{"detail":"my error"}]}`,
			h:      fac(nil, errors.New("my error")),
			status: 500,
		},
		{
			name:   "data-error",
			expect: `{"errors":[{"detail":"my error"}]}`,
			h:      fac(1, errors.New("my error")),
			status: 500,
		},
		{
			name:   "nil-APPERR",
			expect: `{"errors":[{"code":"qq","detail":"my error"}]}`,
			h:      fac(nil, APPERR.SetCode("qq").SetData("my error")),
			status: 200,
		},
		{
			name:   "data-APPERR",
			expect: `{"errors":[{"code":"qq","detail":"my error"}]}`,
			h:      fac(1, APPERR.SetCode("qq").SetData("my error")),
			status: 200,
		},
		{
			name:   "nil-E404",
			expect: `{"errors":[{"code":"qq","detail":"my error"}]}`,
			h:      fac(nil, E404.SetCode("qq").SetData("my error")),
			status: 404,
		},
		{
			name:   "data-E404",
			expect: `{"errors":[{"code":"qq","detail":"my error"}]}`,
			h:      fac(1, E404.SetCode("qq").SetData("my error")),
			status: 404,
		},
	}

	for _, c := range basicCases {
		c.expect += "\n"
		t.Run(c.name, func(t *testing.T) {
			w := run(c.h)
			if w.Code != c.status {
				t.Errorf(
					"expected status %d, got %d",
					c.status,
					w.Code,
				)
			}
			if actual := w.Body.String(); c.expect != actual {
				t.Errorf(
					"expected %#v, got %#v",
					c.expect,
					actual,
				)
			}
		})
	}

	out := "<a href=\"http://a.b\">Moved Permanently</a>.\n\n"
	uri := "http://a.b"
	t.Run("nil-E301", func(t *testing.T) {
		w := run(fac(nil, E301.SetData(uri)))
		if b := w.Body.String(); b != out {
			t.Errorf("unexpected content: %s", b)
		}
		if w.Code != 301 {
			t.Errorf("unexpected status: %d", w.Code)
		}
		if l := w.HeaderMap.Get("Location"); l != uri {
			t.Errorf("unexpected location header: %s", l)
		}
	})
	t.Run("data-E301", func(t *testing.T) {
		w := run(fac(1, E301.SetData(uri)))
		if b := w.Body.String(); b != out {
			t.Errorf("unexpected content: %s", b)
		}
		if w.Code != 301 {
			t.Errorf("unexpected status: %d", w.Code)
		}
		if l := w.HeaderMap.Get("Location"); l != uri {
			t.Errorf("unexpected location header: %s", l)
		}
	})
}
