package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"xwace/xwweb/xwhttp"
)

func main() {
	//standardhttp.Main()
	engine := xwhttp.NewInstance()
	engine.Get("/", html_handler)
	engine.Get("/get_json", json_handler)
	engine.Get("/get_post_html", html_handler)
	engine.Post("/get_post_html", html_handler)
	engine.Get("/get_string", string_handler)

	v1 := engine.Group("/v1")
	v1.Get("/", html_handler)
	v1.Get("/hello", string_handler)

	v2 := engine.Group("/v2")
	v2.Get("/hello/:name", string_handler2) // todo 有点问题待解决
	v2.Post("/login", json_handler)
	engine.Use(xwhttp.Logger())
	v2.Use(middlewareForV2())
	engine.Run(":8080")

}

func string_handler(c *xwhttp.Context) {
	fmt.Println("进入了string_handler")
	c.String(http.StatusOK, "hello [%s], you're at [%s]\n", c.Query("name"), c.Path)
}
func string_handler2(c *xwhttp.Context) {
	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
}
func html_handler(c *xwhttp.Context) {
	c.HTML(http.StatusOK, fmt.Sprintf("<h1>Hello Gee</h1>\n<h2>method=%s</h2>\n<h2>path=%s</h2>", c.Method, c.Path))
}
func json_handler(c *xwhttp.Context) {
	c.JSON(http.StatusOK, xwhttp.H{
		"username": c.PostForm("username"),
		"password": c.PostForm("password"),
	})
}

func middlewareForV2() xwhttp.HandlerFunc {
	return func(c *xwhttp.Context) {
		// Start timer
		t := time.Now()
		// 如果服务器发生错误，返回500状态码，注意这时候并不会panic，下一行日志记录会继续执行
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
