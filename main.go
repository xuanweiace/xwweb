package main

import (
	"net/http"
	"xwace/xwweb/xwhttp"
)

func main() {
	//standardhttp.Main()
	xwhttp.NewInstance()

}

func string_handler(c *xwhttp.Context) {
	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
}
func html_handler(c *xwhttp.Context) {
	c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
}
func json_handler(c *xwhttp.Context) {
	c.JSON(http.StatusOK, xwhttp.H{
		"username": c.PostForm("username"),
		"password": c.PostForm("password"),
	})
}
