package wp

type Terms struct {
	TermId    uint64 `gorm:"column:term_id" db:"term_id" json:"term_id" form:"term_id"`
	Name      string `gorm:"column:name" db:"name" json:"name" form:"name"`
	Slug      string `gorm:"column:slug" db:"slug" json:"slug" form:"slug"`
	TermGroup int64  `gorm:"column:term_group" db:"term_group" json:"term_group" form:"term_group"`
}

func (t Terms) PrimaryKey() string {
	return "term_id"
}
func (t Terms) Table() string {
	return "wp_terms"
}

type TermsMy struct {
	Terms
	TermTaxonomy
}

func (t TermsMy) PrimaryKey() string {
	return "term_id"
}
func (t TermsMy) Table() string {
	return "wp_terms"
}
