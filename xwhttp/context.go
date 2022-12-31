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
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
}

// 实例化的时候直接返回一个指针。不需要包外暴露，这样可以保证context一定是框架做管理的而不是用户(即用户可以从我context里取东西和放自定义的东西，但是不能滥用与生成)
func newContextInstance(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Writer:     writer,
		Req:        request,
		Path:       request.URL.Path,
		Method:     request.Method,
		StatusCode: 200,
		Params:     map[string]string{},
		handlers:   []HandlerFunc{}, // 一个优化是，这里只是赋值为nil，在ServeHTTP中如果判定有中间件，再分配空间。这样可以减少小
		index:      -1,
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
	//fmt.Printf("c.Req.PostForm:%q\n", c.Req.PostForm)
	//c.Req.ParseForm() todo 有先后顺序的 有小坑
	fmt.Printf("c.Req.PostForm:%q\n", c.Req.Form)
	//fmt.Printf("c.Req.PostForm:%q\n", c.Req.)
	//fmt.Printf("c.Req.PostFormValue[%s]:%q\n", key, c.Req.PostFormValue(key))

	return c.Req.FormValue(key)
}
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// 以各种格式返回
func (c *Context) String(status int, format string, values ...interface{}) {
	c.SetStatusCode(status)
	//不要这样，期望对所有info类的信息都通过方法去set，做到修改属性的收敛
	//c.StatusCode = status
	c.Writer.Write([]byte(fmt.Sprintf(format, values...))) // 字符串转换成byte数组的形式
}

func (c *Context) HTML(statuscode int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatusCode(statuscode)
	c.Writer.Write([]byte(html))
}
func (c *Context) JSON(statuscode int, j interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatusCode(statuscode)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(j); err != nil {
		fmt.Fprintf(c.Writer, "后端json-encode错误")
	}
}

// 注意c.index++ 不能去掉，就会一直卡在第一个中间件上了，必然死循环。
// 可以这样想：每调用一次 Next()，c.index 得 +1，不然 c.index 就会一直是 0。
// 死循环之后，是不会走到 for 语句后面的 c.index++ 的。
func (c *Context) Next() {
	c.index++
	//为什么这个需要遍历？每个middleware都调用一下next不就好了吗？
	//不是所有的handler都会调用 Next()。手工调用 Next()，一般用于在请求前后各实现一些行为。
	//如果中间件只作用于请求前，可以省略调用Next()，算是一种兼容性比较好的写法。
	//又因为index和handlers都是context的属性，在context的视角看就是全局变量，因此调用多次Next也不会产生重复调用handler的情况
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}
