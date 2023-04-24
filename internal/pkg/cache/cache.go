package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/cmd/cachemanager"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/dao"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/safety"
	"time"
)

var postContextCache *cache.MapCache[uint64, dao.PostContext]
var categoryAndTagsCaches *cache.MapCache[string, []models.TermsMy]
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

func InitActionsCommonCache() {
	c := config.GetConfig()

	searchPostIdsCache = cachemanager.MapCacheBy[string](dao.SearchPostIds, c.CacheTime.SearchPostCacheTime)

	postListIdsCache = cachemanager.MapCacheBy[string](dao.SearchPostIds, c.CacheTime.PostListCacheTime)

	monthPostsCache = cachemanager.MapCacheBy[string](dao.MonthPost, c.CacheTime.MonthPostCacheTime)

	postContextCache = cachemanager.MapCacheBy[uint64](dao.GetPostContext, c.CacheTime.ContextPostCacheTime)

	postsCache = cachemanager.MapBatchCacheBy(dao.GetPostsByIds, c.CacheTime.PostDataCacheTime)

	postMetaCache = cachemanager.MapBatchCacheBy(dao.GetPostMetaByPostIds, c.CacheTime.PostDataCacheTime)

	categoryAndTagsCaches = cachemanager.MapCacheBy[string](dao.CategoriesAndTags, c.CacheTime.CategoryCacheTime)

	recentPostsCaches = cache.NewVarCache(dao.RecentPosts, c.CacheTime.RecentPostCacheTime)

	recentCommentsCaches = cache.NewVarCache(dao.RecentComments, c.CacheTime.RecentCommentsCacheTime)

	postCommentCaches = cachemanager.MapCacheBy[uint64](dao.PostComments, c.CacheTime.PostCommentsCacheTime)

	maxPostIdCache = cache.NewVarCache(dao.GetMaxPostId, c.CacheTime.MaxPostIdCacheTime)

	usersCache = cachemanager.MapCacheBy[uint64](dao.GetUserById, c.CacheTime.UserInfoCacheTime)

	usersNameCache = cachemanager.MapCacheBy[string](dao.GetUserByName, c.CacheTime.UserInfoCacheTime)

	commentsCache = cachemanager.MapBatchCacheBy(dao.GetCommentByIds, c.CacheTime.CommentsCacheTime)

	allUsernameCache = cache.NewVarCache(dao.AllUsername, c.CacheTime.UserInfoCacheTime)

	feedCache = cache.NewVarCache(feed, time.Hour)

	postFeedCache = cachemanager.MapCacheBy[string](postFeed, time.Hour)

	commentsFeedCache = cache.NewVarCache(commentsFeed, time.Hour)

	newCommentCache = cachemanager.MapCacheBy[string, string](nil, 15*time.Minute)

	InitFeed()
}

type Arch struct {
	data  []models.PostArchive
	fn    func(context.Context) ([]models.PostArchive, error)
	month time.Month
}

var arch = safety.NewVar(Arch{
	fn: dao.Archives,
})

func Archives(ctx context.Context) []models.PostArchive {
	a := arch.Load()
	data := a.data
	l := len(data)
	m := time.Now().Month()
	if l > 0 && a.month != m || l < 1 {
		r, err := a.fn(ctx)
		if err != nil {
			logs.Error(err, "set cache fail")
			return nil
		}
		a.month = m
		a.data = r
		arch.Store(a)
		data = r
	}
	return data
}

// CategoriesTags categories or tags
//
// t is constraints.Tag or constraints.Category
func CategoriesTags(ctx context.Context, t ...string) []models.TermsMy {
	tt := ""
	if len(t) > 0 {
		tt = t[0]
	}
	r, err := categoryAndTagsCaches.GetCache(ctx, tt, time.Second, ctx, tt)
	logs.IfError(err, "get category fail")
	return r
}
func AllCategoryTagsNames(ctx context.Context, t ...string) map[string]struct{} {
	tt := ""
	if len(t) > 0 {
		tt = t[0]
	}
	r, err := categoryAndTagsCaches.GetCache(ctx, tt, time.Second, ctx, tt)
	if err != nil {
		logs.Error(err, "get category fail")
		return nil
	}
	return slice.ToMap(r, func(t models.TermsMy) (string, struct{}) {
		return t.Name, struct{}{}
	}, true)
}
