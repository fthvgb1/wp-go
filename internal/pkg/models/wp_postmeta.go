package models

import (
	"github.com/leeqvip/gophp"
	"github/fthvgb1/wp-go/helper"
)

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

func (p Postmeta) AttachmentMetadata() (r WpAttachmentMetadata, err error) {
	if p.MetaKey == "_wp_attachment_metadata" && p.MetaValue != "" {
		unSerialize, er := gophp.Unserialize([]byte(p.MetaValue))
		if er != nil {
			err = er
			return
		}
		info, ok := unSerialize.(map[string]any)
		if ok {
			r, err = helper.MapToStruct[WpAttachmentMetadata](info)
		}
	}
	return
}
func AttachmentMetadata(s string) (r WpAttachmentMetadata, err error) {
	unSerialize, er := gophp.Unserialize([]byte(s))
	if er != nil {
		err = er
		return
	}
	info, ok := unSerialize.(map[string]any)
	if ok {
		r, err = helper.MapToStruct[WpAttachmentMetadata](info)
	}
	return
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
