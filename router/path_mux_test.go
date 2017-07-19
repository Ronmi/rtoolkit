package router

import "testing"

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
		found bool
	}{
		// rule /
		{
			rule:  "/",
			uri:   "http://localhost/",
			found: true,
		},
		{
			rule:  "/",
			uri:   "http://localhost/a",
			found: true,
		},
		{
			rule:  "/",
			uri:   "http://localhost/a/b",
			found: true,
		},
		// rule /a
		{
			rule:  "/a",
			uri:   "http://localhost/a",
			found: true,
		},
		{
			rule:  "/a",
			uri:   "http://localhost/",
			found: false,
		},
		{
			rule:  "/a",
			uri:   "http://localhost/a/b",
			found: false,
		},
		// rule /a/
		{
			rule:  "/a/",
			uri:   "http://localhost/a/b",
			found: true,
		},
		{
			rule:  "/a/",
			uri:   "http://localhost/a",
			found: false,
		},
		{
			rule:  "/a/",
			uri:   "http://localhost/",
			found: false,
		},
		// rule /*
		{
			rule:  "/*",
			uri:   "http://localhost/a",
			found: true,
		},
		{
			rule:  "/*",
			uri:   "http://localhost/",
			found: true,
		},
		{
			rule:  "/*",
			uri:   "http://localhost/a/b",
			found: false,
		},
		// rule /*/*
		{
			rule:  "/*/*",
			uri:   "http://localhost/a",
			found: false,
		},
		{
			rule:  "/*/*",
			uri:   "http://localhost/a/b/c",
			found: false,
		},
		{
			rule:  "/*/*",
			uri:   "http://localhost/",
			found: false,
		},
		{
			rule:  "/*/*",
			uri:   "http://localhost/a/b",
			found: true,
		},
		// rule /*/
		{
			rule:  "/*/",
			uri:   "http://localhost/a",
			found: false,
		},
		{
			rule:  "/*/",
			uri:   "http://localhost/a/b/c",
			found: true,
		},
		{
			rule:  "/*/",
			uri:   "http://localhost/",
			found: false,
		},
		{
			rule:  "/*/",
			uri:   "http://localhost/a/b",
			found: true,
		},
		// rule /a/*
		{
			rule:  "/a/*",
			uri:   "http://localhost/a",
			found: false,
		},
		{
			rule:  "/a/*",
			uri:   "http://localhost/a/b/c",
			found: false,
		},
		{
			rule:  "/a/*",
			uri:   "http://localhost/",
			found: false,
		},
		{
			rule:  "/a/*",
			uri:   "http://localhost/c/b",
			found: false,
		},
		{
			rule:  "/a/*",
			uri:   "http://localhost/a/b",
			found: true,
		},
		// rule /*/a
		{
			rule:  "/*/a",
			uri:   "http://localhost/a",
			found: false,
		},
		{
			rule:  "/*/a",
			uri:   "http://localhost/a/b/c",
			found: false,
		},
		{
			rule:  "/*/a",
			uri:   "http://localhost/",
			found: false,
		},
		{
			rule:  "/*/a",
			uri:   "http://localhost/a/b",
			found: false,
		},
		{
			rule:  "/*/a",
			uri:   "http://localhost/b/a",
			found: true,
		},
	}

	for _, c := range cases {
		t.Run("Matching", func(t *testing.T) {
			m := ByPath()
			m.HandleFunc(c.rule, h)
			req := makeReq(c.uri)
			_, found := m.mappings.match(req)
			if found != c.found {
				t.Fatalf("expected matching [%s] with rule [%s] to be %t, got %t", c.uri, c.rule, c.found, found)
			}
		})
	}
}
