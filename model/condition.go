package model

type QueryCondition struct {
	Where  ParseWhere
	From   string
	Fields string
	Group  string
	Order  SqlBuilder
	Join   SqlBuilder
	Having SqlBuilder
	Limit  int
	Offset int
	In     [][]any
}

func Conditions(fns ...Condition) QueryCondition {
	r := QueryCondition{}
	for _, fn := range fns {
		fn(&r)
	}
	if r.Fields == "" {
		r.Fields = "*"
	}
	return r
}

type Condition func(c *QueryCondition)

func Where(where ParseWhere) Condition {
	return func(c *QueryCondition) {
		c.Where = where
	}
}
func Fields(fields string) Condition {
	return func(c *QueryCondition) {
		c.Fields = fields
	}
}

func From(from string) Condition {
	return func(c *QueryCondition) {
		c.From = from
	}
}

func Group(group string) Condition {
	return func(c *QueryCondition) {
		c.Group = group
	}
}

func Order[T ~[][]string](order T) Condition {
	return func(c *QueryCondition) {
		c.Order = SqlBuilder(order)
	}
}

func Join[T ~[][]string](join T) Condition {
	return func(c *QueryCondition) {
		c.Join = SqlBuilder(join)
	}
}

func Having[T ~[][]string](having T) Condition {
	return func(c *QueryCondition) {
		c.Having = SqlBuilder(having)
	}
}

func Limit(limit int) Condition {
	return func(c *QueryCondition) {
		c.Limit = limit
	}
}

func Offset(offset int) Condition {
	return func(c *QueryCondition) {
		c.Offset = offset
	}
}

func In(in ...[]any) Condition {
	return func(c *QueryCondition) {
		c.In = append(c.In, in...)
	}
}
