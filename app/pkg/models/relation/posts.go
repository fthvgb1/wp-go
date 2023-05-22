package relation

import (
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/model"
)

var PostsWithAuthor = model.RelationHasOne(func(m *models.Posts) uint64 {
	return m.PostAuthor
}, func(p *models.Users) uint64 {
	return p.Id
}, func(m *models.Posts, p *models.Users) {
	m.Author = p
}, model.Relationship{
	RelationType: model.HasOne,
	Table:        "wp_users user",
	ForeignKey:   "ID",
	Local:        "post_author",
})
