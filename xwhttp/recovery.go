package xwhttp

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func trace(msg string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller 即Callers，trace，defer的func三层。注意Recovery函数不算一层。因为在编译时他本身不算做一个函数了。

	var str strings.Builder
	str.WriteString(msg + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%s", err)
				log.Printf("%s\n-----------------\n", trace(msg))
				c.Fail(http.StatusInternalServerError, msg)
			}
		}()

		c.Next() // 因为是中间件的形式，所以别忘了手动调用c.Next，这样才能实现调用defer。如果不是手动调用，则无法调用到这里的defer，劲儿recovery不生效了。
	}

}
