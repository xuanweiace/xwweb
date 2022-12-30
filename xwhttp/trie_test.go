package xwhttp

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_parse_route(t *testing.T) {
	//*的作用
	fmt.Println(parsePattern("/*"))
	fmt.Println(parsePattern("/*/a"))

	//
	fmt.Println(parsePattern("/a"))
	fmt.Println(parsePattern("/a/"))
	fmt.Println(parsePattern("//a/")) //过滤不掉这种非法情况
	fmt.Println(parsePattern("/a/a"))
	fmt.Println(parsePattern("/:name"))
	fmt.Println(parsePattern("/:name:abc/age"))

	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}
func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/htttttp")

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if ps["name"] != "htttttp" {
		t.Fatal("name should be equal to 'geektutu'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

}
