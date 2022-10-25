package helper

import (
	"crypto/md5"
	"fmt"
	"github.com/dlclark/regexp2"
	"io"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type IntNumber interface {
	~int | ~int64 | ~int32 | ~int8 | ~int16 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func ToAny[T any](v T) any {
	return v
}

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

func RangeSlice[T IntNumber](start, end, step T) []T {
	if step == 0 {
		panic("step can't be 0")
	}
	l := int((end-start+1)/step + 1)
	if l < 0 {
		l = 0 - l
	}
	r := make([]T, 0, l)
	for i := start; ; {
		r = append(r, i)
		i = i + step
		if (step > 0 && i > end) || (step < 0 && i < end) {
			break
		}
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

var tag = regexp.MustCompile(`<(.*?)>`)

func StripTagsX(str, allowable string) string {
	if allowable == "" {
		return allHtmlTag.ReplaceAllString(str, "")
	}
	tags := tag.ReplaceAllString(allowable, "$1|")
	or := strings.TrimRight(tags, "|")
	reg := fmt.Sprintf(`<(/?(%s).*?)>`, or)
	regx := fmt.Sprintf(`\{\[(/?(%s).*?)\]\}`, or)
	cp, err := regexp.Compile(reg)
	if err != nil {
		return str
	}
	rep := cp.ReplaceAllString(str, "{[$1]}")
	tmp := tag.ReplaceAllString(rep, "")
	rex, err := regexp.Compile(regx)
	if err != nil {
		return str
	}
	html := rex.ReplaceAllString(tmp, "<$1>")
	return html
}

var tagx = regexp.MustCompile(`(</?[a-z0-9]+?)( |>)`)
var selfCloseTags = map[string]string{"area": "", "base": "", "basefont": "", "br": "", "col": "", "command": "", "embed": "", "frame": "", "hr": "", "img": "", "input": "", "isindex": "", "link": "", "meta": "", "param": "", "source": "", "track": "", "wbr": ""}

func CloseHtmlTag(str string) string {
	tags := tag.FindAllString(str, -1)
	if len(tags) < 1 {
		return str
	}
	var tagss = make([]string, 0, len(tags))
	for _, s := range tags {
		ss := strings.TrimSpace(tagx.FindString(s))
		if ss[len(ss)-1] != '>' {
			ss = fmt.Sprintf("%s>", ss)
			if _, ok := selfCloseTags[ss]; ok {
				continue
			}
		}
		tagss = append(tagss, ss)
	}
	r := SliceMap(SliceReverse(ClearClosedTag(tagss)), func(s string) string {
		return fmt.Sprintf("</%s>", strings.Trim(s, "<>"))
	})
	return strings.Join(r, "")
}

func ClearClosedTag(s []string) []string {
	i := 0
	for {
		if len(s[i:]) < 2 {
			return s
		}
		l := s[i]
		r := fmt.Sprintf(`</%s>`, strings.Trim(l, "<>"))
		if s[i+1] == r {
			if len(s[i+1:]) > 1 {
				ss := s[:i]
				s = append(ss, s[i+2:]...)

			} else {
				s = s[:i]
			}
			i = 0
			continue
		}
		i++
	}
}

func SliceMap[T, R any](arr []T, fn func(T) R) []R {
	r := make([]R, 0, len(arr))
	for _, t := range arr {
		r = append(r, fn(t))
	}
	return r
}

func SliceFilter[T any](arr []T, fn func(T) bool) []T {
	var r []T
	for _, t := range arr {
		if fn(t) {
			r = append(r, t)
		}
	}
	return r
}

func SliceReduce[T, R any](arr []T, fn func(T, R) R, r R) R {
	for _, t := range arr {
		r = fn(t, r)
	}
	return r
}

func SliceReverse[T any](arr []T) []T {
	var r []T
	for i := len(arr); i > 0; i-- {
		r = append(r, arr[i-1])
	}
	return r
}

func SliceSelfReverse[T any](arr []T) []T {
	l := len(arr)
	half := l / 2
	for i := 0; i < half; i++ {
		arr[i], arr[l-i-1] = arr[l-i-1], arr[i]
	}
	return arr
}

func SimpleSliceToMap[K comparable, V any](arr []V, fn func(V) K) map[K]V {
	return SliceToMap(arr, func(v V) (K, V) {
		return fn(v), v
	}, true)
}

func SliceToMap[K comparable, V, T any](arr []V, fn func(V) (K, T), isCoverPrev bool) map[K]T {
	m := make(map[K]T)
	for _, v := range arr {
		k, r := fn(v)
		if !isCoverPrev {
			if _, ok := m[k]; ok {
				continue
			}
		}
		m[k] = r
	}
	return m
}

func RandNum[T IntNumber](start, end T) T {
	end++
	return T(rand.Int63n(int64(end-start))) + start
}

type anyArr[T any] struct {
	data []T
	fn   func(i, j T) bool
}

func (r anyArr[T]) Len() int {
	return len(r.data)
}

func (r anyArr[T]) Swap(i, j int) {
	r.data[i], r.data[j] = r.data[j], r.data[i]
}

func (r anyArr[T]) Less(i, j int) bool {
	return r.fn(r.data[i], r.data[j])
}

func SampleSort[T any](arr []T, fn func(i, j T) bool) {
	slice := anyArr[T]{
		data: arr,
		fn:   fn,
	}
	sort.Sort(slice)
	return
}
