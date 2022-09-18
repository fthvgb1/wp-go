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
	id := c.Param("id")
	var h = gin.H{}
	var err error
	defer func() {
		c.HTML(http.StatusOK, "detail", h)
		if err != nil {
			c.Error(err)
		}
	}()

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
	showComment := true
	common.PasswordProjectTitle(&post)
	if post.PostPassword != "" && pw != post.PostPassword {
		common.PasswdProjectContent(&post)
		showComment = false
	}
	recent, err := common.RecentPosts()
	archive, err := common.Archives()
	categoryItems, err := common.Categories()
	canComment := false
	if post.CommentStatus == "open" {
		canComment = true
	}
	prev, err := models.FirstOne[models.WpPosts](models.SqlBuilder{
		{"post_date", "<", post.PostDate.Format("2006-01-02 15:04:05")},
		{"post_status", "publish"},
		{"post_type", "post"},
	}, "ID,post_title")
	if prev.Id > 0 {
		if _, ok := common.PostsCache.Load(prev.Id); !ok {
			common.QueryAndSetPostCache([]models.WpPosts{prev})
		}
	}
	next, err := models.FirstOne[models.WpPosts](models.SqlBuilder{
		{"post_date", ">", post.PostDate.Format("2006-01-02 15:04:05")},
		{"post_status", "publish"},
		{"post_type", "post"},
	}, "ID,post_title,post_password")
	if prev.Id > 0 {
		if _, ok := common.PostsCache.Load(next.Id); !ok {
			common.QueryAndSetPostCache([]models.WpPosts{next})
		}
	}
	h = gin.H{
		"title":       fmt.Sprintf("%s-%s", post.PostTitle, models.Options["blogname"]),
		"post":        post,
		"options":     models.Options,
		"recentPosts": recent,
		"archives":    archive,
		"categories":  categoryItems,
		"comment":     showComment,
		"canComment":  canComment,
		"prev":        prev,
		"next":        next,
	}
}
