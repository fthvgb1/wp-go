package model

type QueryCondition struct {
	where  ParseWhere
	from   string
	fields string
	group  string
	order  SqlBuilder
	join   SqlBuilder
	having SqlBuilder
	page   int
	limit  int
	offset int
	in     [][]any
}

func Conditions(fns ...Condition) *QueryCondition {
	r := &QueryCondition{}
	for _, fn := range fns {
		fn(r)
	}
	if r.fields == "" {
		r.fields = "*"
	}
	return r
}

type Condition func(c *QueryCondition)

func Where(where ParseWhere) Condition {
	return func(c *QueryCondition) {
		c.where = where
	}
}
func Fields(fields string) Condition {
	return func(c *QueryCondition) {
		c.fields = fields
	}
}

func From(from string) Condition {
	return func(c *QueryCondition) {
		c.from = from
	}
}

func Group(group string) Condition {
	return func(c *QueryCondition) {
		c.group = group
	}
}

func Order(order SqlBuilder) Condition {
	return func(c *QueryCondition) {
		c.order = order
	}
}

func Join(join SqlBuilder) Condition {
	return func(c *QueryCondition) {
		c.join = join
	}
}

func Having(having SqlBuilder) Condition {
	return func(c *QueryCondition) {
		c.having = having
	}
}

func Page(page int) Condition {
	return func(c *QueryCondition) {
		c.page = page
	}
}

func Limit(limit int) Condition {
	return func(c *QueryCondition) {
		c.limit = limit
	}
}

func Offset(offset int) Condition {
	return func(c *QueryCondition) {
		c.offset = offset
	}
}

func In(in ...[]any) Condition {
	return func(c *QueryCondition) {
		c.in = append(c.in, in...)
	}
}
