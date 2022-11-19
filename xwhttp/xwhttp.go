package xwhttp

import (
	"net/http"
)

/*
对外提供几个方法:
NewInstance Get Post

*/

// httpServer的执行实体
type Engine struct {
	router *router
}

// handler的集合
type router struct {
	handlers map[string]HandlerFunc
}

func (r *router) handle(c *Context) {
	assemble_mapper_key := c.Method + "-" + c.Path
	println("thisis:", assemble_mapper_key)
	if handler, ok := r.handlers[assemble_mapper_key]; ok {
		print("找到了")
		handler(c)
	}
}

// handler定义
type HandlerFunc func(c *Context)

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//new一个Context，包装所有的请求生命周期所需要的对象，后续所有操作都直接用context就可以了。
	c := newContextInstance(writer, request)
	e.router.handle(c)

}

// 创建一个xwHttpServer实例。工厂方法，因此不设置接收器。相当于一个全局函数
// 建议返回的是一个指针，避免把这个instance到处传递的时候都生成的是新的实例。
// 更一般的，最好涉及到工厂方法的类，都用指针去新建，这样可以控制所有的实例都在工厂方法中生成，而不会你随便传个参数就新建实例了。
func NewInstance() *Engine {
	return &Engine{router: &router{handlers: map[string]HandlerFunc{}}}
}

//对于这种结构体里面只包含一个属性的这种情况，其实一级就够了，不需要分成两个层级。

// 用户友好向的API设计，其实就是拆分函数方便用户调用
func (e *Engine) Get(route string, handlerFunc HandlerFunc) {
	assemble_mapper_key := "GET-" + route
	if _, ok := e.router.handlers[assemble_mapper_key]; ok {
		panic("已经有该路由了，无法再次注册，路由为：" + assemble_mapper_key) // 一定要把上下文信息输出出来，不然不好定位错误
	}
	e.router.handlers[assemble_mapper_key] = handlerFunc
}

func (e *Engine) Post(route string, handlerFunc HandlerFunc) {
	assemble_mapper_key := "POST-" + route
	println(assemble_mapper_key)
	if _, ok := e.router.handlers[assemble_mapper_key]; ok {
		panic("已经有该路由了，无法再次注册，路由为：" + assemble_mapper_key) // 一定要把上下文信息输出出来，不然不好定位错误
	}
	e.router.handlers[assemble_mapper_key] = handlerFunc
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}
