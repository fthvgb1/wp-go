package httptool

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"golang.org/x/exp/constraints"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func GetString(u string, q map[string]any, a ...any) (r string, code int, err error) {
	res, err := Get(u, q, a...)
	if res != nil {
		code = res.StatusCode
	}
	if err != nil {
		return "", code, err
	}
	defer res.Body.Close()
	rr, err := io.ReadAll(res.Body)
	if err != nil {
		return "", code, err
	}
	r = string(rr)
	return
}

func Get(u string, q map[string]any, a ...any) (res *http.Response, err error) {
	cli, req, err := GetClient(u, q, a...)
	res, err = cli.Do(req)
	return
}

func GetToJsonAny[T any](u string, q map[string]any, a ...any) (r T, code int, err error) {
	rr, err := Get(u, q, a...)
	if err != nil {
		return
	}
	code = rr.StatusCode
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		return
	}
	rr.Body.Close()
	err = json.Unmarshal(b, &r)
	return
}

func PostWwwString(u string, form map[string]any, a ...any) (r string, code int, err error) {
	rr, err := Post(u, 1, form, a...)
	if err != nil {
		return "", 0, err
	}
	code = rr.StatusCode
	rs, err := io.ReadAll(rr.Body)
	if err != nil {
		return "", code, err
	}
	rr.Body.Close()
	r = string(rs)
	return
}
func PostFormDataString(u string, form map[string]any, a ...any) (r string, code int, err error) {
	rr, err := Post(u, 2, form, a...)
	if err != nil {
		return "", 0, err
	}
	code = rr.StatusCode
	rs, err := io.ReadAll(rr.Body)
	if err != nil {
		return "", code, err
	}
	rr.Body.Close()
	r = string(rs)
	return
}

func GetClient(u string, q map[string]any, a ...any) (res *http.Client, req *http.Request, err error) {
	parse, err := url.Parse(u)
	if err != nil {
		return nil, nil, err
	}
	cli := http.Client{}
	values := parse.Query()
	setValue(q, values)
	parse.RawQuery = values.Encode()
	req = &http.Request{
		Method: "GET",
		URL:    parse,
	}
	setArgs(&cli, req, a...)
	return &cli, req, err
}

// Post request
//
// u url
//
// types 1 x-www-form-urlencoded, 2 form-data, 3 json, 4 binary
func Post(u string, types int, form map[string]any, a ...any) (res *http.Response, err error) {
	cli, req, err := PostClient(u, types, form, a...)
	res, err = cli.Do(req)
	return
}

func PostClient(u string, types int, form map[string]any, a ...any) (cli *http.Client, req *http.Request, err error) {
	parse, err := url.Parse(u)
	if err != nil {
		return
	}
	cli = &http.Client{}
	req = &http.Request{
		Method: "POST",
		URL:    parse,
		Header: http.Header{},
	}
	switch types {
	case 1:
		values := url.Values{}
		setValue(form, values)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		b := strings.NewReader(values.Encode())
		req.Body = io.NopCloser(b)
	case 2:
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		err = setFormData(form, writer)
		if err != nil {
			return
		}
		err = writer.Close()
		if err != nil {
			return
		}
		req.Body = io.NopCloser(payload)
		req.Header.Add("Content-Type", writer.FormDataContentType())
	case 3:
		fo, err := json.Marshal(form)
		if err != nil {
			return nil, nil, err
		}
		b := bytes.NewReader(fo)
		req.Body = io.NopCloser(b)
		req.Header.Add("Content-Type", "application/json")
	case 4:
		b, ok := maps.GetStrAnyVal[[]byte](form, "binary")
		if !ok {
			return nil, nil, errors.New("no binary value")
		}
		req.Body = io.NopCloser(bytes.NewReader(b))
		req.Header.Add("Content-Type", "application/octet-stream")
	}
	setArgs(cli, req, a...)
	return
}

func PostToJsonAny[T any](u string, types int, form map[string]any, a ...any) (r T, code int, err error) {
	rr, err := Post(u, types, form, a...)
	if err != nil {
		return
	}
	code = rr.StatusCode
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		return
	}
	rr.Body.Close()
	err = json.Unmarshal(b, &r)
	return
}

