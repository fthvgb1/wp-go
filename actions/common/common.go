package common

import (
	"context"
	"fmt"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/config"
	"github/fthvgb1/wp-go/config/wpconfig"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/models/wp"
	"sync"
	"time"
)

var postContextCache *cache.MapCache[uint64, PostContext]
var archivesCaches *Arch
var categoryCaches *cache.SliceCache[wp.TermsMy]
var recentPostsCaches *cache.SliceCache[wp.Posts]
var recentCommentsCaches *cache.SliceCache[wp.Comments]
var postCommentCaches *cache.MapCache[uint64, []uint64]
var postsCache *cache.MapCache[uint64, wp.Posts]
var monthPostsCache *cache.MapCache[string, []uint64]
var postListIdsCache *cache.MapCache[string, PostIds]
var searchPostIdsCache *cache.MapCache[string, PostIds]
var maxPostIdCache *cache.SliceCache[uint64]
var TotalRaw int
var usersCache *cache.MapCache[uint64, wp.Users]
var usersNameCache *cache.MapCache[string, wp.Users]
var commentsCache *cache.MapCache[uint64, wp.Comments]

func InitActionsCommonCache() {
	c := config.Conf.Load()
	archivesCaches = &Arch{
		mutex:        &sync.Mutex{},
		setCacheFunc: archives,
	}

	searchPostIdsCache = cache.NewMapCacheByFn[string, PostIds](searchPostIds, c.SearchPostCacheTime)

	postListIdsCache = cache.NewMapCacheByFn[string, PostIds](searchPostIds, c.PostListCacheTime)

	monthPostsCache = cache.NewMapCacheByFn[string, []uint64](monthPost, c.MonthPostCacheTime)

	postContextCache = cache.NewMapCacheByFn[uint64, PostContext](getPostContext, c.ContextPostCacheTime)

	postsCache = cache.NewMapCacheByBatchFn[uint64, wp.Posts](getPostsByIds, c.PostDataCacheTime)

	categoryCaches = cache.NewSliceCache[wp.TermsMy](categories, c.CategoryCacheTime)

	recentPostsCaches = cache.NewSliceCache[wp.Posts](recentPosts, c.RecentPostCacheTime)

	recentCommentsCaches = cache.NewSliceCache[wp.Comments](recentComments, c.RecentCommentsCacheTime)

	postCommentCaches = cache.NewMapCacheByFn[uint64, []uint64](postComments, c.PostCommentsCacheTime)

	maxPostIdCache = cache.NewSliceCache[uint64](getMaxPostId, c.MaxPostIdCacheTime)

	usersCache = cache.NewMapCacheByFn[uint64, wp.Users](getUserById, c.UserInfoCacheTime)

	usersNameCache = cache.NewMapCacheByFn[string, wp.Users](getUserByName, c.UserInfoCacheTime)

	commentsCache = cache.NewMapCacheByBatchFn[uint64, wp.Comments](getCommentByIds, c.CommentsCacheTime)
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
	usersNameCache.ClearExpired()
}
func FlushCache() {
	searchPostIdsCache.Flush()
	postsCache.Flush()
	postsCache.Flush()
	postListIdsCache.Flush()
	monthPostsCache.Flush()
	postContextCache.Flush()
	usersCache.Flush()
	commentsCache.Flush()
	usersCache.Flush()
}

type PostIds struct {
	Ids    []uint64
	Length int
}

type Arch struct {
	data         []wp.PostArchive
	mutex        *sync.Mutex
	setCacheFunc func(context.Context) ([]wp.PostArchive, error)
	month        time.Month
}

func (c *Arch) getArchiveCache(ctx context.Context) []wp.PostArchive {
	l := len(c.data)
	m := time.Now().Month()
	if l > 0 && c.month != m || l < 1 {
		r, err := c.setCacheFunc(ctx)
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
	prev wp.Posts
	next wp.Posts
}

func archives(ctx context.Context) ([]wp.PostArchive, error) {
	return models.Find[wp.PostArchive](ctx, models.SqlBuilder{
		{"post_type", "post"}, {"post_status", "publish"},
	}, "YEAR(post_date) AS `year`, MONTH(post_date) AS `month`, count(ID) as posts", "year,month", models.SqlBuilder{{"year", "desc"}, {"month", "desc"}}, nil, nil, 0)
}

func Archives(ctx context.Context) (r []wp.PostArchive) {
	return archivesCaches.getArchiveCache(ctx)
}

func Categories(ctx context.Context) []wp.TermsMy {
	r, err := categoryCaches.GetCache(ctx, time.Second, ctx)
	logs.ErrPrintln(err, "get category ")
	return r
}

func categories(a ...any) (terms []wp.TermsMy, err error) {
	ctx := a[0].(context.Context)
	var in = []any{"category"}
	terms, err = models.Find[wp.TermsMy](ctx, models.SqlBuilder{
		{"tt.count", ">", "0", "int"},
		{"tt.taxonomy", "in", ""},
	}, "t.term_id", "", models.SqlBuilder{
		{"t.name", "asc"},
	}, models.SqlBuilder{
		{"t", "inner join", "wp_term_taxonomy tt", "t.term_id = tt.term_id"},
	}, nil, 0, in)
	for i := 0; i < len(terms); i++ {
		if v, ok := wpconfig.Terms.Load(terms[i].Terms.TermId); ok {
			terms[i].Terms = v
		}
		if v, ok := wpconfig.TermTaxonomies.Load(terms[i].Terms.TermId); ok {
			terms[i].TermTaxonomy = v
		}
	}
	return
}

func PasswordProjectTitle(post *wp.Posts) {
	if post.PostPassword != "" {
		post.PostTitle = fmt.Sprintf("密码保护：%s", post.PostTitle)
	}
}

func PasswdProjectContent(post *wp.Posts) {
	if post.PostContent != "" {
		format := `
<form action="/login" class="post-password-form" method="post">
<p>此内容受密码保护。如需查阅，请在下列字段中输入您的密码。</p>
<p><label for="pwbox-%d">密码： <input name="post_password" id="pwbox-%d" type="password" size="20"></label> <input type="submit" name="Submit" value="提交"></p>
</form>`
		post.PostContent = fmt.Sprintf(format, post.Id, post.Id)
	}
}
