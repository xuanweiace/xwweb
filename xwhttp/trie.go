package xwhttp

import (
	"strings"
)

// router初始化时一个空node，insert是在该node上添加children
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang。 若为""则代表该node不是完整路由
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// 要求插入的必须是一个规范化的路由，不能是"/a//c/d"这种
// 要求插入的模糊路由不能重叠，比如插入两条路由分别是"/a/:name/"和"/a/:age/"那就寄了，search的时候会按照第一个找到的返回
// n代表一个根节点，在这个节点上进行插入
// 对于非全局函数，尽量不用递归
func (n *node) insert(parts []string) {
	rt := n
	for _, part := range parts {
		found := false
		for _, child := range rt.children {
			if child.part == part {
				found = true
				rt = child // 虽然是在rt的for循环里面修改了rt但是是可以的因为紧接着就退出循环了
				break
			}
		}
		if !found {
			rt.children = append(rt.children, &node{
				pattern:  "",
				part:     part,
				children: make([]*node, 0),
				isWild:   part[0] == ':' || part[0] == '*',
			})
			rt = rt.children[len(rt.children)-1] //这一步别忘了！！！
		}
	}
	//完全匹配结束，则在对应节点填充pattern字段标记这是一个完整路由
	rt.pattern = joinParts(parts)
	//fmt.Printf("rt:%v\n", rt)
}

// parts里肯定放不包含*和: 是精确路由
// 最后只需要返回匹配到的叶子结点就好了 因为可以通过叶子结点保存的信息复现出整个路径
func (n *node) search(parts []string) *node {
	rt := n

	for _, part := range parts {
		//fmt.Printf("root:%q\n", rt)
		//优先找精确路由
		var n0, n1, n2 *node //依次是，精确匹配、模糊匹配、通配符匹配。。。其实完全可以只用一个变量啊 没必要3个变量！！！
		for _, child := range rt.children {
			if child.part == part {
				n0 = child
			} else if child.part[0] == ':' {
				n1 = child
			} else if child.part[0] == '*' {
				n2 = child
			}
		}
		if n0 != nil {
			rt = n0
		} else if n1 != nil {
			rt = n1
		} else if n2 != nil {
			rt = n2
			break // 找到通配符了直接返回，不再向后匹配
		} else {
			return nil
		}
	}
	return rt
}

func parsePattern(pattern string) []string {
	li := strings.Split(pattern, "/")

	parts := make([]string, 0)
	//去掉了所有空路由部分 比如"/a//c/d"这种
	for _, item := range li {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { // 万能通配符，后面的都被匹配进去了
				break
			}
		}
	}
	return parts
}

// 想过用unparse还是assemble，最终还是用join吧，和python保持一致
func joinParts(parts []string) string {
	pattern := "/" + strings.Join(parts, "/")
	return pattern
}
