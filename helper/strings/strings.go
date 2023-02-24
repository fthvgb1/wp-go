package strings

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/exp/constraints"
	"io"
	"strconv"
	"strings"
)

func Join(s ...string) (str string) {
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

func ToInteger[T constraints.Integer](s string, defaults T) T {
	if s == "" {
		return defaults
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaults
	}
	return T(i)
}
func ToInt[T constraints.Integer](s string) T {
	defaults := T(0)
	if s == "" {
		return defaults
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaults
	}
	return T(i)
}

func Md5(str string) string {
	h := md5.New()
	_, err := io.WriteString(h, str)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func BuilderJoin(s *strings.Builder, str ...string) {
	for _, ss := range str {
		s.WriteString(ss)
	}
}

func BuilderFormat(s *strings.Builder, format string, args ...any) {
	s.WriteString(fmt.Sprintf(format, args...))
}

type Builder struct {
	*strings.Builder
}

func NewBuilder() *Builder {
	return &Builder{&strings.Builder{}}
}

func (b *Builder) WriteString(s ...string) (count int) {
	for _, ss := range s {
		i, _ := b.Builder.WriteString(ss)
		count += i
	}
	return
}
func (b *Builder) Sprintf(format string, a ...any) int {
	i, _ := fmt.Fprintf(b, format, a...)
	return i
}

// CutSpecialDuplicate '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP)
func CutSpecialDuplicate(s, char string) string {
	return strings.Join(strings.Fields(s), char)
}
