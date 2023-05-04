package widget

func CommonArgs() map[string]string {
	return map[string]string{
		"{$before_widget}": `<aside id="%s" class="%s">`,
		"{$after_widget}":  "</aside>",
		"{$before_title}":  `<h2 class="widget-title">`,
		"{$after_title}":   "</h2>",
	}
}
