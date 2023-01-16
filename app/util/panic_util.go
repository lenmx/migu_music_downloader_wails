package util

import (
	"bytes"
	"fmt"
	"runtime"
)

//func PanicCatch(callback func(err string)) {
//	if _err := recover(); _err != nil && callback != nil {
//		fmt.Println("catch err: ", _err)
//		callback(panicTrace(_err))
//	}
//}

func PanicTrace(err interface{}) string {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "%v\n", err)

	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	}

	return buf.String()
}
