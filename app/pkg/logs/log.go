package logs

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/safety"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

var logs = safety.NewVar[*log.Logger](nil)
var logFile = safety.NewVar[*os.File](nil)

func InitLogger() error {
	preFD := logFile.Load()
	l := &log.Logger{}
	c := config.GetConfig()
	if c.LogOutput == "" {
		c.LogOutput = "stderr"
	}
	var out io.Writer
	switch c.LogOutput {
	case "stdout":
		out = os.Stdout
	case "stderr":
		out = os.Stderr
	default:
		file, err := os.OpenFile(c.LogOutput, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			return err
		}
		out = file
		logFile.Store(file)
	}
	logs.Store(l)
	if preFD != nil {
		_ = preFD.Close()
	}
	l.SetFlags(log.Ldate | log.Ltime)
	l.SetOutput(out)
	return nil
}

func Errs(err error, depth int, desc string, args ...any) {
	var pcs [1]uintptr
	runtime.Callers(depth, pcs[:])
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
	logs.Load().Println(s.String())
}

func Error(err error, desc string, args ...any) {
	Errs(err, 3, desc, args...)
}

func IfError(err error, desc string, args ...any) {
	if err == nil {
		return
	}
	Errs(err, 3, desc, args...)
}
