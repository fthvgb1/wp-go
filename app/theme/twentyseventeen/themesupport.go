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

type Image struct {
	PostTitle string `json:"post_title"`
	File      string `json:"file"`
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
type Menus struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}

type NavMenus struct {
	Top    Menus `json:"top"`
	Social Menus `json:"social"`
}
type StarterContent struct {
	Widgets     Widgets                      `json:"widgets"`
	Posts       map[string]map[string]string `json:"posts"`
	Attachments map[string]Image             `json:"attachments"`
	Options     Options                      `json:"options"`
	ThemeMods   ThemeMods                    `json:"theme_mods"`
	NavMenus    NavMenus                     `json:"nav_menus"`
}

var themesupport = themeSupport{
	CustomLineHeight: true,
	StarterContent: StarterContent{
		Widgets: Widgets{
			Sidebar1: []string{"text_business_info", "search", "text_about"},
			Sidebar2: []string{"text_business_info"},
			Sidebar3: []string{"text_about", "search"},
		},
		Posts: map[string]map[string]string{
			"0": {
				"home": "home",
			},
			"about": {
				"thumbnail": "{{image-sandwich}}",
			},
			"contact": {
				"thumbnail": "{{image-espresso}}",
			},
			"blog": {
				"thumbnail": "{{image-coffee}}",
			},
			"homepage-section": {
				"thumbnail": "{{image-espresso}}",
			},
		},
		Attachments: map[string]Image{
			"image-espresso": {
				PostTitle: "浓缩咖啡",
				File:      "assets/images/espresso.jpg",
			},
			"image-sandwich": {
				PostTitle: "三明治",
				File:      "assets/images/sandwich.jpg",
			},
			"image-coffee": {
				PostTitle: "咖啡",
				File:      "assets/images/coffee.jpg",
			},
		},
		Options: Options{
			ShowOnFront:  "page",
			PageOnFront:  "{{home}}",
			PageForPosts: "{{blog}}",
		},
		ThemeMods: ThemeMods{
			Panel1: "{{homepage-section}}",
			Panel2: "{{about}}",
			Panel3: "{{blog}}",
			Panel4: "{{contact}}",
		},
		NavMenus: NavMenus{
			Top: Menus{
				Name: "顶部菜单",
				Items: []string{
					"link_home",
					"page_about",
					"page_blog",
					"page_contact",
				},
			},
			Social: Menus{
				Name: "社交网络链接菜单",
				Items: []string{
					"link_yelp",
					"link_facebook",
					"link_twitter",
					"link_instagram",
					"link_email",
				},
			},
		},
	},
}
