package model

import (
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/helper/slice"
	"testing"
)

type TermTaxonomy struct {
	TermTaxonomyId uint64 `gorm:"column:term_taxonomy_id" db:"term_taxonomy_id" json:"term_taxonomy_id" form:"term_taxonomy_id"`
	TermId         uint64 `gorm:"column:term_id" db:"term_id" json:"term_id" form:"term_id"`
	Taxonomy       string `gorm:"column:taxonomy" db:"taxonomy" json:"taxonomy" form:"taxonomy"`
	Description    string `gorm:"column:description" db:"description" json:"description" form:"description"`
	Parent         uint64 `gorm:"column:parent" db:"parent" json:"parent" form:"parent"`
	Count          int64  `gorm:"column:count" db:"count" json:"count" form:"count"`
	Term           *models.Terms
}

type CommentMeta struct {
	MetaId    uint64 `db:"meta_id"`
	CommentId uint64 `db:"comment_id"`
	MetaKey   string `db:"meta_key"`
	MetaValue string `db:"meta_value"`
}

var termMyHasOneTerm = RelationHasOne(func(m *TermTaxonomy) uint64 {
	return m.TermTaxonomyId
}, func(p *models.Terms) uint64 {
	return p.TermId
}, func(m *TermTaxonomy, p *models.Terms) {
	m.Term = p
}, Relationship{
	RelationType: HasOne,
	Table:        "wp_terms",
	ForeignKey:   "term_id",
	Local:        "term_id",
})

var postHasManyShip = RelationHasMany(func(m *post) uint64 {
	return m.Id
}, func(p *TermRelationships) uint64 {
	return p.ObjectID
}, func(m *post, i *[]TermRelationships) {
	m.Ships = i
}, Relationship{
	RelationType: HasMany,
	Table:        "wp_term_relationships",
	ForeignKey:   "object_id",
	Local:        "ID",
})

var shipHasManyTermMy = RelationHasMany(func(m *TermRelationships) uint64 {
	return m.TermTaxonomyId
}, func(p *TermTaxonomy) uint64 {
	return p.TermTaxonomyId
}, func(m *TermRelationships, i *[]TermTaxonomy) {
	m.TermTaxonomy = i
}, Relationship{
	RelationType: HasMany,
	Table:        "wp_term_taxonomy",
	ForeignKey:   "term_taxonomy_id",
	Local:        "term_taxonomy_id",
})

func (w TermTaxonomy) PrimaryKey() string {
	return "term_taxonomy_id"
}

func (w TermTaxonomy) Table() string {
	return "wp_term_taxonomy"
}

func postAuthorId(p *post) uint64 {
	return p.PostAuthor
}

func postId(p *post) uint64 {
	return p.Id
}

func userId(u *user) uint64 {
	return u.Id
}

func metasPostId(m *models.PostMeta) uint64 {
	return m.PostId
}

func PostAuthor() (func(any) []any, func(any, any), any, any, Relationship) {
	var u user
	var uu []user
	return GetWithID[post](func(t *post) uint64 {
			return t.PostAuthor
		}),
		SetHasOne(func(p *post, u *user) {
			p.User = u
		}, func(t *post) uint64 {
			return t.PostAuthor
		}, func(u *user) uint64 {
			return u.Id
		}),
		&u, &uu,
		Relationship{
			RelationType: HasOne,
			Table:        "wp_users user",
			ForeignKey:   "ID",
			Local:        "post_author",
		}
}
func PostMetas() (func(any) []any, func(any, any), any, any, Relationship) {
	var u []models.PostMeta
	return GetWithID(func(t *post) any {
			return t.Id
		}), SetHasMany(func(t *post, v *[]models.PostMeta) {
			t.PostMeta = v
		}, func(t *post) uint64 {
			return t.Id
		}, func(m *models.PostMeta) uint64 {
			return m.PostId
		}), &u, &u, Relationship{
			RelationType: HasMany,
			Table:        "wp_postmeta meta",
			ForeignKey:   "post_id",
			Local:        "ID",
		}
}

