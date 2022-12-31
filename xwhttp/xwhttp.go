package xwhttp

import (
	"log"
	"net/http"
	"strings"
)

/*
对外提供几个方法:
NewInstance Get Post

*/

// httpServer的执行实体
type Engine struct {
	*RouteGroup // 居然可以直接e.RouteGroup=&RouteGroup{}这样操作。。。
	router      *router
	groups      []*RouteGroup
}
type RouteGroup struct {
	engine      *Engine
	prefix      string
	middlewares []HandlerFunc
}

// handler定义
type HandlerFunc func(c *Context)

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//new一个Context，包装所有的请求生命周期所需要的对象，后续所有操作都直接用context就可以了。
	c := newContextInstance(writer, request)
	for _, group := range e.groups {
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			c.handlers = append(c.handlers, group.middlewares...)
		}
	}
	e.router.handle(c)
}

// 创建一个xwHttpServer实例。工厂方法，因此不设置接收器。相当于一个全局函数
// 建议返回的是一个指针，避免把这个instance到处传递的时候都生成的是新的实例。
// 更一般的，最好涉及到工厂方法的类，都用指针去新建，这样可以控制所有的实例都在工厂方法中生成，而不会你随便传个参数就新建实例了。
func NewInstance() *Engine {
	e := &Engine{router: newRouter()}
	rg := &RouteGroup{
		engine: e,
		prefix: "",
	}
	e.groups = []*RouteGroup{rg}
	e.RouteGroup = rg // 别忘了这个。。。初始化现在变得复杂了很多
	return e
}

//如果只实现到动态路由这一步，对于engine里只有router的这种，结构体里面只包含一个属性的这种情况，其实一级就够了，不需要分成两个层级。（即直接engine.addRoute()，而不用engine.router.addRoute()这样）

// 用户友好向的API设计，其实就是拆分函数方便用户调用
func (g *RouteGroup) Get(route string, handlerFunc HandlerFunc) {
	route_path := g.prefix + route
	if route_path != "/" {
		route_path = strings.TrimSuffix(route_path, "/") // 避免group里设置"/"这种默认路由，存到我们的router里规范是不能带末尾"/"的
	}
	log.Printf("Route GET - %s", route_path)
	g.engine.router.addRoute("GET", route_path, handlerFunc)
}

func (g *RouteGroup) Post(route string, handlerFunc HandlerFunc) {
	route_path := g.prefix + route
	log.Printf("Route POST - %s", route_path)
	g.engine.router.addRoute("POST", route_path, handlerFunc)
}

func (g *RouteGroup) Group(route string) *RouteGroup {
	rg := &RouteGroup{
		engine: g.engine,
		prefix: g.prefix + route,
	}
	rg.engine.groups = append(rg.engine.groups, rg)
	return rg
}

func (g *RouteGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}
