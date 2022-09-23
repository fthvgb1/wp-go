package helper

import (
	"crypto/md5"
	"fmt"
	"github.com/dlclark/regexp2"
	"io"
	"reflect"
	"regexp"
	"strings"
)

func IsContainInArr[T comparable](a T, arr []T) bool {
	for _, v := range arr {
		if a == v {
			return true
		}
	}
	return false
}

func StructColumn[T any, M any](arr []M, field string) (r []T) {
	for i := 0; i < len(arr); i++ {
		v := reflect.ValueOf(arr[i]).FieldByName(field).Interface()
		if val, ok := v.(T); ok {
			r = append(r, val)
		}
	}
	return
}

func RangeSlice[T ~int | ~uint | ~int64 | ~int8 | ~int16 | ~int32 | ~uint64](start, end, step T) []T {
	r := make([]T, 0, int((end-start+1)/step+1))
	for i := start; i <= end; {
		r = append(r, i)
		i = i + step
	}
	return r
}

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

func SlicePagination[T any](arr []T, page, pageSize int) []T {
	start := (page - 1) * pageSize
	l := len(arr)
	if start > l {
		start = l
	}
	end := page * pageSize
	if l < end {
		end = l
	}
	return arr[start:end]
}

func StringMd5(str string) string {
	h := md5.New()
	_, err := io.WriteString(h, str)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

var allHtmlTag = regexp.MustCompile("</?.*>")

func StripTags(str, allowable string) string {
	html := ""
	if allowable == "" {
		return allHtmlTag.ReplaceAllString(str, "")
	}
	r := strings.Split(allowable, ">")
	re := ""
	for _, reg := range r {
		if reg == "" {
			continue
		}
		tag := strings.TrimLeft(reg, "<")
		ree := fmt.Sprintf(`%s|\/%s`, tag, tag)
		re = fmt.Sprintf("%s|%s", re, ree)
	}
	ree := strings.Trim(re, "|")
	reg := fmt.Sprintf("<(?!%s).*?>", ree)
	compile, err := regexp2.Compile(reg, regexp2.IgnoreCase)
	if err != nil {
		return str
	}
	html, err = compile.Replace(str, "", 0, -1)
	if err != nil {
		return str
	}
	return html
}
