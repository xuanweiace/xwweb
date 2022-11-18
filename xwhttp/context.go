package xwhttp

import (
	"fmt"
	"net/http"
)

// 为什么要把返回值也房间来？
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

func (c Context) String(status int, format string, values ...interface{}) {
	c.StatusCode = status
	c.Writer.Write([]byte(fmt.Sprintf(format, values...))) // 字符串转换成byte数组的形式
}
