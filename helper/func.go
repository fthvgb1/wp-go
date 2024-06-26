package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func ToAny[T any](v T) any {
	return v
}

func Or[T any](is bool, left, right T) T {
	if is {
		return left
	}
	return right
}

func StructColumnToSlice[T any, M any](arr []M, field string) (r []T) {
	for i := 0; i < len(arr); i++ {
		v := reflect.ValueOf(arr[i]).FieldByName(field).Interface()
		if val, ok := v.(T); ok {
			r = append(r, val)
		}
	}
	return
}

func UrlScheme(u string, isHttps bool) string {
	return Or(isHttps,
		strings.Replace(u, "http://", "https://", 1),
		strings.Replace(u, "https://", "http://", 1),
	)
}

func CutUrlHost(u string) string {
	ur, err := url.Parse(u)
	if err != nil {
		return u
	}
	ur.Scheme = ""
	ur.Host = ""
	return ur.String()
}

func Defaults[T comparable](vals ...T) T {
	var val T
	for _, v := range vals {
		if v != val {
			return v
		}
	}
	return val
}
func DefaultVal[T any](v, defaults T) T {
	var zero T
	if reflect.DeepEqual(zero, v) {
		return defaults
	}
	return v
}

func IsZero[T comparable](t T) bool {
	var vv T
	return vv != t
}
func IsZeros(v any) bool {
	switch v.(type) {
	case int64, int, int8, int16, int32, uint64, uint, uint8, uint16, uint32:
		i := fmt.Sprintf("%d", v)
		return str.ToInt[int64](i) == 0
	case float32, float64:
		f := fmt.Sprintf("%v", v)
		ff, _ := strconv.ParseFloat(f, 64)
		return ff == float64(0)
	case bool:
		return v.(bool) == false
	case string:
		s := v.(string)
		return s == ""
	}
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

func ToBool[T comparable](t T) bool {
	v := any(t)
	switch v.(type) {
	case string:
		s := v.(string)
		return s != "" && s != "0"
	}
	var vv T
	return vv != t
}

func ToBoolInt(t any) int8 {
	if IsZeros(t) {
		return 0
	}
	return 1
}

func GetContextVal[V, K any](ctx context.Context, k K, defaults V) V {
	v := ctx.Value(k)
	if v == nil {
		return defaults
	}
	vv, ok := v.(V)
	if !ok {
		return defaults
	}
	return vv
}

func IsImplements[T, A any](i A) (T, bool) {
	var a any = i
	t, ok := a.(T)
	return t, ok
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsFile(file string) bool {
	info, err := os.Stat(file)
	return err == nil && !info.IsDir()
}

func GetAnyVal[T any](v any, defaults T) T {
	vv, ok := v.(T)
	if !ok {
		return defaults
	}
	return vv
}

func ParseArgs[T any](defaults T, a ...any) T {
	for _, aa := range a {
		v, ok := aa.(T)
		if ok {
			return v
		}
	}
	return defaults
}

func RunFnWithTimeout(ctx context.Context, t time.Duration, call func(), a ...any) (err error) {
	ctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()
	done := make(chan struct{}, 1)
	go func() {
		call()
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		msg := ParseArgs("", a...)
		if msg != "" {
			return errors.New(str.Join(msg, ":", ctx.Err().Error()))
		}
		return ctx.Err()
	case <-done:
		close(done)
	}
	return nil
}

func RunFnWithTimeouts[A, V any](ctx context.Context, t time.Duration, ar A, call func(A) (V, error), a ...any) (v V, err error) {
	ctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()
	done := make(chan struct{}, 1)
	go func() {
		v, err = call(ar)
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		msg := ParseArgs("", a...)
		if msg != "" {
			return v, errors.New(str.Join(msg, ":", ctx.Err().Error()))
		}
		return v, ctx.Err()
	case <-done:
		close(done)
	}
	return v, err
}

func JsonDecode[T any](byts []byte) (T, error) {
	var v T
	err := json.Unmarshal(byts, &v)
	return v, err
}

func AsError[T any](err error) (T, bool) {
	var v T
	ok := errors.As(err, &v)
	return v, ok
}

func IsDirExistAndMkdir(dir string, perm os.FileMode) error {
	stat, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} else {
			return os.MkdirAll(dir, perm)
		}
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is exist but not dir", dir)
	}
	return nil
}

func ReadDir(dir string) ([]string, error) {
	fii, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range fii {
		name := entry.Name()
		if name == "." || name == ".." {
			continue
		}
		files = append(files, filepath.Join(dir, name))
	}
	return files, nil
}
