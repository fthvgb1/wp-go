package actions

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/vars"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

func PostComment(c *gin.Context) {
	jar, _ := cookiejar.New(nil)
	cli := &http.Client{
		Jar:     jar,
		Timeout: time.Second * 3,
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	defer func() {
		if err != nil {
			c.String(http.StatusConflict, err.Error())
		}
	}()
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", vars.Conf.PostCommentUrl, strings.NewReader(string(body)))
	if err != nil {
		return
	}
	defer req.Body.Close()
	for k, v := range c.Request.Header {
		req.Header.Set(k, v[0])
	}
	res, err := cli.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode == http.StatusOK && res.Request.Response.StatusCode == http.StatusFound {
		for _, cookie := range res.Request.Response.Cookies() {
			c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
		}
		c.Redirect(http.StatusFound, res.Request.Response.Header.Get("Location"))
		return
	}
	s, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = errors.New(string(s))
}
