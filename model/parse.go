package model

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"strconv"
	"strings"
)

func (w SqlBuilder) parseWhereField(ss []string, s *strings.Builder) {
	if strings.Contains(ss[0], ".") && !strings.Contains(ss[0], "(") {
		x := slice.Map(strings.Split(ss[0], "."), func(t string) string {
			return str.Join("`", t, "`")
		})
		s.WriteString(strings.Join(x, "."))
	} else if !strings.Contains(ss[0], ".") && !strings.Contains(ss[0], "(") {
		s.WriteString("`")
		s.WriteString(ss[0])
		s.WriteString("`")
	} else {
		s.WriteString(ss[0])
	}
}

func (w SqlBuilder) parseIn(ss []string, s *strings.Builder, c *int, args *[]any, in *[][]any) (t bool) {
	if slice.IsContained([]string{"in", "not in"}, strings.ToLower(ss[1])) && len(*in) > 0 {
		s.WriteString(" (")
		x := strings.TrimRight(strings.Repeat("?,", len((*in)[*c])), ",")
		s.WriteString(x)
		*args = append(*args, (*in)[*c]...)
		s.WriteString(")")
		*c++
		t = true
	}
	return t
}

func (w SqlBuilder) parseType(ss []string, args *[]any) error {
	if len(ss) == 4 && ss[3] == "int" {
		i, err := strconv.ParseInt(ss[2], 10, 64)
		if err != nil {
			return err
		}
		*args = append(*args, i)
	} else if len(ss) == 4 && ss[3] == "float" {
		i, err := strconv.ParseFloat(ss[2], 64)
		if err != nil {
			return err
		}
		*args = append(*args, i)
	} else {
		*args = append(*args, ss[2])
	}
	return nil
}

// ParseWhere 解析为where条件，支持3种风格,具体用法参照query_test中的 Find 的测试方法
//
// 1. 1个为一组 {{"field operator value"}} 为纯字符串条件，不对参数做处理
//
// 2. 2个为一组 {{"field1","value1"},{"field2","value2"}} => where field1='value1' and field2='value2'
//
// 3. 3个或4个为一组 {{"field","operator","value"[,"int|float"]}} =>  where field operator 'string'|int|float
//
//	{{"a",">","1","int"}} => where 'a'> 1
//
//	{{ "a",">","1"}} => where 'a'>'1'
//
// 另外如果是操作符为in的话为 {{"field","in",""}} => where field in (?,..) in的条件传给 in参数
//
// 4. 5的倍数为一组{{"and|or","field","operator","value","int|float"}}会忽然掉第一组的and|or
//
//	{{"and","field","=","value1","","and","field","=","value2",""}} => where (field = 'value1' and field = 'value2')
//
//	{{"and","field","=","num1","int","or","field","=","num2","int"}} => where (field = num1 or field = num2')
func (w SqlBuilder) ParseWhere(in *[][]any) (string, []any, error) {
	var s strings.Builder
	args := make([]any, 0, len(w))
	c := 0
	for _, ss := range w {
		switch len(ss) {
		case 1:
			s.WriteString(ss[0])
			s.WriteString(" and ")
		case 2:
			w.parseWhereField(ss, &s)
			s.WriteString("=? and ")
			args = append(args, ss[1])
		case 3, 4:
			w.parseWhereField(ss, &s)
			s.WriteString(ss[1])
			if w.parseIn(ss, &s, &c, &args, in) {
				s.WriteString(" and ")
				continue
			}
			s.WriteString(" ? and ")
			err := w.parseType(ss, &args)
			if err != nil {
				return "", nil, err
			}
		}

		if len(ss) >= 5 && len(ss)%5 == 0 {
			j := len(ss) / 5
			for i := 0; i < j; i++ {
				start := i * 5
				end := start + 5
				st := s.String()
				if strings.Contains(st, "and ") && ss[start] == "or" {
					st = strings.TrimRight(st, "and ")
					s.Reset()
					s.WriteString(st)
					s.WriteString(" ")
					s.WriteString(ss[start])
					s.WriteString(" ")
				}
				if i == 0 {
					s.WriteString("( ")
				}
				w.parseWhereField(ss[start+1:end], &s)
				s.WriteString(ss[start+2])
				if w.parseIn(ss[start+1:end], &s, &c, &args, in) {
					s.WriteString(" and ")
					continue
				}
				s.WriteString(" ? and ")
				err := w.parseType(ss[start+1:end], &args)
				if err != nil {
					return "", nil, err
				}
			}
			st := s.String()
			st = strings.TrimRight(st, "and ")
			s.Reset()
			s.WriteString(st)
			s.WriteString(") and ")
		}
	}
	ss := strings.TrimRight(s.String(), "and ")
	if ss != "" {
		s.Reset()
		s.WriteString(" where ")
		s.WriteString(ss)
		ss = s.String()
	}
	if len(*in) > c {
		*in = (*in)[c:]
	}
	return ss, args, nil
}

func (w SqlBuilder) parseOrderBy() string {
	s := strings.Builder{}
	for _, ss := range w {
		if len(ss) == 2 && ss[0] != "" && slice.IsContained([]string{"asc", "desc"}, strings.ToLower(ss[1])) {
			s.WriteString(" ")
			s.WriteString(ss[0])
			s.WriteString(" ")
			s.WriteString(ss[1])
			s.WriteString(",")
		}
	}
	ss := strings.TrimRight(s.String(), ",")
	if ss != "" {
		s.Reset()
		s.WriteString(" order by ")
		s.WriteString(ss)
		ss = s.String()
	}
	return ss
}
func (w SqlBuilder) parseJoin() string {
	s := strings.Builder{}
	for _, ss := range w {
		l := len(ss)
		for j := 0; j < l; j++ {
			s.WriteString(" ")
			if (l == 4 && j == 3) || (l == 3 && j == 2) {
				s.WriteString("on ")
			}
			s.WriteString(ss[j])
			s.WriteString(" ")
		}

	}
	return s.String()
}
