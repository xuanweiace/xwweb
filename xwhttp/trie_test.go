package xwhttp

import (
	"fmt"
	"testing"
)

func Test_parse_route(t *testing.T) {
	//*的作用
	fmt.Println(parsePattern("/*"))
	fmt.Println(parsePattern("/*/a"))

	//
	fmt.Println(parsePattern("/a"))
	fmt.Println(parsePattern("/a/"))
	fmt.Println(parsePattern("//a/")) //过滤不掉这种非法情况
	fmt.Println(parsePattern("/a/a"))
	fmt.Println(parsePattern("/:name"))
	fmt.Println(parsePattern("/:name:abc/age"))
}
