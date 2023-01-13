package helper

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

func StrJoin(s ...string) (str string) {
	if len(s) == 1 {
		return s[0]
	} else if len(s) > 1 {
		b := strings.Builder{}
		for _, s2 := range s {
			b.WriteString(s2)
		}
		str = b.String()
	}
	return
}

func StringMd5(str string) string {
	h := md5.New()
	_, err := io.WriteString(h, str)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
