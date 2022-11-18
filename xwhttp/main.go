package xwhttp

import (
	"fmt"
	"net/http"
)

func Main() {
	server := NewInstance()
	server.Get("/", default_handler)
	server.Get("/hello", hello_handler)
	server.Run("localhost:8080")
}

// func (container *HandlerContainer) Server
func default_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}
func hello_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello URL.Path = %q\n", r.URL.Path)
	for k, v := range r.Header {
		fmt.Fprintf(w, "k=%s, v=%s\n", k, v)
	}
}
