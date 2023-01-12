package wp

type Postmeta struct {
	MetaId    uint64 `db:"meta_id" json:"meta_id" form:"meta_id"`
	PostId    uint64 `db:"post_id" json:"post_id" form:"post_id"`
	MetaKey   string `db:"meta_key" json:"meta_key" form:"meta_key"`
	MetaValue string `db:"meta_value" json:"meta_value" form:"meta_value"`
}

func (p Postmeta) PrimaryKey() string {
	return "meta_id"
}

func (p Postmeta) Table() string {
	return "wp_postmeta"
}
