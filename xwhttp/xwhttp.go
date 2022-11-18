package xwhttp

import (
	"net/http"
)

/*
对外提供几个方法:
NewInstance Get Post

*/

// httpServer的执行实体
type xw struct {
	handlers HandlerContainer
}
type Xw = xw //用别名的防守避免使用反模式，不然会有Warning: Exported function with the unexported return type 从导出函数(exported function)返回未导出类型(unexported type)的值

// handler的集合
type HandlerContainer struct {
	mapper map[string]HandlerFunc //对GET和POST算不同的路由，要注册两次handler
} //handler定义
type HandlerFunc func(response http.ResponseWriter, r *http.Request)

func (c HandlerContainer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	assemble_mapper_key := request.Method + "-" + request.URL.Path
	for k, v := range c.mapper {
		if assemble_mapper_key == k {
			v(writer, request)
		}
	}
}

// 创建一个xwHttpServer实例。工厂方法，因此不设置接收器。相当于一个全局函数
// 建议返回的是一个指针，避免把这个instance到处传递的时候都生成的是新的实例。
// 更一般的，最好涉及到工厂方法的类，都用指针去新建，这样可以控制所有的实例都在工厂方法中生成，而不会你随便传个参数就新建实例了。
func NewInstance() *Xw {
	return &xw{handlers: HandlerContainer{mapper: map[string]HandlerFunc{}}}
}

//对于这种结构体里面只包含一个属性的这种情况，其实一级就够了，不需要分成两个层级。

// 用户友好向的API设计，其实就是拆分函数方便用户调用
func (s xw) Get(route string, handlerFunc HandlerFunc) {
	assemble_mapper_key := "GET-" + route
	if _, ok := s.handlers.mapper[assemble_mapper_key]; ok {
		panic("已经有该路由了，无法再次注册，路由为：" + assemble_mapper_key) // 一定要把上下文信息输出出来，不然不好定位错误
	}
	s.handlers.mapper[assemble_mapper_key] = handlerFunc
}

func (s xw) Post(route string, handlerFunc HandlerFunc) {
	assemble_mapper_key := "POST-" + route
	if _, ok := s.handlers.mapper[assemble_mapper_key]; ok {
		panic("已经有该路由了，无法再次注册，路由为：" + assemble_mapper_key) // 一定要把上下文信息输出出来，不然不好定位错误
	}
	s.handlers.mapper[assemble_mapper_key] = handlerFunc
}

func (s xw) Run(addr string) {
	http.ListenAndServe(addr, s.handlers)
}
