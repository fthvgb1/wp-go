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
	"sync"
	"time"
)

var postContextCache *cache.MapCache[uint64, dao.PostContext]
var archivesCaches *Arch
var categoryAndTagsCaches *cache.MapCache[int, []models.TermsMy]
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

func InitActionsCommonCache() {
	c := config.GetConfig()
	archivesCaches = &Arch{
		mutex: &sync.Mutex{},
		fn:    dao.Archives,
	}

	searchPostIdsCache = cachemanager.MapCacheBy[string](dao.SearchPostIds, c.CacheTime.SearchPostCacheTime)

	postListIdsCache = cachemanager.MapCacheBy[string](dao.SearchPostIds, c.CacheTime.PostListCacheTime)

	monthPostsCache = cachemanager.MapCacheBy[string](dao.MonthPost, c.CacheTime.MonthPostCacheTime)

	postContextCache = cachemanager.MapCacheBy[uint64](dao.GetPostContext, c.CacheTime.ContextPostCacheTime)

	postsCache = cachemanager.MapBatchCacheBy(dao.GetPostsByIds, c.CacheTime.PostDataCacheTime)

	postMetaCache = cachemanager.MapBatchCacheBy(dao.GetPostMetaByPostIds, c.CacheTime.PostDataCacheTime)

	categoryAndTagsCaches = cachemanager.MapCacheBy[int](dao.CategoriesAndTags, c.CacheTime.CategoryCacheTime)

	recentPostsCaches = cache.NewVarCache(dao.RecentPosts, c.CacheTime.RecentPostCacheTime)

	recentCommentsCaches = cache.NewVarCache(dao.RecentComments, c.CacheTime.RecentCommentsCacheTime)

	postCommentCaches = cachemanager.MapCacheBy[uint64](dao.PostComments, c.CacheTime.PostCommentsCacheTime)

	maxPostIdCache = cache.NewVarCache(dao.GetMaxPostId, c.CacheTime.MaxPostIdCacheTime)

	usersCache = cachemanager.MapCacheBy[uint64](dao.GetUserById, c.CacheTime.UserInfoCacheTime)

	usersNameCache = cachemanager.MapCacheBy[string](dao.GetUserByName, c.CacheTime.UserInfoCacheTime)

	commentsCache = cachemanager.MapBatchCacheBy(dao.GetCommentByIds, c.CacheTime.CommentsCacheTime)

	allUsernameCache = cache.NewVarCache(dao.AllUsername, c.CacheTime.UserInfoCacheTime)

	headerImagesCache = cachemanager.MapCacheBy[string](getHeaderImages, c.CacheTime.ThemeHeaderImagCacheTime)

	feedCache = cache.NewVarCache(feed, time.Hour)

	postFeedCache = cachemanager.MapCacheBy[string](postFeed, time.Hour)

	commentsFeedCache = cache.NewVarCache(commentsFeed, time.Hour)

	newCommentCache = cachemanager.MapCacheBy[string, string](nil, 15*time.Minute)

	InitFeed()
}

func Archives(ctx context.Context) (r []models.PostArchive) {
	return archivesCaches.getArchiveCache(ctx)
}

type Arch struct {
	data  []models.PostArchive
	mutex *sync.Mutex
	fn    func(context.Context) ([]models.PostArchive, error)
	month time.Month
}

func (a *Arch) getArchiveCache(ctx context.Context) []models.PostArchive {
	l := len(a.data)
	m := time.Now().Month()
	if l > 0 && a.month != m || l < 1 {
		r, err := a.fn(ctx)
		if err != nil {
			logs.ErrPrintln(err, "set cache err[%s]")
			return nil
		}
		a.mutex.Lock()
		defer a.mutex.Unlock()
		a.month = m
		a.data = r
	}
	return a.data
}

// CategoriesTags categories or tags
//
// t is constraints.Tag or constraints.Category
func CategoriesTags(ctx context.Context, t ...int) []models.TermsMy {
	tt := 0
	if len(t) > 0 {
		tt = t[0]
	}
	r, err := categoryAndTagsCaches.GetCache(ctx, tt, time.Second, ctx, tt)
	logs.ErrPrintln(err, "get category err")
	return r
}
func AllCategoryTagsNames(ctx context.Context, t ...int) map[string]struct{} {
	tt := 0
	if len(t) > 0 {
		tt = t[0]
	}
	r, err := categoryAndTagsCaches.GetCache(ctx, tt, time.Second, ctx, tt)
	logs.ErrPrintln(err, "get category err")
	return slice.ToMap(r, func(t models.TermsMy) (string, struct{}) {
		return t.Name, struct{}{}
	}, true)
}
