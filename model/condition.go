package model

type QueryCondition struct {
	Where  ParseWhere
	From   string
	Fields string
	Group  string
	Order  SqlBuilder
	Join   SqlBuilder
	Having SqlBuilder
	Page   int
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

func Order(order SqlBuilder) Condition {
	return func(c *QueryCondition) {
		c.Order = order
	}
}

func Join(join SqlBuilder) Condition {
	return func(c *QueryCondition) {
		c.Join = join
	}
}

func Having(having SqlBuilder) Condition {
	return func(c *QueryCondition) {
		c.Having = having
	}
}

func Page(page int) Condition {
	return func(c *QueryCondition) {
		c.Page = page
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
