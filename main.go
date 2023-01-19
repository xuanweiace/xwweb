package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"xwace/xwweb/xwhttp"
)

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
func main() {
	//standardhttp.Main()
	engine := xwhttp.Default()

	engine.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	engine.LoadHTMLGlob("templates/*")
	engine.Get("/", html_handler_wrap("css.tmpl", nil))
	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	engine.Get("/students", html_handler_wrap("arr.tmpl", xwhttp.H{
		"key1": "key2",
		"stu":  []*student{stu1, stu2},
	}))
	engine.Get("/date", html_handler_wrap("custom_func.tmpl", xwhttp.H{
		"key3": "key4",
		"now":  time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
	}))
	engine.Get("/get_json", json_handler)
	//engine.Get("/get_post_html", html_handler)
	//engine.Post("/get_post_html", html_handler)
	engine.Get("/get_string", string_handler)

	v1 := engine.Group("/v1")
	//v1.Get("/", html_handler)
	v1.Get("/hello", string_handler)

	v2 := engine.Group("/v2")
	v2.Get("/hello/:name", string_handler2) // todo 有点问题待解决
	v2.Post("/login", json_handler)
	engine.Use(xwhttp.Logger())
	v2.Use(middlewareForV2())

	engine.Static("/asset", "./static")
	fmt.Println(os.Getwd())

	err := engine.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}

}

func string_handler(c *xwhttp.Context) {
	fmt.Println("进入了string_handler")
	names := []string{"testsae"}
	c.String(http.StatusOK, "hello [%s], you're at %s tmp:%s\n", c.Param("name"), c.Path, names[100])
}
func string_handler2(c *xwhttp.Context) {
	c.String(http.StatusOK, "hello %s, you're at %s \n", c.Param("name"), c.Path)
}

func html_handler_wrap(name string, data any) xwhttp.HandlerFunc {
	fmt.Println("测试地址是否一样1：", &name, &data)
	html_handler := func(c *xwhttp.Context) {
		fmt.Println("测试地址是否一样2：", &name, &data)
		c.HTML(http.StatusOK, name, data) // 疑问 这样用闭包，参数会被覆盖成一样吗？吗？
	}
	return html_handler
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
