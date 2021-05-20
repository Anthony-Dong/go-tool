package utils

import (
	"fmt"
	"io"
	"runtime"
)

func PanicPrint(recoverError interface{}, writer io.Writer) {
	stackBuffer := make([]byte, 64<<10) // 最多打印64k的堆栈信息
	if _, err := fmt.Fprintf(writer, "panic: %v\n\n%s\n", recoverError, stackBuffer[:runtime.Stack(stackBuffer, false)]); err != nil {
		fmt.Printf("panic: %v\n\n%s\n", recoverError, stackBuffer[:runtime.Stack(stackBuffer, false)]) // 出现异常用 std 输出
	}
}
