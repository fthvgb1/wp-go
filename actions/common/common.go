package common

import (
	"context"
	"fmt"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/config"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"sync"
	"time"
)

var postContextCache *cache.MapCache[uint64, PostContext]
var archivesCaches *Arch
var categoryCaches *cache.SliceCache[models.WpTermsMy]
var recentPostsCaches *cache.SliceCache[models.WpPosts]
var recentCommentsCaches *cache.SliceCache[models.WpComments]
var postCommentCaches *cache.MapCache[uint64, []uint64]
var postsCache *cache.MapCache[uint64, models.WpPosts]
var monthPostsCache *cache.MapCache[string, []uint64]
var postListIdsCache *cache.MapCache[string, PostIds]
var searchPostIdsCache *cache.MapCache[string, PostIds]
var maxPostIdCache *cache.SliceCache[uint64]
var TotalRaw int
var usersCache *cache.MapCache[uint64, models.WpUsers]
var commentsCache *cache.MapCache[uint64, models.WpComments]

func InitActionsCommonCache() {
	archivesCaches = &Arch{
		mutex:        &sync.Mutex{},
		setCacheFunc: archives,
	}

	searchPostIdsCache = cache.NewMapCacheByFn[string, PostIds](searchPostIds, config.Conf.SearchPostCacheTime)

	postListIdsCache = cache.NewMapCacheByFn[string, PostIds](searchPostIds, config.Conf.PostListCacheTime)

	monthPostsCache = cache.NewMapCacheByFn[string, []uint64](monthPost, config.Conf.MonthPostCacheTime)

	postContextCache = cache.NewMapCacheByFn[uint64, PostContext](getPostContext, config.Conf.ContextPostCacheTime)

	postsCache = cache.NewMapCacheByBatchFn[uint64, models.WpPosts](getPostsByIds, config.Conf.PostDataCacheTime)

	categoryCaches = cache.NewSliceCache[models.WpTermsMy](categories, config.Conf.CategoryCacheTime)

	recentPostsCaches = cache.NewSliceCache[models.WpPosts](recentPosts, config.Conf.RecentPostCacheTime)

	recentCommentsCaches = cache.NewSliceCache[models.WpComments](recentComments, config.Conf.RecentCommentsCacheTime)

	postCommentCaches = cache.NewMapCacheByFn[uint64, []uint64](postComments, config.Conf.PostCommentsCacheTime)

	maxPostIdCache = cache.NewSliceCache[uint64](getMaxPostId, config.Conf.MaxPostIdCacheTime)

	usersCache = cache.NewMapCacheByBatchFn[uint64, models.WpUsers](getUsers, config.Conf.UserInfoCacheTime)

	commentsCache = cache.NewMapCacheByBatchFn[uint64, models.WpComments](getCommentByIds, config.Conf.CommentsCacheTime)
}

func ClearCache() {
	searchPostIdsCache.ClearExpired()
	postsCache.ClearExpired()
	postsCache.ClearExpired()
	postListIdsCache.ClearExpired()
	monthPostsCache.ClearExpired()
	postContextCache.ClearExpired()
	usersCache.ClearExpired()
	commentsCache.ClearExpired()
}

type PostIds struct {
	Ids    []uint64
	Length int
}

type Arch struct {
	data         []models.PostArchive
	mutex        *sync.Mutex
	setCacheFunc func() ([]models.PostArchive, error)
	month        time.Month
}

func (c *Arch) getArchiveCache() []models.PostArchive {
	l := len(c.data)
	m := time.Now().Month()
	if l > 0 && c.month != m || l < 1 {
		r, err := c.setCacheFunc()
		if err != nil {
			logs.ErrPrintln(err, "set cache err[%s]")
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
	prev models.WpPosts
	next models.WpPosts
}

func archives() ([]models.PostArchive, error) {
	return models.Find[models.PostArchive](models.SqlBuilder{
		{"post_type", "post"}, {"post_status", "publish"},
	}, "YEAR(post_date) AS `year`, MONTH(post_date) AS `month`, count(ID) as posts", "year,month", models.SqlBuilder{{"year", "desc"}, {"month", "desc"}}, nil, 0)
}

func Archives() (r []models.PostArchive) {
	return archivesCaches.getArchiveCache()
}

func Categories(ctx context.Context) []models.WpTermsMy {
	r, err := categoryCaches.GetCache(ctx, time.Second)
	logs.ErrPrintln(err, "get category ")
	return r
}

func categories(...any) (terms []models.WpTermsMy, err error) {
	var in = []any{"category"}
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
