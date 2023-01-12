package actions

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/internal/wpconfig"
	"github/fthvgb1/wp-go/phpass"
	"net/http"
	"strings"
)

func Login(c *gin.Context) {
	password := c.PostForm("post_password")
	ref := c.Request.Referer()
	if ref == "" {
		ref = "/"
	}
	if password == "" || strings.Replace(password, " ", "", -1) == "" {
		c.Redirect(http.StatusFound, ref)
		return
	}
	s := sessions.Default(c)
	s.Set("post_password", password)
	err := s.Save()
	if err != nil {
		c.Error(err)
		return
	}
	pass, err := phpass.NewPasswordHash(8, true).HashPassword(password)
	if err != nil {
		c.Error(err)
		return
	}
	cohash := fmt.Sprintf("wp-postpass_%s", helper.StringMd5(wpconfig.Options.Value("siteurl")))
	c.SetCookie(cohash, pass, 24*3600, "/", "", false, false)

	c.Redirect(http.StatusFound, ref)
}
