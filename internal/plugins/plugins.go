package plugins

const (
	Home = iota + 1
	Archive
	Category
	Tag
	Search
	Detail

	Ok
	Empty404
	Error
	InternalErr
)

var IndexSceneMap = map[int]struct{}{
	Home:     {},
	Archive:  {},
	Category: {},
	Tag:      {},
	Search:   {},
}
