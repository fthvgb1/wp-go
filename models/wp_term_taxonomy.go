package models

var WpTermTaxonomyM = WpTermTaxonomy{}

type WpTermTaxonomy struct {
	model[WpTermTaxonomy]
	TermTaxonomyId uint64 `gorm:"column:term_taxonomy_id" db:"term_taxonomy_id" json:"term_taxonomy_id" form:"term_taxonomy_id"`
	TermId         uint64 `gorm:"column:term_id" db:"term_id" json:"term_id" form:"term_id"`
	Taxonomy       string `gorm:"column:taxonomy" db:"taxonomy" json:"taxonomy" form:"taxonomy"`
	Description    string `gorm:"column:description" db:"description" json:"description" form:"description"`
	Parent         uint64 `gorm:"column:parent" db:"parent" json:"parent" form:"parent"`
	Count          int64  `gorm:"column:count" db:"count" json:"count" form:"count"`
}

func (w WpTermTaxonomy) PrimaryKey() string {
	return "term_taxonomy_id"
}

func (w WpTermTaxonomy) Table() string {
	return "wp_term_taxonomy"
}
