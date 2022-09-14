package route

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/models"
	"net/http"
	"strconv"
	"sync"
)

var PostsCache sync.Map

func index(c *gin.Context) {
	page := 1
	pageSize := 10
	p := c.Query("paged")
	if pa, err := strconv.Atoi(p); err != nil {
		page = pa
	}
	status := []interface{}{"publish", "private"}
	posts, _, err := models.SimplePagination[models.WpPosts](models.SqlBuilder{{
		"post_type", "post",
	}, {"post_status", "in", ""}}, "ID", page, pageSize, models.SqlBuilder{{"post_date", "desc"}}, nil, status)
	if err != nil {
		return
	}
	var all []uint64
	var allPosts []models.WpPosts
	var needQuery []interface{}
	for _, wpPosts := range posts {
		all = append(all, wpPosts.Id)
		if _, ok := PostsCache.Load(wpPosts.Id); !ok {
			needQuery = append(needQuery, wpPosts.Id)
		}
	}
	if len(needQuery) > 0 {
		rawPosts, err := models.Find[models.WpPosts](models.SqlBuilder{{
			"Id", "in", "",
		}}, "a.*,d.name category_name", "", nil, models.SqlBuilder{{
			"a", "left join", "wp_term_relationships b", "a.Id=b.object_id",
		}, {
			"left join", "wp_term_taxonomy c", "b.term_taxonomy_id=c.term_taxonomy_id",
		}, {
			"left join", "wp_terms d", "c.term_id=d.term_id",
		}}, 0, needQuery)
		if err != nil {
			return
		}
		for _, post := range rawPosts {
			PostsCache.Store(post.Id, post)
		}
	}

	for _, id := range all {
		post, _ := PostsCache.Load(id)
		pp := post.(models.WpPosts)
		allPosts = append(allPosts, pp)
	}
	recent, _ := recentPosts()
	archive, _ := archives()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"posts":       allPosts,
		"options":     models.Options,
		"recentPosts": recent,
		"archives":    archive,
	})
}

func recentPosts() (r []models.WpPosts, err error) {
	r, err = models.Find[models.WpPosts](models.SqlBuilder{{
		"post_type", "post",
	}, {"post_status", "publish"}}, "ID,post_title", "", models.SqlBuilder{{"post_date", "desc"}}, nil, 5)
	return
}

func archives() (r []models.PostArchive, err error) {
	r, err = models.Find[models.PostArchive](models.SqlBuilder{
		{"post_type", "post"}, {"post_status", "publish"},
	}, "YEAR(post_date) AS `year`, MONTH(post_date) AS `month`, count(ID) as posts", "year,month", models.SqlBuilder{{"year", "desc"}, {"month", "desc"}}, nil, 0)
	return
}
