package models

type PostMeta struct {
	MetaId    uint64 `db:"meta_id" json:"meta_id" form:"meta_id"`
	PostId    uint64 `db:"post_id" json:"post_id" form:"post_id"`
	MetaKey   string `db:"meta_key" json:"meta_key" form:"meta_key"`
	MetaValue string `db:"meta_value" json:"meta_value" form:"meta_value"`
}

func (p PostMeta) PrimaryKey() string {
	return "meta_id"
}

func (p PostMeta) Table() string {
	return "wp_postmeta"
}

type WpAttachmentMetadata struct {
	Width     int                         `json:"width,omitempty"`
	Height    int                         `json:"height,omitempty"`
	File      string                      `json:"file,omitempty"`
	FileSize  int                         `json:"filesize,omitempty"`
	Sizes     map[string]MetaDataFileSize `json:"sizes,omitempty"`
	ImageMeta ImageMeta                   `json:"image_meta"`
}

type ImageMeta struct {
	Aperture         string   `json:"aperture,omitempty"`
	Credit           string   `json:"credit,omitempty"`
	Camera           string   `json:"camera,omitempty"`
	Caption          string   `json:"caption,omitempty"`
	CreatedTimestamp string   `json:"created_timestamp,omitempty"`
	Copyright        string   `json:"copyright,omitempty"`
	FocalLength      string   `json:"focal_length,omitempty"`
	Iso              string   `json:"iso,omitempty"`
	ShutterSpeed     string   `json:"shutter_speed,omitempty"`
	Title            string   `json:"title,omitempty"`
	Orientation      string   `json:"orientation,omitempty"`
	Keywords         []string `json:"keywords,omitempty"`
}

type MetaDataFileSize struct {
	File     string `json:"file,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	MimeType string `json:"mime-type,omitempty"`
	FileSize int    `json:"filesize,omitempty"`
}
