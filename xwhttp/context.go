package xwhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 定义这个的意义是什么？
type H map[string]interface{}

// 为什么要把返回值也放进来？
// 因为设计思想就是：Context 随着每一个请求的出现而产生，请求的结束而销毁，和当前请求强相关的信息都应由 Context 承载

// 想好Context是建造者模式（先都set好了然后统一build），还是直接可以set值的
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

// 实例化的时候直接返回一个指针。不需要包外暴露，这样可以保证context一定是框架做管理的而不是用户(即用户可以从我context里取东西和放自定义的东西，但是不能滥用与生成)
func newContextInstance(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Writer:     writer,
		Req:        request,
		Path:       request.URL.Path,
		Method:     request.Method,
		StatusCode: 200,
	}
}
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}
func (c *Context) SetStatusCode(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 从req中获取
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 以各种格式返回
func (c Context) String(status int, format string, values ...interface{}) {
	c.SetStatusCode(status)
	//不要这样，期望对所有info类的信息都通过方法去set，做到修改属性的收敛
	//c.StatusCode = status
	c.Writer.Write([]byte(fmt.Sprintf(format, values...))) // 字符串转换成byte数组的形式
}

func (c Context) HTML(statuscode int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatusCode(statuscode)
	c.Writer.Write([]byte(html))
}
func (c Context) JSON(statuscode int, j interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatusCode(statuscode)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(j); err != nil {
		fmt.Fprintf(c.Writer, "后端json-encode错误")
	}
}
