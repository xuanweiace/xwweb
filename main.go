package main

import (
	"fmt"
	"net/http"
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
	engine.Run(":8080")
}

func string_handler(c *xwhttp.Context) {
	c.String(http.StatusOK, "hello [%s], you're at [%s]\n", c.Query("name"), c.Path)
}
func html_handler(c *xwhttp.Context) {
	c.HTML(http.StatusOK, fmt.Sprintf("<h1>Hello Gee</h1>\n<h2>method=%s</h2>", c.Method))
}
func json_handler(c *xwhttp.Context) {
	c.JSON(http.StatusOK, xwhttp.H{
		"username": c.PostForm("username"),
		"password": c.PostForm("password"),
	})
}
