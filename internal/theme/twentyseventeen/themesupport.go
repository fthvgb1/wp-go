package twentyseventeen

type themeSupport struct {
	CustomLineHeight bool           `json:"custom-line-height"`
	StarterContent   StarterContent `json:"starter-content"`
}

type Widgets struct {
	Sidebar1 []string `json:"sidebar-1"`
	Sidebar2 []string `json:"sidebar-2"`
	Sidebar3 []string `json:"sidebar-3"`
}
type About struct {
	Thumbnail string `json:"thumbnail"`
}
type Contact struct {
	Thumbnail string `json:"thumbnail"`
}
type Blog struct {
	Thumbnail string `json:"thumbnail"`
}
type HomepageSection struct {
	Thumbnail string `json:"thumbnail"`
}
type Posts struct {
	Num0            string          `json:"0"`
	About           About           `json:"about"`
	Contact         Contact         `json:"contact"`
	Blog            Blog            `json:"blog"`
	HomepageSection HomepageSection `json:"homepage-section"`
}
type ImageEspresso struct {
	PostTitle string `json:"post_title"`
	File      string `json:"file"`
}
type ImageSandwich struct {
	PostTitle string `json:"post_title"`
	File      string `json:"file"`
}
type ImageCoffee struct {
	PostTitle string `json:"post_title"`
	File      string `json:"file"`
}
type Attachments struct {
	ImageEspresso ImageEspresso `json:"image-espresso"`
	ImageSandwich ImageSandwich `json:"image-sandwich"`
	ImageCoffee   ImageCoffee   `json:"image-coffee"`
}
type Options struct {
	ShowOnFront  string `json:"show_on_front"`
	PageOnFront  string `json:"page_on_front"`
	PageForPosts string `json:"page_for_posts"`
}
type ThemeMods struct {
	Panel1 string `json:"panel_1"`
	Panel2 string `json:"panel_2"`
	Panel3 string `json:"panel_3"`
	Panel4 string `json:"panel_4"`
}
type Top struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}
type Social struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}
type NavMenus struct {
	Top    Top    `json:"top"`
	Social Social `json:"social"`
}
type StarterContent struct {
	Widgets     Widgets     `json:"widgets"`
	Posts       Posts       `json:"posts"`
	Attachments Attachments `json:"attachments"`
	Options     Options     `json:"options"`
	ThemeMods   ThemeMods   `json:"theme_mods"`
	NavMenus    NavMenus    `json:"nav_menus"`
}

var themesupport themeSupport
