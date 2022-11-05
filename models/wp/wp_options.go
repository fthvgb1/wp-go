package wp

type WpOptions struct {
	OptionId    uint64 `gorm:"column:option_id" db:"option_id" json:"option_id" form:"option_id"`
	OptionName  string `gorm:"column:option_name" db:"option_name" json:"option_name" form:"option_name"`
	OptionValue string `gorm:"column:option_value" db:"option_value" json:"option_value" form:"option_value"`
	Autoload    string `gorm:"column:autoload" db:"autoload" json:"autoload" form:"autoload"`
}

func (w WpOptions) PrimaryKey() string {
	return "option_id"
}

func (w WpOptions) Table() string {
	return "wp_options"
}
