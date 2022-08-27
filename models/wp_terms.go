package models

var WpTermsM = WpTerms{}

type WpTerms struct {
	model[WpTerms]
	TermId    uint64 `gorm:"column:term_id" db:"term_id" json:"term_id" form:"term_id"`
	Name      string `gorm:"column:name" db:"name" json:"name" form:"name"`
	Slug      string `gorm:"column:slug" db:"slug" json:"slug" form:"slug"`
	TermGroup int64  `gorm:"column:term_group" db:"term_group" json:"term_group" form:"term_group"`
}

func (t WpTerms) PrimaryKey() string {
	return "term_id"
}
func (t WpTerms) Table() string {
	return "wp_terms"
}
