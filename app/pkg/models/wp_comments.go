package models

import "time"

type Comments struct {
	CommentId          uint64    `gorm:"column:comment_ID" db:"comment_ID" json:"comment_ID" form:"comment_ID"`
	CommentPostId      uint64    `gorm:"column:comment_post_ID" db:"comment_post_ID" json:"comment_post_ID" form:"comment_post_ID"`
	CommentAuthor      string    `gorm:"column:comment_author" db:"comment_author" json:"comment_author" form:"comment_author"`
	CommentAuthorEmail string    `gorm:"column:comment_author_email" db:"comment_author_email" json:"comment_author_email" form:"comment_author_email"`
	CommentAuthorUrl   string    `gorm:"column:comment_author_url" db:"comment_author_url" json:"comment_author_url" form:"comment_author_url"`
	CommentAuthorIp    string    `gorm:"column:comment_author_IP" db:"comment_author_IP" json:"comment_author_IP" form:"comment_author_IP"`
	CommentDate        time.Time `gorm:"column:comment_date" db:"comment_date" json:"comment_date" form:"comment_date"`
	CommentDateGmt     time.Time `gorm:"column:comment_date_gmt" db:"comment_date_gmt" json:"comment_date_gmt" form:"comment_date_gmt"`
	CommentContent     string    `gorm:"column:comment_content" db:"comment_content" json:"comment_content" form:"comment_content"`
	CommentKarma       int       `gorm:"column:comment_karma" db:"comment_karma" json:"comment_karma" form:"comment_karma"`
	CommentApproved    string    `gorm:"column:comment_approved" db:"comment_approved" json:"comment_approved" form:"comment_approved"`
	CommentAgent       string    `gorm:"column:comment_agent" db:"comment_agent" json:"comment_agent" form:"comment_agent"`
	CommentType        string    `gorm:"column:comment_type" db:"comment_type" json:"comment_type" form:"comment_type"`
	CommentParent      uint64    `gorm:"column:comment_parent" db:"comment_parent" json:"comment_parent" form:"comment_parent"`
	UserId             uint64    `gorm:"column:user_id" db:"user_id" json:"user_id" form:"user_id"`
	//扩展字段
	PostTitle  string    `db:"post_title"`
	UpdateTime time.Time `gorm:"update_time" form:"update_time" json:"update_time" db:"update_time"`
}

func (w Comments) PrimaryKey() string {
	return "comment_ID"
}

func (w Comments) Table() string {
	return "wp_comments"
}

type PostComments struct {
	Comments
	Children []uint64
}
