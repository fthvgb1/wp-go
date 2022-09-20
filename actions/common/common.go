package common

import (
	"context"
	"fmt"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/vars"
	"log"
	"strings"
	"sync"
	"time"
)

var PostsCache sync.Map
var PostContextCache sync.Map

var archivesCaches *Arch
var categoryCaches *cache.SliceCache[models.WpTermsMy]
var recentPostsCaches *cache.SliceCache[models.WpPosts]

func InitCache() {
	archivesCaches = &Arch{
		mutex:        &sync.Mutex{},
		setCacheFunc: archives,
	}
	categoryCaches = cache.NewSliceCache[models.WpTermsMy](categories, vars.Conf.CategoryCacheTime)
	recentPostsCaches = cache.NewSliceCache[models.WpPosts](recentPosts, vars.Conf.RecentPostCacheTime)
}

type Arch struct {
	data         []models.PostArchive
	mutex        *sync.Mutex
	setCacheFunc func() ([]models.PostArchive, error)
	month        time.Month
}

func (c *Arch) GetCache() []models.PostArchive {
	l := len(c.data)
	m := time.Now().Month()
	if l > 0 && c.month != m || l < 1 {
		r, err := c.setCacheFunc()
		if err != nil {
			log.Printf("set cache err[%s]", err)
			return nil
		}
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.month = m
		c.data = r
	}
	return c.data
}

type PostContext struct {
	Prev       models.WpPosts
	Next       models.WpPosts
	expireTime time.Duration
	setTime    time.Time
}

func GetContextPost(id uint64, t time.Time) (prev, next models.WpPosts, err error) {
	post, ok := PostContextCache.Load(id)
	if ok {
		c := post.(PostContext)
		isExp := c.expireTime/time.Second+time.Duration(c.setTime.Unix()) < time.Duration(time.Now().Unix())
		if !isExp && (c.Prev.Id > 0 || c.Next.Id > 0) {
			return c.Prev, c.Next, nil
		}
	}
	prev, next, err = getPostContext(t)
	post = PostContext{
		Prev:       prev,
		Next:       next,
		expireTime: vars.Conf.ContextPostCacheTime,
		setTime:    time.Now(),
	}
	PostContextCache.Store(id, post)
	return
}

func getPostContext(t time.Time) (prev, next models.WpPosts, err error) {
	next, err = models.FirstOne[models.WpPosts](models.SqlBuilder{
		{"post_date", ">", t.Format("2006-01-02 15:04:05")},
		{"post_status", "in", ""},
		{"post_type", "post"},
	}, "ID,post_title,post_password", nil, []interface{}{"publish", "private"})
	if _, ok := PostsCache.Load(next.Id); !ok {

	}
	prev, err = models.FirstOne[models.WpPosts](models.SqlBuilder{
		{"post_date", "<", t.Format("2006-01-02 15:04:05")},
		{"post_status", "in", ""},
		{"post_type", "post"},
	}, "ID,post_title", models.SqlBuilder{{"post_date", "desc"}}, []interface{}{"publish", "private"})
	return
}

func GetPostFromCache(Id uint64) (r models.WpPosts) {
	p, ok := PostsCache.Load(Id)
	if ok {
		r = *p.(*models.WpPosts)
	}
	return
}

func QueryAndSetPostCache(postIds []models.WpPosts) (err error) {
	var all []uint64
	var needQuery []interface{}
	for _, wpPosts := range postIds {
		all = append(all, wpPosts.Id)
		if _, ok := PostsCache.Load(wpPosts.Id); !ok {
			needQuery = append(needQuery, wpPosts.Id)
		}
	}
	if len(needQuery) > 0 {
		rawPosts, er := models.Find[models.WpPosts](models.SqlBuilder{{
			"Id", "in", "",
		}}, "a.*,ifnull(d.name,'') category_name,ifnull(taxonomy,'') `taxonomy`", "", nil, models.SqlBuilder{{
			"a", "left join", "wp_term_relationships b", "a.Id=b.object_id",
		}, {
			"left join", "wp_term_taxonomy c", "b.term_taxonomy_id=c.term_taxonomy_id",
		}, {
			"left join", "wp_terms d", "c.term_id=d.term_id",
		}}, 0, needQuery)
		if er != nil {
			err = er
			return
		}
		postsMap := make(map[uint64]*models.WpPosts)
		for i, post := range rawPosts {
			v, ok := postsMap[post.Id]
			if !ok {
				v = &rawPosts[i]
			}
			if post.Taxonomy == "category" {
				v.Categories = append(v.Categories, post.CategoryName)
			} else if post.Taxonomy == "post_tag" {
				v.Tags = append(v.Tags, post.CategoryName)
			}
			postsMap[post.Id] = v
		}
		for _, pp := range postsMap {
			if len(pp.Categories) > 0 {
				t := make([]string, 0, len(pp.Categories))
				for _, cat := range pp.Categories {
					t = append(t, fmt.Sprintf(`<a href="/p/category/%s" rel="category tag">%s</a>`, cat, cat))
				}
				pp.CategoriesHtml = strings.Join(t, "、")
			}
			if len(pp.Tags) > 0 {
				t := make([]string, 0, len(pp.Tags))
				for _, cat := range pp.Tags {
					t = append(t, fmt.Sprintf(`<a href="/p/tag/%s" rel="tag">%s</a>`, cat, cat))
				}
				pp.TagsHtml = strings.Join(t, "、")
			}
			PostsCache.Store(pp.Id, pp)
		}
	}
	return
}

func archives() ([]models.PostArchive, error) {
	return models.Find[models.PostArchive](models.SqlBuilder{
		{"post_type", "post"}, {"post_status", "publish"},
	}, "YEAR(post_date) AS `year`, MONTH(post_date) AS `month`, count(ID) as posts", "year,month", models.SqlBuilder{{"year", "desc"}, {"month", "desc"}}, nil, 0)
}

func Archives() (r []models.PostArchive) {
	return archivesCaches.GetCache()
}

func Categories(ctx context.Context) []models.WpTermsMy {
	r, _ := categoryCaches.GetCache(ctx, time.Second)
	return r
}

func categories(...any) (terms []models.WpTermsMy, err error) {
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

func RecentPosts(ctx context.Context) (r []models.WpPosts) {
	r, _ = recentPostsCaches.GetCache(ctx, time.Second)
	return
}
func recentPosts(...any) (r []models.WpPosts, err error) {
	r, err = models.Find[models.WpPosts](models.SqlBuilder{{
		"post_type", "post",
	}, {"post_status", "publish"}}, "ID,post_title,post_password", "", models.SqlBuilder{{"post_date", "desc"}}, nil, 5)
	for i, post := range r {
		if post.PostPassword != "" {
			PasswordProjectTitle(&r[i])
		}
	}
	return
}

func PasswordProjectTitle(post *models.WpPosts) {
	if post.PostPassword != "" {
		post.PostTitle = fmt.Sprintf("密码保护：%s", post.PostTitle)
	}
}

func PasswdProjectContent(post *models.WpPosts) {
	if post.PostContent != "" {
		format := `
<form action="/login" class="post-password-form" method="post">
<p>此内容受密码保护。如需查阅，请在下列字段中输入您的密码。</p>
<p><label for="pwbox-%d">密码： <input name="post_password" id="pwbox-%d" type="password" size="20"></label> <input type="submit" name="Submit" value="提交"></p>
</form>`
		post.PostContent = fmt.Sprintf(format, post.Id, post.Id)
	}
}
