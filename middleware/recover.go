package middleware

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic:%+v stack:%s", err, stack())
		}
	}()

	c.Next()
}

func stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}
