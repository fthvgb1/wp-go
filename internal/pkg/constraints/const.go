package constraints

const (
	Home = iota + 1
	Archive
	Category
	Tag
	Search
	Author
	Detail

	Ok
	Error404
	ParamError
	InternalErr

	Defaults = "default"

	HeadScript   = "headScript"
	FooterScript = "footerScript"
)
