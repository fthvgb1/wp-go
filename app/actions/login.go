package actions

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/phphelper"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
	pass, err := phphelper.NewPasswordHash(8, true).HashPassword(password)
	if err != nil {
		c.Error(err)
		return
	}
	cohash := fmt.Sprintf("wp-postpass_%s", str.Md5(wpconfig.GetOption("siteurl")))
	c.SetCookie(cohash, pass, 24*3600, "/", "", false, false)

	c.Redirect(http.StatusFound, ref)
}
