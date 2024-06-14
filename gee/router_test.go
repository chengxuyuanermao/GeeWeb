package gee

import (
	"reflect"
	"testing"
)

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/hello"), []string{"p", "hello"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/:hello"), []string{"p", ":hello"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/:hello/play"), []string{"p", ":hello", "play"})
	if !ok {
		t.Fatal("not deep equal")
	}
}

func newTestRouter() *Router {
	r := NewRouter()
	r.AddRoute("GET", "/hello/cwx", nil)
	r.AddRoute("GET", "/hello/:name", nil)
	r.AddRoute("GET", "/hello/:name/index", nil)
	return r
}

func TestGetRouter(t *testing.T) {
	r := newTestRouter()
	n, params := r.getRoute("GET", "/hello/geetutu/index")
	if n == nil {
		t.Fatal("node is nil")
	}
	t.Logf("node is %v", *n)

	if len(params) == 0 {
		t.Fatal("params err")
	}
	t.Logf("params is %q", params)
}
