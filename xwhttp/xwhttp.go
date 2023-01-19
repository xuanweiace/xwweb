package xwhttp

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

/*
对外提供几个方法:
NewInstance Get Post

*/

// httpServer的执行实体
type Engine struct {
	*RouteGroup  // 居然可以直接e.RouteGroup=&RouteGroup{}这样操作。。。
	router       *router
	groups       []*RouteGroup
	htmlTemplete *template.Template // html模板渲染
	funcMap      template.FuncMap   // html模板渲染 注意是个map
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
	c := newContextInstance(writer, request, e)
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

// 其实相当于我们g.Get注册路由之前做一些处理
// 调用类似：e.Static("/assets", "D:\mystudy\go\src\xwace\static")
// 注意 这段函数执行的时候，请求还没有到来（甚至web服务器都没有启动），这也就是生命周期的概念，即执行这段代码的时候，要判断是否在req的生命周期内，也就是python的with的意思
func (g *RouteGroup) Static(urlRelativePath string, serverPath string) {
	// 得到完整url
	urlAbsolutePath := path.Join(g.prefix + urlRelativePath)
	//文件系统服务器，类似我们Engine是一个web服务器一样(即只要是http.Handler接口，我们都可以让他成为服务器来接管http请求，通过ServeHTTP方法)
	//换句话说，http.Handler接口可以是一个调用链，都通过内部不断调用ServeHTTP来完成，就和我们的c.Next()调用中间件一样
	//（注意，xwhttp是一个web框架，其中Engine是我们的server引擎，因为它实现了ServeHTTP方法。只不过我们没有显示调用它，而是注册到go原生的http里，在请求到来的时候让原生http框架帮我们调用。现在请求被我们接管，来到我们这边我们需要显示调用）
	fs := http.Dir(serverPath)
	fileServer := http.StripPrefix(urlAbsolutePath, http.FileServer(fs))
	//注意这个函数不能拿出去单独写了，因为没有fileServer这个变量给我们用（这就是闭包的作用，相当于免费的全局变量用）
	handler := func(c *Context) {
		filepath := c.Param("filepath")
		fmt.Println("filepath:", filepath)
		fmt.Println("fs:", fs)
		if _, err := fs.Open(filepath); err != nil {
			fmt.Println("不存在啊！！！")
			c.String(http.StatusNotFound, "err=%v", err)
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}

	urlPattern := path.Join(urlRelativePath, "/*filepath")
	// Register GET handlers
	g.Get(urlPattern, handler)
}

// 关于html模板渲染
func (e *Engine) SetFuncMap(mp template.FuncMap) {
	e.funcMap = mp
}
func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplete = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

func (e *Engine) Run(addr string) error {
	err := http.ListenAndServe(addr, e)
	return err
}
