package models

import "time"

type WpPosts struct {
	Id                  uint64    `gorm:"column:ID" db:"ID" json:"ID" form:"ID"`
	PostAuthor          uint64    `gorm:"column:post_author" db:"post_author" json:"post_author" form:"post_author"`
	PostDate            time.Time `gorm:"column:post_date" db:"post_date" json:"post_date" form:"post_date"`
	PostDateGmt         time.Time `gorm:"column:post_date_gmt" db:"post_date_gmt" json:"post_date_gmt" form:"post_date_gmt"`
	PostContent         string    `gorm:"column:post_content" db:"post_content" json:"post_content" form:"post_content"`
	PostTitle           string    `gorm:"column:post_title" db:"post_title" json:"post_title" form:"post_title"`
	PostExcerpt         string    `gorm:"column:post_excerpt" db:"post_excerpt" json:"post_excerpt" form:"post_excerpt"`
	PostStatus          string    `gorm:"column:post_status" db:"post_status" json:"post_status" form:"post_status"`
	CommentStatus       string    `gorm:"column:comment_status" db:"comment_status" json:"comment_status" form:"comment_status"`
	PingStatus          string    `gorm:"column:ping_status" db:"ping_status" json:"ping_status" form:"ping_status"`
	PostPassword        string    `gorm:"column:post_password" db:"post_password" json:"post_password" form:"post_password"`
	PostName            string    `gorm:"column:post_name" db:"post_name" json:"post_name" form:"post_name"`
	ToPing              string    `gorm:"column:to_ping" db:"to_ping" json:"to_ping" form:"to_ping"`
	Pinged              string    `gorm:"column:pinged" db:"pinged" json:"pinged" form:"pinged"`
	PostModified        time.Time `gorm:"column:post_modified" db:"post_modified" json:"post_modified" form:"post_modified"`
	PostModifiedGmt     time.Time `gorm:"column:post_modified_gmt" db:"post_modified_gmt" json:"post_modified_gmt" form:"post_modified_gmt"`
	PostContentFiltered string    `gorm:"column:post_content_filtered" db:"post_content_filtered" json:"post_content_filtered" form:"post_content_filtered"`
	PostParent          uint64    `gorm:"column:post_parent" db:"post_parent" json:"post_parent" form:"post_parent"`
	Guid                string    `gorm:"column:guid" db:"guid" json:"guid" form:"guid"`
	MenuOrder           int       `gorm:"column:menu_order" db:"menu_order" json:"menu_order" form:"menu_order"`
	PostType            string    `gorm:"column:post_type" db:"post_type" json:"post_type" form:"post_type"`
	PostMimeType        string    `gorm:"column:post_mime_type" db:"post_mime_type" json:"post_mime_type" form:"post_mime_type"`
	CommentCount        int64     `gorm:"column:comment_count" db:"comment_count" json:"comment_count" form:"comment_count"`
}

func (w WpPosts) PrimaryKey() string {
	return "ID"
}

func (w WpPosts) Table() string {
	return "wp_posts"
}
