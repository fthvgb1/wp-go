package actions

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/actions/common"
	"github/fthvgb1/wp-go/models"
	"net/http"
	"strconv"
)

func Detail(c *gin.Context) {
	var err error
	recent := common.RecentPosts()
	archive := common.Archives()
	categoryItems := common.Categories()
	var h = gin.H{
		"title":       models.Options["blogname"],
		"options":     models.Options,
		"recentPosts": recent,
		"archives":    archive,
		"categories":  categoryItems,
	}
	defer func() {
		c.HTML(http.StatusOK, "posts/detail.gohtml", h)
		if err != nil {
			c.Error(err)
		}
	}()
	id := c.Param("id")
	Id := 0
	if id != "" {
		Id, err = strconv.Atoi(id)
		if err != nil {
			return
		}
	}
	ID := uint64(Id)
	post := common.GetPostFromCache(ID)
	if post.Id == 0 {
		er := common.QueryAndSetPostCache([]models.WpPosts{{Id: ID}})
		if er != nil {
			err = er
			return
		}
		post = common.GetPostFromCache(ID)
		if post.Id == 0 {
			return
		}
	}
	pw := sessions.Default(c).Get("post_password")
	showComment := false
	if post.CommentCount > 0 || post.CommentStatus == "open" {
		showComment = true
	}
	common.PasswordProjectTitle(&post)
	if post.PostPassword != "" && pw != post.PostPassword {
		common.PasswdProjectContent(&post)
		showComment = false
	}
	prev, next, err := common.GetContextPost(post.Id, post.PostDate)

	h["title"] = fmt.Sprintf("%s-%s", post.PostTitle, models.Options["blogname"])
	h["post"] = post
	h["showComment"] = showComment
	h["prev"] = prev
	h["next"] = next
}
