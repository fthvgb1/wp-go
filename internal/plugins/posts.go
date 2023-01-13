package plugins

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/internal/pkg/models"
)

func NewPostPlugin(ctx *gin.Context, scene uint) *Plugin[models.Posts] {
	p := NewPlugin[models.Posts](nil, -1, nil, scene, ctx)
	p.Push(Digest)
	return p
}

func ApplyPlugin(p *Plugin[models.Posts], post *models.Posts) {
	p.post = post
	p.Next()
	p.index = -1
}

func PasswordProjectTitle(post *models.Posts) {
	if post.PostPassword != "" {
		post.PostTitle = fmt.Sprintf("密码保护：%s", post.PostTitle)
	}
}

func PasswdProjectContent(post *models.Posts) {
	if post.PostContent != "" {
		format := `
<form action="/login" class="post-password-form" method="post">
<p>此内容受密码保护。如需查阅，请在下列字段中输入您的密码。</p>
<p><label for="pwbox-%d">密码： <input name="post_password" id="pwbox-%d" type="password" size="20"></label> <input type="submit" name="Submit" value="提交"></p>
</form>`
		post.PostContent = fmt.Sprintf(format, post.Id, post.Id)
	}
}
