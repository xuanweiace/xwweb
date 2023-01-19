package xwhttp

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_bibao(t *testing.T) {
	x := 1
	f := func() {
		y := x
		fmt.Println(x)
		fmt.Println(y)
	}
	x = 5
	f()
}

type A struct {
	x int
}

func Test_nil(t *testing.T) {
	var s string
	mp := map[string]string{}
	val, ok := mp[s]
	fmt.Println(val, ok)
	val2, ok := mp["abc"]
	fmt.Println(val2, ok)
	fmt.Println(val2 == "") // true
	fmt.Println(s == "")    // true

	//fmt.Println(val2 == nil) // 报错

	var a A
	fmt.Println(a)                      // {0}
	fmt.Println(a == A{}, a == A{x: 0}) // true true
	var p *A
	fmt.Println(p) // <nil>

}

type Func func(int, *A)
type Int int32

type B struct {
	x int
}

func Test_f啊unc(t *testing.T) {
	f := Func(func(i int, a *A) {
		fmt.Println("123")
	})
	fmt.Println(f)
	//不能f()这样调用
	fmt.Println(int32(5)) // go不支持类型自动转换（隐式转换），只支持强制类型转换（显示转换）
	fmt.Println(Int(5))
	a := A{x: 123}
	b := B(a)
	fmt.Println(b)
	//fmt.Println(b == A{x: 123}) 报错
	//fmt.Println(a.(A))
}

func printFunc(a, b int) {
	fmt.Println(a, b)
}
func printFunc2(a ...int) {
	fmt.Println(a)
}
func printFunc3(a int, b string) {
	fmt.Println(a, b)
}
func gen_2int() (int, int) {
	return 1, 2
}
func gen_int_string() (int, string) {
	return 1, "a"
}
func Test_这样传参也可以啊(t *testing.T) {
	printFunc(gen_2int())
	printFunc2(gen_2int())
	//printFunc2(gen_int_string()) 报错
	printFunc3(gen_int_string())
}

func Test_分配在栈上(t *testing.T) {
	x := map[string]string{}
	fmt.Printf("%p\n", x)
	fmt.Printf("%p\n", &x)

	y := 100
	z := 200
	fmt.Printf("%p %p\n", &y, &z)
	//fmt.Printf("%T", y)

}

func Test_for的坑(t *testing.T) {
	var m = []int{1, 2, 3, 4, 5}

	for a, b := range m {
		i, v := a, b
		go func() {
			time.Sleep(time.Second * 3)
			fmt.Println(i, v)
		}()
	}

	time.Sleep(time.Second * 10)

}
func Test_文件修改能否被Read到(t *testing.T) {
	f, err := os.Open("./other_test.go")
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("已经打开")
	time.Sleep(15e9)
	bs := make([]byte, 10000, 10000)
	n, _ := f.Read(bs)

	fmt.Println("n:", n)
	fmt.Println("bs:")
	fmt.Println(bs)
	fmt.Println(string(bs)) // 经测试，程序启动后修改的当前文件，则修改后的内容不会被输出
	//time.AfterFunc()
}

func change(s []string) {
	fmt.Printf("%p\n", s)
	s[0] = "asd"
}
func Test_关于切片复制(t *testing.T) {
	a := []string{"a", "b"}
	fmt.Printf("%p\n", a)
	change(a[:])
	fmt.Println(a)
}
