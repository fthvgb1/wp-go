package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/dao"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/cache/cachemanager"
	"github.com/fthvgb1/wp-go/cache/reload"
	"time"
)

func InitActionsCommonCache() {
	c := config.GetConfig()

	cachemanager.NewMemoryMapCache(nil, dao.SearchPostIds, c.CacheTime.SearchPostCacheTime, "searchPostIds", func() time.Duration {
		return config.GetConfig().CacheTime.SearchPostCacheTime
	})

	cachemanager.NewMemoryMapCache(nil, dao.SearchPostIds, c.CacheTime.PostListCacheTime, "listPostIds", func() time.Duration {
		return config.GetConfig().CacheTime.PostListCacheTime
	})

	cachemanager.NewMemoryMapCache(nil, dao.MonthPost, c.CacheTime.MonthPostCacheTime, "monthPostIds", func() time.Duration {
		return config.GetConfig().CacheTime.MonthPostCacheTime
	})

	cachemanager.NewMemoryMapCache(nil, dao.GetPostContext, c.CacheTime.ContextPostCacheTime, "postContext", func() time.Duration {
		return config.GetConfig().CacheTime.ContextPostCacheTime
	})

	cachemanager.NewMemoryMapCache(dao.GetPostsByIds, nil, c.CacheTime.PostDataCacheTime, "postData", func() time.Duration {
		return config.GetConfig().CacheTime.PostDataCacheTime
	})

	cachemanager.NewMemoryMapCache(dao.GetPostMetaByPostIds, nil, c.CacheTime.PostDataCacheTime, "postMetaData", func() time.Duration {
		return config.GetConfig().CacheTime.PostDataCacheTime
	})

	cachemanager.NewMemoryMapCache(nil, dao.CategoriesAndTags, c.CacheTime.CategoryCacheTime, "categoryAndTagsData", func() time.Duration {
		return config.GetConfig().CacheTime.CategoryCacheTime
	})

	cachemanager.NewVarMemoryCache(dao.RecentPosts, c.CacheTime.RecentPostCacheTime, "recentPosts", func() time.Duration {
		return config.GetConfig().CacheTime.RecentPostCacheTime
	})

	cachemanager.NewVarMemoryCache(RecentComment, c.CacheTime.RecentCommentsCacheTime, "recentComments", func() time.Duration {
		return config.GetConfig().CacheTime.RecentCommentsCacheTime
	})

	cachemanager.NewMemoryMapCache(nil, dao.CommentNum, 30*time.Second, "commentNumber", func() time.Duration {
		return config.GetConfig().CacheTime.CommentsIncreaseUpdateTime
	})

	cachemanager.NewMemoryMapCache(nil, PostTopComments, 30*time.Second, "PostCommentsIds", func() time.Duration {
		return config.GetConfig().CacheTime.CommentsIncreaseUpdateTime
	})
	cachemanager.NewMemoryMapCache(nil, dao.PostTopCommentNum, 30*time.Second, "postTopCommentsNum", func() time.Duration {
		return config.GetConfig().CacheTime.CommentsIncreaseUpdateTime
	})

	cachemanager.NewMemoryMapCache(dao.GetCommentByIds, nil, time.Hour, "postCommentData", func() time.Duration {
		return config.GetConfig().CacheTime.CommentsCacheTime
	})

	cachemanager.NewMemoryMapCache(dao.CommentChildren, nil, time.Minute, "commentChildren", func() time.Duration {
		return config.GetConfig().CacheTime.CommentsIncreaseUpdateTime
	})

	cachemanager.NewVarMemoryCache(dao.GetMaxPostId, c.CacheTime.MaxPostIdCacheTime, "maxPostId", func() time.Duration {
		return config.GetConfig().CacheTime.MaxPostIdCacheTime
	})

	cachemanager.NewMemoryMapCache(nil, dao.GetUserById, c.CacheTime.UserInfoCacheTime, "userData", func() time.Duration {
		return config.GetConfig().CacheTime.UserInfoCacheTime
	})

	cachemanager.NewMemoryMapCache(nil, dao.GetUserByName, c.CacheTime.UserInfoCacheTime, "usernameToUserData", func() time.Duration {
		return config.GetConfig().CacheTime.UserInfoCacheTime
	})

	cachemanager.NewVarMemoryCache(dao.AllUsername, c.CacheTime.UserInfoCacheTime, "allUsername", func() time.Duration {
		return config.GetConfig().CacheTime.UserInfoCacheTime
	})

	cachemanager.NewVarMemoryCache(SiteFeed, time.Hour, "siteFeed")

	cachemanager.NewMemoryMapCache(nil, PostFeed, time.Hour, "postFeed")

	cachemanager.NewVarMemoryCache(CommentsFeed, time.Hour, "commentsFeed")

	cachemanager.NewMemoryMapCache[string, string](nil, nil, 15*time.Minute, "NewComment")

	InitFeed()
}

type Arch struct {
	data  []models.PostArchive
	fn    func(context.Context) ([]models.PostArchive, error)
	month time.Month
}

var arch = reload.Vars(Arch{
	fn: dao.Archives,
}, "archives-year-month-data")

func Archives(ctx context.Context) []models.PostArchive {
	a := arch.Load()
	data := a.data
	l := len(data)
	m := time.Now().Month()
	if l < 1 || a.month != m {
		r, err := a.fn(ctx)
		if err != nil {
			logs.Error(err, "set cache Archives fail")
			return nil
		}
		a.month = m
		a.data = r
		arch.Store(a)
		data = r
	}
	return data
}
