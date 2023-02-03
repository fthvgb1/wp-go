package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/dao"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"sync"
	"time"
)

var postContextCache *cache.MapCache[uint64, dao.PostContext]
var archivesCaches *Arch
var categoryAndTagsCaches *cache.VarCache[[]models.TermsMy]
var recentPostsCaches *cache.VarCache[[]models.Posts]
var recentCommentsCaches *cache.VarCache[[]models.Comments]
var postCommentCaches *cache.MapCache[uint64, []uint64]
var postsCache *cache.MapCache[uint64, models.Posts]

var postMetaCache *cache.MapCache[uint64, map[string]any]

var monthPostsCache *cache.MapCache[string, []uint64]
var postListIdsCache *cache.MapCache[string, dao.PostIds]
var searchPostIdsCache *cache.MapCache[string, dao.PostIds]
var maxPostIdCache *cache.VarCache[uint64]

var usersCache *cache.MapCache[uint64, models.Users]
var usersNameCache *cache.MapCache[string, models.Users]
var commentsCache *cache.MapCache[uint64, models.Comments]

var feedCache *cache.VarCache[[]string]

var postFeedCache *cache.MapCache[string, string]

var commentsFeedCache *cache.VarCache[[]string]

var newCommentCache *cache.MapCache[string, string]

var allUsernameCache *cache.VarCache[map[string]struct{}]

var headerImagesCache *cache.MapCache[string, []models.PostThumbnail]

var ctx context.Context

func InitActionsCommonCache() {
	c := config.Conf.Load()
	archivesCaches = &Arch{
		mutex:        &sync.Mutex{},
		setCacheFunc: dao.Archives,
	}

	searchPostIdsCache = cache.NewMemoryMapCacheByFn[string](dao.SearchPostIds, c.SearchPostCacheTime)

	postListIdsCache = cache.NewMemoryMapCacheByFn[string](dao.SearchPostIds, c.PostListCacheTime)

	monthPostsCache = cache.NewMemoryMapCacheByFn[string](dao.MonthPost, c.MonthPostCacheTime)

	postContextCache = cache.NewMemoryMapCacheByFn[uint64](dao.GetPostContext, c.ContextPostCacheTime)

	postsCache = cache.NewMemoryMapCacheByBatchFn(dao.GetPostsByIds, c.PostDataCacheTime)

	postMetaCache = cache.NewMemoryMapCacheByBatchFn(dao.GetPostMetaByPostIds, c.PostDataCacheTime)

	categoryAndTagsCaches = cache.NewVarCache(dao.CategoriesAndTags, c.CategoryCacheTime)

	recentPostsCaches = cache.NewVarCache(dao.RecentPosts, c.RecentPostCacheTime)

	recentCommentsCaches = cache.NewVarCache(dao.RecentComments, c.RecentCommentsCacheTime)

	postCommentCaches = cache.NewMemoryMapCacheByFn[uint64](dao.PostComments, c.PostCommentsCacheTime)

	maxPostIdCache = cache.NewVarCache(dao.GetMaxPostId, c.MaxPostIdCacheTime)

	usersCache = cache.NewMemoryMapCacheByFn[uint64](dao.GetUserById, c.UserInfoCacheTime)

	usersNameCache = cache.NewMemoryMapCacheByFn[string](dao.GetUserByName, c.UserInfoCacheTime)

	commentsCache = cache.NewMemoryMapCacheByBatchFn(dao.GetCommentByIds, c.CommentsCacheTime)

	feedCache = cache.NewVarCache(feed, time.Hour)

	postFeedCache = cache.NewMemoryMapCacheByFn[string](postFeed, time.Hour)

	commentsFeedCache = cache.NewVarCache(commentsFeed, time.Hour)

	newCommentCache = cache.NewMemoryMapCacheByFn[string, string](nil, 15*time.Minute)

	allUsernameCache = cache.NewVarCache(dao.AllUsername, c.UserInfoCacheTime)

	headerImagesCache = cache.NewMemoryMapCacheByFn[string](getHeaderImages, c.ThemeHeaderImagCacheTime)

	ctx = context.Background()

	InitFeed()
}

func ClearCache() {
	searchPostIdsCache.ClearExpired(ctx)
	searchPostIdsCache.ClearExpired(ctx)
	postsCache.ClearExpired(ctx)
	postMetaCache.ClearExpired(ctx)
	postListIdsCache.ClearExpired(ctx)
	monthPostsCache.ClearExpired(ctx)
	postContextCache.ClearExpired(ctx)
	usersCache.ClearExpired(ctx)
	commentsCache.ClearExpired(ctx)
	usersNameCache.ClearExpired(ctx)
	postFeedCache.ClearExpired(ctx)
	newCommentCache.ClearExpired(ctx)
	headerImagesCache.ClearExpired(ctx)
}
func FlushCache() {
	searchPostIdsCache.Flush(ctx)
	postsCache.Flush(ctx)
	postMetaCache.Flush(ctx)
	postListIdsCache.Flush(ctx)
	monthPostsCache.Flush(ctx)
	postContextCache.Flush(ctx)
	usersCache.Flush(ctx)
	commentsCache.Flush(ctx)
	usersCache.Flush(ctx)
	postFeedCache.Flush(ctx)
	newCommentCache.Flush(ctx)
	headerImagesCache.Flush(ctx)
}

func Archives(ctx context.Context) (r []models.PostArchive) {
	return archivesCaches.getArchiveCache(ctx)
}

type Arch struct {
	data         []models.PostArchive
	mutex        *sync.Mutex
	setCacheFunc func(context.Context) ([]models.PostArchive, error)
	month        time.Month
}

func (c *Arch) getArchiveCache(ctx context.Context) []models.PostArchive {
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

func Categories(ctx context.Context) []models.TermsMy {
	r, err := categoryAndTagsCaches.GetCache(ctx, time.Second, ctx)
	logs.ErrPrintln(err, "get category err")
	r = slice.Filter(r, func(my models.TermsMy) bool {
		return my.Taxonomy == "category"
	})
	return r
}

func Tags(ctx context.Context) []models.TermsMy {
	r, err := categoryAndTagsCaches.GetCache(ctx, time.Second, ctx)
	logs.ErrPrintln(err, "get category err")
	r = slice.Filter(r, func(my models.TermsMy) bool {
		return my.Taxonomy == "post_tag"
	})
	return r
}
func AllTagsNames(ctx context.Context) map[string]struct{} {
	r, err := categoryAndTagsCaches.GetCache(ctx, time.Second, ctx)
	logs.ErrPrintln(err, "get category err")
	return slice.FilterAndToMap(r, func(t models.TermsMy) (string, struct{}, bool) {
		if t.Taxonomy == "post_tag" {
			return t.Name, struct{}{}, true
		}
		return "", struct{}{}, false
	})
}

func AllCategoryNames(ctx context.Context) map[string]struct{} {
	r, err := categoryAndTagsCaches.GetCache(ctx, time.Second, ctx)
	logs.ErrPrintln(err, "get category err")
	return slice.FilterAndToMap(r, func(t models.TermsMy) (string, struct{}, bool) {
		if t.Taxonomy == "category" {
			return t.Name, struct{}{}, true
		}
		return "", struct{}{}, false
	})
}
