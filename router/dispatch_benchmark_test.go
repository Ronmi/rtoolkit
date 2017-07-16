package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
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

func benchDispatch(n int, b *testing.B) {
	rules := createRules(n, "{var1}", "{var2}")
	m := mux.NewRouter()
	for _, r := range rules {
		m.HandleFunc(r, h)
	}
	sz := len(rules)
	min, max := 0, sz-1
	mid := (min + max) / 2
	tmpl := "http://localhost/lv1/lv2/%d/var1/var2"
	minReq := makeReq(fmt.Sprintf(tmpl, min))
	midReq := makeReq(fmt.Sprintf(tmpl, mid))
	maxReq := makeReq(fmt.Sprintf(tmpl, max))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(httptest.NewRecorder(), minReq)
		m.ServeHTTP(httptest.NewRecorder(), midReq)
		m.ServeHTTP(httptest.NewRecorder(), maxReq)
	}
}

func BenchmarkDispatch100(b *testing.B) {
	benchDispatch(100, b)
}

func BenchmarkDispatch200(b *testing.B) {
	benchDispatch(200, b)
}

func BenchmarkDispatch500(b *testing.B) {
	benchDispatch(500, b)
}

func BenchmarkDispatch1000(b *testing.B) {
	benchDispatch(1000, b)
}
