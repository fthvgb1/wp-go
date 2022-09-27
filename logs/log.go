package logs

import (
	"log"
	"strings"
)

func ErrPrintln(err error, desc string, args ...any) {
	s := strings.Builder{}
	tmp := "%s err:[%s]"
	if desc == "" {
		tmp = "%s%s"
	}
	s.WriteString(tmp)
	argss := []any{desc, err}
	if len(args) > 0 {
		s.WriteString(" args:")
		for _, arg := range args {
			s.WriteString("%v ")
			argss = append(argss, arg)
		}
	}
	if err != nil {
		log.Printf(s.String(), argss...)
	}
}
