package model

import (
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"testing"
)

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
				In([]any{190, 3022}),
				WithCtx(&ctx),
				WithFn(true, false, Conditions(
					Fields("ID,user_login,user_pass"),
				), PostAuthor2()),
				Fields("posts.*"),
				From("wp_posts posts"),
				WithFn(true, false, nil, Meta2()),
			)
			got, err := Finds[post](ctx, q)
			_ = got
			if err != nil {
				t.Errorf("err:%v", err)
			}
		}
	})
}