var postHaveManyTerms = RelationHasMany(func(m *post) uint64 {
	return m.Id
}, func(p *struct {
	ObjectId uint64 `db:"object_id"`
	models.Terms
}) uint64 {
	return p.ObjectId
}, func(m *post, i *[]struct {
	ObjectId uint64 `db:"object_id"`
	models.Terms
}) {
	v := slice.Map(*i, func(t struct {
		ObjectId uint64 `db:"object_id"`
		models.Terms
	}) models.Terms {
		return t.Terms
	})
	m.Terms = &v
}, Relationship{
	RelationType: HasOne,
	Table:        "wp_terms",
	ForeignKey:   "term_id",
	Local:        "term_id",
	Middle: &Relationship{
		RelationType: HasOne,
		Table:        "wp_term_taxonomy taxonomy",
		ForeignKey:   "term_taxonomy_id",
		Local:        "term_taxonomy_id",
		Middle: &Relationship{
			RelationType: HasMany,
			Table:        "wp_term_relationships",
			ForeignKey:   "object_id",
			Local:        "ID",
		},
	},
})

var postHaveManyCommentMetas = func() RelationFn {
	type metas struct {
		CommentPostID uint64 `db:"comment_post_ID"`
		CommentMeta
	}
	return RelationHasMany(func(m *post) uint64 {
		return m.Id
	}, func(p *metas) uint64 {
		return p.CommentPostID
	}, func(m *post, i *[]metas) {
		v := slice.Map(*i, func(t metas) CommentMeta {
			return t.CommentMeta
		})
		m.CommentMetas = &v
	}, Relationship{
		RelationType: HasOne,
		Table:        "wp_commentmeta",
		ForeignKey:   "comment_id",
		Local:        "comment_ID",
		Middle: &Relationship{
			RelationType: HasMany,
			Table:        "wp_comments comments",
			ForeignKey:   "comment_post_ID",
			Local:        "ID",
		},
	})
}()

func Meta2() RelationFn {
	return RelationHasMany(postId, metasPostId, func(m *post, i *[]models.PostMeta) {
		m.PostMeta = i
	}, Relationship{
		RelationType: "hasMany",
		Table:        "wp_postmeta meta",
		ForeignKey:   "post_id",
		Local:        "ID",
	})
}

func PostAuthor2() RelationFn {
	return RelationHasOne(postAuthorId, userId, func(p *post, u *user) {
		p.User = u
	}, Relationship{
		RelationType: "hasOne",
		Table:        "wp_users user",
		ForeignKey:   "ID",
		Local:        "post_author",
	})
}

func TestGets2(t *testing.T) {
	t.Run("one", func(t *testing.T) {
		{
			q := Conditions(
				Where(SqlBuilder{{"posts.id = 190"}}),
				WithCtx(&ctx),
				WithFn(true, true, Conditions(
					Fields("ID,user_login,user_pass"),
				), PostAuthor2()),
				Fields("posts.*"),
				From("wp_posts posts"),
				WithFn(true, true, nil, Meta2()),
				WithFn(true, false, Conditions(
					WithFn(true, false, Conditions(
						WithFn(true, false, nil, termMyHasOneTerm),
					), shipHasManyTermMy),
				), postHasManyShip),
				WithFn(true, false, nil, postHaveManyTerms),
			)
			got, err := Gets[post](ctx, q)
			_ = got
			if err != nil {
				t.Errorf("err:%v", err)
			}
		}
	})
	t.Run("many", func(t *testing.T) {
		{
			q := Conditions(
				Where(SqlBuilder{{"posts.id", "in", ""}}),
				In([]any{190, 3022, 291, 2858}),
				WithCtx(&ctx),
				WithFn(true, false, Conditions(
					Fields("ID,user_login,user_pass"),
				), PostAuthor2()),
				Fields("posts.*"),
				From("wp_posts posts"),
				WithFn(true, false, nil, Meta2()),
				/*WithFn(true, false, Conditions(
					WithFn(true, false, Conditions(
						WithFn(true, false, nil, termMyHasOneTerm),
					), shipHasManyTermMy),
				), postHasManyShip),*/
				WithFn(true, false, nil, postHaveManyTerms),
				WithFn(true, false, nil, postHaveManyCommentMetas),
			)
			got, err := Finds[post](ctx, q)
			_ = got
			if err != nil {
				t.Errorf("err:%v", err)
			}
		}
	})
}
