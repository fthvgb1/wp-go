package common

import (
	"fmt"
	"github/fthvgb1/wp-go/models"
)

func Archives() (r []models.PostArchive, err error) {
	r, err = models.Find[models.PostArchive](models.SqlBuilder{
		{"post_type", "post"}, {"post_status", "publish"},
	}, "YEAR(post_date) AS `year`, MONTH(post_date) AS `month`, count(ID) as posts", "year,month", models.SqlBuilder{{"year", "desc"}, {"month", "desc"}}, nil, 0)
	return
}

func Categories() (terms []models.WpTermsMy, err error) {
	var in = []interface{}{"category"}
	terms, err = models.Find[models.WpTermsMy](models.SqlBuilder{
		{"tt.count", ">", "0", "int"},
		{"tt.taxonomy", "in", ""},
	}, "t.term_id", "", models.SqlBuilder{
		{"t.name", "asc"},
	}, models.SqlBuilder{
		{"t", "inner join", "wp_term_taxonomy tt", "t.term_id = tt.term_id"},
	}, 0, in)
	for i := 0; i < len(terms); i++ {
		if v, ok := models.Terms[terms[i].WpTerms.TermId]; ok {
			terms[i].WpTerms = v
		}
		if v, ok := models.TermTaxonomy[terms[i].WpTerms.TermId]; ok {
			terms[i].WpTermTaxonomy = v
		}
	}
	return
}

func RecentPosts() (r []models.WpPosts, err error) {
	r, err = models.Find[models.WpPosts](models.SqlBuilder{{
		"post_type", "post",
	}, {"post_status", "publish"}}, "ID,post_title,post_password", "", models.SqlBuilder{{"post_date", "desc"}}, nil, 5)
	return
}

func PasswdProject(post *models.WpPosts) {
	if post.PostTitle != "" {
		post.PostTitle = fmt.Sprintf("密码保护：%s", post.PostTitle)
	}
	if post.PostContent != "" {
		format := `
<form action="/login" class="post-password-form" method="post">
<p>此内容受密码保护。如需查阅，请在下列字段中输入您的密码。</p>
<p><label for="pwbox-%d">密码： <input name="post_password" id="pwbox-%d" type="password" size="20"></label> <input type="submit" name="Submit" value="提交"></p>
</form>`
		post.PostContent = fmt.Sprintf(format, post.Id, post.Id)
	}
}
