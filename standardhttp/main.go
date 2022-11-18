package standardhttp

import (
	"fmt"
	"log"
	"net/http"
)

func Main() {
	//方法1: 调用 http.HandleFunc 实现了路由和Handler的映射，也就是只能针对具体的路由写处理逻辑
	//http.HandleFunc("/", default_handler)
	//http.HandleFunc("/hello", hello_handler)
	//log.Fatal(http.ListenAndServe("localhost:8080", nil))

	//方法2: 我们拦截了所有的HTTP请求，拥有了统一的控制入口。在这里我们可以自由定义路由映射的规则，也可以统一添加一些处理逻辑，例如日志、异常处理等。
	handlers := HandlerContainer{}
	log.Fatal(http.ListenAndServe("localhost:8080", handlers))

}

type HandlerContainer struct{}

func (h HandlerContainer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/":
		default_handler(writer, request)
	case "/hello":
		hello_handler(writer, request)
	default:
		fmt.Fprintf(writer, "404 not found： %s\n", request.URL)

	}
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
