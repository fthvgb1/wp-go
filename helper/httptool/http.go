package httptool

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

func GetString(u string, q map[string]string, timeout int64, a ...any) (r string, code int, err error) {
	res, err := Get(u, q, timeout, a...)
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

func Get(u string, q map[string]string, timeout int64, a ...any) (res *http.Response, err error) {
	parse, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	cli := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	values := parse.Query()
	for k, v := range q {
		values.Add(k, v)
	}
	parse.RawQuery = values.Encode()
	req := http.Request{
		Method: "GET",
		URL:    parse,
	}
	if len(a) > 0 {
		for _, arg := range a {
			h, ok := arg.(map[string]string)
			if ok && len(h) > 0 {
				for k, v := range h {
					req.Header.Add(k, v)
				}
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
		}
	}
	res, err = cli.Do(&req)
	return
}