func setArgs(cli *http.Client, req *http.Request, a ...any) {
	if len(a) < 1 {
		return
	}
	for _, arg := range a {
		h, ok := arg.(map[string]string)
		if ok && len(h) > 0 {
			for k, v := range h {
				req.Header.Add(k, v)
			}
		}
		hh, ok := arg.(http.Header)
		if ok {
			req.Header = hh
		}
		t, ok := arg.(time.Duration)
		if ok {
			cli.Timeout = t
		}
		checkRedirect, ok := arg.(func(req *http.Request, via []*http.Request) error)
		if ok {
			cli.CheckRedirect = checkRedirect
		}
		jar, ok := arg.(http.CookieJar)
		if ok {
			cli.Jar = jar
		}
		c, ok := arg.(string)
		if ok && c != "" {
			req.Header.Add("cookie", c)
		}
	}
}

func set[T constraints.Integer | constraints.Float](a []T, k string, values url.Values) {
	if !strings.Contains(k, "[]") {
		k = str.Join(k, "[]")
	}
	for _, vv := range a {
		values.Add(k, number.ToString(vv))
	}
}

func setFormData(m map[string]any, values *multipart.Writer) (err error) {
	for k, v := range m {
		switch v.(type) {
		case *os.File:
			f := v.(*os.File)
			if f == nil {
				continue
			}
			ff, err := values.CreateFormFile(k, f.Name())
			if err != nil {
				return err
			}
			_, err = io.Copy(ff, f)
			if err != nil {
				return err
			}
		case string:
			err = values.WriteField(k, v.(string))
		case int64, int, int8, int32, int16, uint64, uint, uint8, uint32, uint16, float32, float64:
			err = values.WriteField(k, fmt.Sprintf("%v", v))
		case []string:
			if !strings.Contains(k, "[]") {
				k = str.Join(k, "[]")
			}
			for _, vv := range v.([]string) {
				err = values.WriteField(k, vv)
			}
		case *[]string:
			if !strings.Contains(k, "[]") {
				k = str.Join(k, "[]")
			}
			for _, vv := range *(v.(*[]string)) {
				err = values.WriteField(k, vv)
			}
		case []int64:
			err = sets(v.([]int64), k, values)
		case []int:
			err = sets(v.([]int), k, values)
		case []int8:
			err = sets(v.([]int8), k, values)
		case []int16:
			err = sets(v.([]int16), k, values)
		case []int32:
			err = sets(v.([]int32), k, values)
		case []uint64:
			err = sets(v.([]uint64), k, values)
		case []uint:
			err = sets(v.([]uint), k, values)
		case []uint8:
			err = sets(v.([]uint8), k, values)
		case []uint16:
			err = sets(v.([]uint16), k, values)
		case []uint32:
			err = sets(v.([]uint32), k, values)
		case []float32:
			err = sets(v.([]float32), k, values)
		case []float64:
			err = sets(v.([]float64), k, values)
		}
	}
	return
}

func sets[T constraints.Integer | constraints.Float](a []T, k string, values *multipart.Writer) error {
	if !strings.Contains(k, "[]") {
		k = str.Join(k, "[]")
	}
	for _, vv := range a {
		err := values.WriteField(k, number.ToString(vv))
		if err != nil {
			return err
		}
	}
	return nil
}

func setValue(m map[string]any, values url.Values) {
	for k, v := range m {
		switch v.(type) {
		case string:
			values.Add(k, v.(string))
		case int64, int, int8, int32, int16, uint64, uint, uint8, uint32, uint16, float32, float64:
			values.Add(k, fmt.Sprintf("%v", v))
		case []string:
			if !strings.Contains(k, "[]") {
				k = str.Join(k, "[]")
			}
			for _, vv := range v.([]string) {
				values.Add(k, vv)
			}
		case *[]string:
			if !strings.Contains(k, "[]") {
				k = str.Join(k, "[]")
			}
			for _, vv := range *(v.(*[]string)) {
				values.Add(k, vv)
			}
		case []int64:
			set(v.([]int64), k, values)
		case []int:
			set(v.([]int), k, values)
		case []int8:
			set(v.([]int8), k, values)
		case []int16:
			set(v.([]int16), k, values)
		case []int32:
			set(v.([]int32), k, values)
		case []uint64:
			set(v.([]uint64), k, values)
		case []uint:
			set(v.([]uint), k, values)
		case []uint8:
			set(v.([]uint8), k, values)
		case []uint16:
			set(v.([]uint16), k, values)
		case []uint32:
			set(v.([]uint32), k, values)
		case []float32:
			set(v.([]float32), k, values)
		case []float64:
			set(v.([]float64), k, values)
		}
	}
}
