package cache

import (
	"context"
	"github/fthvgb1/wp-go/cache"
	"github/fthvgb1/wp-go/internal/pkg/config"
	"github/fthvgb1/wp-go/internal/pkg/dao"
	"github/fthvgb1/wp-go/internal/pkg/logs"
	models2 "github/fthvgb1/wp-go/internal/pkg/models"
	"sync"
	"time"
)

var postContextCache *cache.MapCache[uint64, common.PostContext]
var archivesCaches *Arch
var categoryCaches *cache.SliceCache[models2.TermsMy]
var recentPostsCaches *cache.SliceCache[models2.Posts]
var recentCommentsCaches *cache.SliceCache[models2.Comments]
var postCommentCaches *cache.MapCache[uint64, []uint64]
var postsCache *cache.MapCache[uint64, models2.Posts]

var postMetaCache *cache.MapCache[uint64, map[string]any]

var monthPostsCache *cache.MapCache[string, []uint64]
var postListIdsCache *cache.MapCache[string, common.PostIds]
var searchPostIdsCache *cache.MapCache[string, common.PostIds]
var maxPostIdCache *cache.SliceCache[uint64]

var usersCache *cache.MapCache[uint64, models2.Users]
var usersNameCache *cache.MapCache[string, models2.Users]
var commentsCache *cache.MapCache[uint64, models2.Comments]

func InitActionsCommonCache() {
	c := config.Conf.Load()
	archivesCaches = &Arch{
		mutex:        &sync.Mutex{},
		setCacheFunc: common.Archives,
	}

	searchPostIdsCache = cache.NewMapCacheByFn[string, common.PostIds](common.SearchPostIds, c.SearchPostCacheTime)

	postListIdsCache = cache.NewMapCacheByFn[string, common.PostIds](common.SearchPostIds, c.PostListCacheTime)

	monthPostsCache = cache.NewMapCacheByFn[string, []uint64](common.MonthPost, c.MonthPostCacheTime)

	postContextCache = cache.NewMapCacheByFn[uint64, common.PostContext](common.GetPostContext, c.ContextPostCacheTime)

	postsCache = cache.NewMapCacheByBatchFn[uint64, models2.Posts](common.GetPostsByIds, c.PostDataCacheTime)

	postMetaCache = cache.NewMapCacheByBatchFn[uint64, map[string]any](common.GetPostMetaByPostIds, c.PostDataCacheTime)

	categoryCaches = cache.NewSliceCache[models2.TermsMy](common.Categories, c.CategoryCacheTime)

	recentPostsCaches = cache.NewSliceCache[models2.Posts](common.RecentPosts, c.RecentPostCacheTime)

	recentCommentsCaches = cache.NewSliceCache[models2.Comments](common.RecentComments, c.RecentCommentsCacheTime)

	postCommentCaches = cache.NewMapCacheByFn[uint64, []uint64](common.PostComments, c.PostCommentsCacheTime)

	maxPostIdCache = cache.NewSliceCache[uint64](common.GetMaxPostId, c.MaxPostIdCacheTime)

	usersCache = cache.NewMapCacheByFn[uint64, models2.Users](common.GetUserById, c.UserInfoCacheTime)

	usersNameCache = cache.NewMapCacheByFn[string, models2.Users](common.GetUserByName, c.UserInfoCacheTime)

	commentsCache = cache.NewMapCacheByBatchFn[uint64, models2.Comments](common.GetCommentByIds, c.CommentsCacheTime)
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
}

func Archives(ctx context.Context) (r []models2.PostArchive) {
	return archivesCaches.getArchiveCache(ctx)
}

type Arch struct {
	data         []models2.PostArchive
	mutex        *sync.Mutex
	setCacheFunc func(context.Context) ([]models2.PostArchive, error)
	month        time.Month
}

func (c *Arch) getArchiveCache(ctx context.Context) []models2.PostArchive {
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

func Categories(ctx context.Context) []models2.TermsMy {
	r, err := categoryCaches.GetCache(ctx, time.Second, ctx)
	logs.ErrPrintln(err, "get category ")
	return r
}
