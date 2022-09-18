package actions

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
)

func Login(c *gin.Context) {
	password := c.PostForm("post_password")
	ref := c.Request.Referer()
	if ref == "" {
		ref = "/"
	}
	if password == "" || strings.Replace(password, " ", "", -1) == "" {
		c.Redirect(304, ref)
		return
	}
	s := sessions.Default(c)
	s.Set("post_password", password)
	err := s.Save()
	if err != nil {
		c.Error(err)
		return
	}
	c.Redirect(302, ref)
}
