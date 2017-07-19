package router

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestPathMuxRegister(t *testing.T) {
	successCases := []string{
		"/",
		"/a",
		"/a/",
		"/a/b",
		"/*",
		"/a/*",
		"/*/a",
		"//a/*",
	}
	failedCases := []string{
		"a",
		"/a*/a",
		"/a/*a",
		"/a//*",
	}

	m := ByPath()
	for _, c := range successCases {
		t.Run("Valid", func(t *testing.T) {
			defer func() {
				err := recover()
				if err == nil {
					return
				}

				t.Errorf("unexpected failure for %s: %s", c, err)
			}()
			m.HandleFunc(c, h)
		})
	}
	for _, c := range failedCases {
		t.Run("Invalid", func(t *testing.T) {
			defer func() {
				err := recover()
				if err != nil {
					return
				}

				t.Errorf("unexpected success for %s", c)
			}()
			m.HandleFunc(c, h)
		})
	}
}

func TestPathMuxMatching(t *testing.T) {
	cases := []struct {
		rule  string
		uri   string
		data  []string
		found bool
	}{
		// rule /
		{
			rule:  "/",
			uri:   "http://localhost/",
			data:  []string{},
			found: true,
		},
		{
			rule:  "/",
			uri:   "http://localhost/a",
			data:  []string{},
			found: true,
		},
		{
			rule:  "/",
			uri:   "http://localhost/a/b",
			data:  []string{},
			found: true,
		},
		// rule /a
		{
			rule:  "/a",
			uri:   "http://localhost/a",
			data:  []string{},
			found: true,
		},
		{
			rule:  "/a",
			uri:   "http://localhost/",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/a",
			uri:   "http://localhost/a/b",
			data:  []string{},
			found: false,
		},
		// rule /a/
		{
			rule:  "/a/",
			uri:   "http://localhost/a/b",
			data:  []string{},
			found: true,
		},
		{
			rule:  "/a/",
			uri:   "http://localhost/a",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/a/",
			uri:   "http://localhost/",
			data:  []string{},
			found: false,
		},
		// rule /*
		{
			rule:  "/*",
			uri:   "http://localhost/a",
			data:  []string{"a"},
			found: true,
		},
		{
			rule:  "/*",
			uri:   "http://localhost/",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/*",
			uri:   "http://localhost/a/b",
			data:  []string{},
			found: false,
		},
		// rule /*/*
		{
			rule:  "/*/*",
			uri:   "http://localhost/a",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/*/*",
			uri:   "http://localhost/a/b/c",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/*/*",
			uri:   "http://localhost/",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/*/*",
			uri:   "http://localhost/a/b",
			data:  []string{"a", "b"},
			found: true,
		},
		// rule /*/
		{
			rule:  "/*/",
			uri:   "http://localhost/a",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/*/",
			uri:   "http://localhost/a/b/c",
			data:  []string{"a"},
			found: true,
		},
		{
			rule:  "/*/",
			uri:   "http://localhost/",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/*/",
			uri:   "http://localhost/a/b",
			data:  []string{"a"},
			found: true,
		},
		// rule /a/*
		{
			rule:  "/a/*",
			uri:   "http://localhost/a",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/a/*",
			uri:   "http://localhost/a/b/c",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/a/*",
			uri:   "http://localhost/",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/a/*",
			uri:   "http://localhost/c/b",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/a/*",
			uri:   "http://localhost/a/b",
			data:  []string{"b"},
			found: true,
		},
		// rule /*/a
		{
			rule:  "/*/a",
			uri:   "http://localhost/a",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/*/a",
			uri:   "http://localhost/a/b/c",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/*/a",
			uri:   "http://localhost/",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/*/a",
			uri:   "http://localhost/a/b",
			data:  []string{},
			found: false,
		},
		{
			rule:  "/*/a",
			uri:   "http://localhost/b/a",
			data:  []string{"b"},
			found: true,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Matching[%s][%s]", c.rule, c.uri), func(t *testing.T) {
			m := ByPath()
			m.HandleFunc(c.rule, h)
			req := makeReq(c.uri)
			_, data, found := m.mappings.match(req)
			if found != c.found {
				t.Fatalf("expected matching [%s] with rule [%s] to be %t, got %t", c.uri, c.rule, c.found, found)
			}

			if !reflect.DeepEqual(data, c.data) {
				t.Fatalf("expected variables are %#v, got %#v", c.data, data)
			}
		})
	}
}

func TestFillPathVars(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		c := context.WithValue(context.Background(), PathVarKey, []string{"a", "b"})
		var a, b string

		cnt := FillPathVariable(c, &a, &b)

		if cnt != 2 {
			t.Errorf("expected to fill 2 vars, got %d", cnt)
		}

		if a != "a" {
			t.Errorf("expected a to be 'a', got '%s'", a)
		}

		if b != "b" {
			t.Errorf("expected b to be 'b', got '%s'", b)
		}
	})

	t.Run("NoData", func(t *testing.T) {
		c := context.WithValue(context.Background(), "ToT", []string{"a", "b"})
		var a, b string

		cnt := FillPathVariable(c, &a, &b)

		if cnt != 0 {
			t.Errorf(
				"expected to fail filling, got %d with '%s' and '%s'",
				cnt,
				a,
				b,
			)
		}
	})
}
