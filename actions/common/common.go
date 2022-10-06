package common

import (
	"context"
	"database/sql"
	"fmt"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/models"
	"github/fthvgb1/wp-go/vars"
	"strconv"
	"sync"
	"time"
)

var postContextCache *cache.MapCache[uint64, PostContext]
var archivesCaches *Arch
var categoryCaches *cache.SliceCache[models.WpTermsMy]
var recentPostsCaches *cache.SliceCache[models.WpPosts]
var recentCommentsCaches *cache.SliceCache[models.WpComments]
var postCommentCaches *cache.MapCache[uint64, []models.WpComments]
var postsCache *cache.MapCache[uint64, models.WpPosts]
var monthPostsCache *cache.MapCache[string, []uint64]
var postListIdsCache *cache.MapCache[string, PostIds]
var searchPostIdsCache *cache.MapCache[string, PostIds]
var maxPostIdCache *cache.SliceCache[uint64]
var TotalRaw int
var usersCache *cache.MapCache[uint64, models.WpUsers]

func InitActionsCommonCache() {
	archivesCaches = &Arch{
		mutex:        &sync.Mutex{},
		setCacheFunc: archives,
	}

	searchPostIdsCache = cache.NewMapCacheByFn[string, PostIds](searchPostIds, vars.Conf.SearchPostCacheTime)

	postListIdsCache = cache.NewMapCacheByFn[string, PostIds](searchPostIds, vars.Conf.PostListCacheTime)

	monthPostsCache = cache.NewMapCacheByFn[string, []uint64](monthPost, vars.Conf.MonthPostCacheTime)

	postContextCache = cache.NewMapCacheByFn[uint64, PostContext](getPostContext, vars.Conf.ContextPostCacheTime)

	postsCache = cache.NewMapCacheByBatchFn[uint64, models.WpPosts](getPostsByIds, vars.Conf.PostDataCacheTime)
	postsCache.SetCacheFunc(getPostById)

	categoryCaches = cache.NewSliceCache[models.WpTermsMy](categories, vars.Conf.CategoryCacheTime)

	recentPostsCaches = cache.NewSliceCache[models.WpPosts](recentPosts, vars.Conf.RecentPostCacheTime)

	recentCommentsCaches = cache.NewSliceCache[models.WpComments](recentComments, vars.Conf.RecentCommentsCacheTime)

	postCommentCaches = cache.NewMapCacheByFn[uint64, []models.WpComments](postComments, vars.Conf.CommentsCacheTime)

	maxPostIdCache = cache.NewSliceCache[uint64](getMaxPostId, vars.Conf.MaxPostIdCacheTime)

	usersCache = cache.NewMapCacheByBatchFn[uint64, models.WpUsers](getUsers, time.Hour)
	usersCache.SetCacheFunc(getUser)
}

func ClearCache() {
	searchPostIdsCache.ClearExpired()
	postsCache.ClearExpired()
	postsCache.ClearExpired()
	postListIdsCache.ClearExpired()
	monthPostsCache.ClearExpired()
	postContextCache.ClearExpired()
}

func GetMonthPostIds(ctx context.Context, year, month string, page, limit int, order string) (r []models.WpPosts, total int, err error) {
	res, err := monthPostsCache.GetCache(ctx, fmt.Sprintf("%s%s", year, month), time.Second, year, month)
	if err != nil {
		return
	}
	if order == "desc" {
		res = helper.SliceReverse(res)
	}
	total = len(res)
	rr := helper.SlicePagination(res, page, limit)
	r, err = GetPostsByIds(ctx, rr)
	return
}

func monthPost(args ...any) (r []uint64, err error) {
	year, month := args[0].(string), args[1].(string)
	where := models.SqlBuilder{
		{"post_type", "in", ""},
		{"post_status", "in", ""},
		{"year(post_date)", year},
		{"month(post_date)", month},
	}
	postType := []any{"post"}
	status := []any{"publish"}
	ids, err := models.Find[models.WpPosts](where, "ID", "", models.SqlBuilder{{"Id", "asc"}}, nil, 0, postType, status)
	if err != nil {
		return
	}
	for _, post := range ids {
		r = append(r, post.Id)
	}
	return
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

func PostComments(ctx context.Context, Id uint64) ([]models.WpComments, error) {
	return postCommentCaches.GetCache(ctx, Id, time.Second, Id)
}

func postComments(args ...any) ([]models.WpComments, error) {
	postId := args[0].(uint64)
	return models.Find[models.WpComments](models.SqlBuilder{
		{"comment_approved", "1"},
		{"comment_post_ID", "=", strconv.FormatUint(postId, 10), "int"},
	}, "*", "", models.SqlBuilder{
		{"comment_date_gmt", "asc"},
		{"comment_ID", "asc"},
	}, nil, 0)
}

func RecentComments(ctx context.Context, n int) (r []models.WpComments) {
	r, err := recentCommentsCaches.GetCache(ctx, time.Second)
	if len(r) > n {
		r = r[0:n]
	}
	logs.ErrPrintln(err, "get recent comment")
	return
}
func recentComments(...any) (r []models.WpComments, err error) {
	return models.Find[models.WpComments](models.SqlBuilder{
		{"comment_approved", "1"},
		{"post_status", "publish"},
	}, "comment_ID,comment_author,comment_post_ID,post_title", "", models.SqlBuilder{{"comment_date_gmt", "desc"}}, models.SqlBuilder{
		{"a", "left join", "wp_posts b", "a.comment_post_ID=b.ID"},
	}, 10)
}

func GetContextPost(ctx context.Context, id uint64, date time.Time) (prev, next models.WpPosts, err error) {
	postCtx, err := postContextCache.GetCache(ctx, id, time.Second, date)
	if err != nil {
		return models.WpPosts{}, models.WpPosts{}, err
	}
	prev = postCtx.prev
	next = postCtx.next
	return
}

func getPostContext(arg ...any) (r PostContext, err error) {
	t := arg[0].(time.Time)
	next, err := models.FirstOne[models.WpPosts](models.SqlBuilder{
		{"post_date", ">", t.Format("2006-01-02 15:04:05")},
		{"post_status", "in", ""},
		{"post_type", "post"},
	}, "ID,post_title,post_password", nil, []any{"publish"})
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return
	}
	prev, err := models.FirstOne[models.WpPosts](models.SqlBuilder{
		{"post_date", "<", t.Format("2006-01-02 15:04:05")},
		{"post_status", "in", ""},
		{"post_type", "post"},
	}, "ID,post_title", models.SqlBuilder{{"post_date", "desc"}}, []any{"publish"})
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return
	}
	r = PostContext{
		prev: prev,
		next: next,
	}
	return
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

func RecentPosts(ctx context.Context, n int) (r []models.WpPosts) {
	r, err := recentPostsCaches.GetCache(ctx, time.Second)
	if n < len(r) {
		r = r[:n]
	}
	logs.ErrPrintln(err, "get recent post")
	return
}
func recentPosts(...any) (r []models.WpPosts, err error) {
	r, err = models.Find[models.WpPosts](models.SqlBuilder{{
		"post_type", "post",
	}, {"post_status", "publish"}}, "ID,post_title,post_password", "", models.SqlBuilder{{"post_date", "desc"}}, nil, 10)
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
