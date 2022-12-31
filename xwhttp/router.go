package xwhttp

import (
	"fmt"
	"net/http"
	"strings"
)

// 好像也可以因为加了动态路由就去掉了handler，对于模糊匹配也可以支持。
type router struct {
	handlers   map[string]HandlerFunc
	get_roots  *node // todo
	post_roots *node
}

func newRouter() *router {
	return &router{
		handlers:   map[string]HandlerFunc{},
		get_roots:  &node{children: make([]*node, 0)}, //注意这里不能是nil
		post_roots: &node{children: make([]*node, 0)},
	}
}
func (r *router) handle(c *Context) {
	fmt.Println("c.Path:", c.Path)
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		//可以直接赋值或者copy的，如果要这么写的话，需要在newContextInstance的时候make()才可以
		for k, v := range params {
			c.Params[k] = v
		}
		assemble_mapper_key := c.Method + "-" + n.pattern
		fmt.Println("当前请求的路径是:", assemble_mapper_key)
		if handler, ok := r.handlers[assemble_mapper_key]; ok {
			handler(c)
		} else {
			fmt.Println("没找着", assemble_mapper_key)
		}
	} else {
		//这里也可以用 fmt.Fprintf(c.Writer, "404 NOT FOUND: %s\n", c.Path) 但是不推荐，因为没法返回状态码。因此推荐所有的response操作都通过context提供的方法执行
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

func (r *router) addRoute(method string, pattern string, handlerFunc HandlerFunc) {
	parts := parsePattern(pattern)
	//fmt.Println("parts:", parts)
	assemble_mapper_key := method + "-" + pattern
	//fmt.Println("add:assemble_mapper_key=", assemble_mapper_key)
	if _, ok := r.handlers[assemble_mapper_key]; ok {
		panic("已经有该路由了，无法再次注册，路由为：" + assemble_mapper_key) // 一定要把上下文信息输出出来，不然不好定位错误
	}
	r.handlers[assemble_mapper_key] = handlerFunc
	if method == "GET" {
		r.get_roots.insert(parts)
	} else if method == "POST" {
		r.post_roots.insert(parts)
	} else {
		panic("暂不支持该HTTP method：" + method)
	}
}

func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	params := make(map[string]string)
	realParts := parsePattern(pattern)
	fmt.Println("realPartsL:", realParts)
	var n *node
	if method == "GET" {
		n = r.get_roots.search(realParts)
		fmt.Printf("node:%q\n", n)

	} else if method == "POST" {
		n = r.post_roots.search(realParts)
	} else {
		return nil, nil
	}
	if n != nil {
		parts := parsePattern(n.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = realParts[i]
			} else if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(realParts[i:], "/")
				break
			}
		}
		//fmt.Printf("node:%q par:%q", n, params)
		return n, params
	}
	return nil, nil
}
