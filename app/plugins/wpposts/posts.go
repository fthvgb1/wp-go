package wpposts

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/models"
)

func PasswordProjectTitle(post *models.Posts) {
	post.PostTitle = fmt.Sprintf("密码保护：%s", post.PostTitle)
}

func PasswdProjectContent(post *models.Posts) {
	format := `
<form action="/login" class="post-password-form" method="post">
<p>此内容受密码保护。如需查阅，请在下列字段中输入您的密码。</p>
<p><label for="pwbox-%d">密码： <input name="post_password" id="pwbox-%d" type="password" size="20"></label> <input type="submit" name="Submit" value="提交"></p>
</form>`
	post.PostContent = fmt.Sprintf(format, post.Id, post.Id)
}
