package logs

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func ErrPrintln(err error, desc string, args ...any) {
	if err == nil {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	f := runtime.CallersFrames([]uintptr{pcs[0]})
	ff, _ := f.Next()
	s := strings.Builder{}
	_, _ = fmt.Fprintf(&s, "%s:%d %s err:[%s]", ff.File, ff.Line, desc, err)
	if len(args) > 0 {
		s.WriteString(" args:")
		for _, arg := range args {
			_, _ = fmt.Fprintf(&s, "%v", arg)
		}
	}
	if err != nil {
		log.Println(s.String())
	}
}
