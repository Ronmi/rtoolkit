package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func h(w http.ResponseWriter, r *http.Request) {}

func makeReq(uri string) *http.Request {
	u, _ := url.Parse(uri)
	i, o := io.Pipe()
	_ = o.Close()
	ret := &http.Request{
		Method:        "GET",
		URL:           u,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Body:          i,
		ContentLength: 0,
		Host:          u.Hostname(),
		RemoteAddr:    "127.0.0.1:12345",
		RequestURI:    uri,
	}

	return ret
}

func createRules(n int, var1, var2 string) []string {
	routingRules := make([]string, n)

	for i := 0; i < n; i++ {
		routingRules[i] = fmt.Sprintf("/lv1/lv2/%d/%s/%s", i+1, var1, var2)
	}

	return routingRules
}

func benchDispatch(m http.Handler, req *http.Request, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func createMux(n int) *PathMux {
	rules := createRules(n, "*", "*")
	m := ByPath()
	for _, r := range rules {
		m.HandleFunc(r, h)
	}

	return m
}

func BenchmarkRouterInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ByPath()
	}
}

func BenchmarkRouterRegister(b *testing.B) {
	f := func(m *PathMux, rules []string, n int) {
		for i := 0; i < n; i++ {
			for _, r := range rules {
				m.HandleFunc(r, h)
			}
		}
	}
	cases := []int{
		100, 200, 500, 1000,
	}

	for _, c := range cases {
		b.Run(
			fmt.Sprintf("%d Rules", c),
			func(b *testing.B) {
				rules := createRules(c, "*", "*")
				m := ByPath()
				b.ResetTimer()
				f(m, rules, b.N)
			},
		)
	}
}

func BenchmarkRouterDispatch(b *testing.B) {
	tmpl := "http://localhost/lv1/lv2/%d/var1/var2"
	cases := []struct {
		n    int
		pos  int
		name string
	}{
		{100, 0, "HEAD"},
		{100, 50, "MID"},
		{100, 99, "TAIL"},
		{100, 100, "404"},
		{200, 0, "HEAD"},
		{200, 100, "MID"},
		{200, 199, "TAIL"},
		{200, 200, "404"},
		{500, 0, "HEAD"},
		{500, 250, "MID"},
		{500, 499, "TAIL"},
		{500, 500, "404"},
		{1000, 0, "HEAD"},
		{1000, 500, "MID"},
		{1000, 999, "TAIL"},
		{1000, 1000, "404"},
	}

	for _, c := range cases {
		b.Run(
			fmt.Sprintf("%d-%s", c.n, c.name),
			func(b *testing.B) {
				benchDispatch(
					createMux(c.n),
					makeReq(fmt.Sprintf(tmpl, c.pos)),
					b,
				)
			},
		)
	}
}
