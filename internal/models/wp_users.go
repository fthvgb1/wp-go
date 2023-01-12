package models

import "time"

type Users struct {
	Id                uint64    `gorm:"column:ID" db:"ID" json:"ID"`
	UserLogin         string    `gorm:"column:user_login" db:"user_login" json:"user_login"`
	UserPass          string    `gorm:"column:user_pass" db:"user_pass" json:"user_pass"`
	UserNicename      string    `gorm:"column:user_nicename" db:"user_nicename" json:"user_nicename"`
	UserEmail         string    `gorm:"column:user_email" db:"user_email" json:"user_email"`
	UserUrl           string    `gorm:"column:user_url" db:"user_url" json:"user_url"`
	UserRegistered    time.Time `gorm:"column:user_registered" db:"user_registered" json:"user_registered"`
	UserActivationKey string    `gorm:"column:user_activation_key" db:"user_activation_key" json:"user_activation_key"`
	UserStatus        int       `gorm:"column:user_status" db:"user_status" json:"user_status"`
	DisplayName       string    `gorm:"column:display_name" db:"display_name" json:"display_name"`
}

func (u Users) Table() string {
	return "wp_users"
}

func (u Users) PrimaryKey() string {
	return "ID"
}
