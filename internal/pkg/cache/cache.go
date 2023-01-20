package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/dao"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"sync"
	"time"
)

var postContextCache *cache.MapCache[uint64, common.PostContext]
var archivesCaches *Arch
var categoryCaches *cache.VarCache[[]models.TermsMy]
var recentPostsCaches *cache.VarCache[[]models.Posts]
var recentCommentsCaches *cache.VarCache[[]models.Comments]
var postCommentCaches *cache.MapCache[uint64, []uint64]
var postsCache *cache.MapCache[uint64, models.Posts]

var postMetaCache *cache.MapCache[uint64, map[string]any]

var monthPostsCache *cache.MapCache[string, []uint64]
var postListIdsCache *cache.MapCache[string, common.PostIds]
var searchPostIdsCache *cache.MapCache[string, common.PostIds]
var maxPostIdCache *cache.VarCache[uint64]

var usersCache *cache.MapCache[uint64, models.Users]
var usersNameCache *cache.MapCache[string, models.Users]
var commentsCache *cache.MapCache[uint64, models.Comments]

var feedCache *cache.VarCache[[]string]

var postFeedCache *cache.MapCache[string, string]

var commentsFeedCache *cache.VarCache[[]string]

var newCommentCache *cache.MapCache[string, string]

func InitActionsCommonCache() {
	c := config.Conf.Load()
	archivesCaches = &Arch{
		mutex:        &sync.Mutex{},
		setCacheFunc: common.Archives,
	}

	searchPostIdsCache = cache.NewMapCacheByFn[string](common.SearchPostIds, c.SearchPostCacheTime)

	postListIdsCache = cache.NewMapCacheByFn[string](common.SearchPostIds, c.PostListCacheTime)

	monthPostsCache = cache.NewMapCacheByFn[string](common.MonthPost, c.MonthPostCacheTime)

	postContextCache = cache.NewMapCacheByFn[uint64](common.GetPostContext, c.ContextPostCacheTime)

	postsCache = cache.NewMapCacheByBatchFn(common.GetPostsByIds, c.PostDataCacheTime)

	postMetaCache = cache.NewMapCacheByBatchFn(common.GetPostMetaByPostIds, c.PostDataCacheTime)

	categoryCaches = cache.NewVarCache(common.Categories, c.CategoryCacheTime)

	recentPostsCaches = cache.NewVarCache(common.RecentPosts, c.RecentPostCacheTime)

	recentCommentsCaches = cache.NewVarCache(common.RecentComments, c.RecentCommentsCacheTime)

	postCommentCaches = cache.NewMapCacheByFn[uint64](common.PostComments, c.PostCommentsCacheTime)

	maxPostIdCache = cache.NewVarCache(common.GetMaxPostId, c.MaxPostIdCacheTime)

	usersCache = cache.NewMapCacheByFn[uint64](common.GetUserById, c.UserInfoCacheTime)

	usersNameCache = cache.NewMapCacheByFn[string](common.GetUserByName, c.UserInfoCacheTime)

	commentsCache = cache.NewMapCacheByBatchFn(common.GetCommentByIds, c.CommentsCacheTime)

	feedCache = cache.NewVarCache(feed, time.Hour)

	postFeedCache = cache.NewMapCacheByFn[string](postFeed, time.Hour)

	commentsFeedCache = cache.NewVarCache(commentsFeed, time.Hour)

	newCommentCache = cache.NewMapCacheByFn[string, string](nil, 15*time.Minute)

	InitFeed()
}

func ClearCache() {
	searchPostIdsCache.ClearExpired()
	postsCache.ClearExpired()
	postMetaCache.ClearExpired()
	postListIdsCache.ClearExpired()
	monthPostsCache.ClearExpired()
	postContextCache.ClearExpired()
	usersCache.ClearExpired()
	commentsCache.ClearExpired()
	usersNameCache.ClearExpired()
	postFeedCache.ClearExpired()
	newCommentCache.ClearExpired()
}
func FlushCache() {
	searchPostIdsCache.Flush()
	postsCache.Flush()
	postMetaCache.Flush()
	postListIdsCache.Flush()
	monthPostsCache.Flush()
	postContextCache.Flush()
	usersCache.Flush()
	commentsCache.Flush()
	usersCache.Flush()
	postFeedCache.Flush()
	newCommentCache.Flush()
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
	r, err := categoryCaches.GetCache(ctx, time.Second, ctx)
	logs.ErrPrintln(err, "get category ")
	return r
}
