package html

import (
	"html/template"
	"reflect"
	"testing"
)

func Test_htmlSpecialChars(t *testing.T) {
	type args struct {
		text  string
		flags int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{text: "<a href='test'>Test</a>", flags: EntQuotes},
			want: "&lt;a href=&#039;test&#039;&gt;Test&lt;/a&gt;",
		}, {
			name: "t2",
			args: args{text: "<a href='test'>Test</a>", flags: EntCompat},
			want: "&lt;a href='test'&gt;Test&lt;/a&gt;",
		}, {
			name: "t3",
			args: args{text: "<a href='test'>T est</a>", flags: EntCompat | EntSpace},
			want: "&lt;a&nbsp;href='test'&gt;T&nbsp;est&lt;/a&gt;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SpecialChars(tt.args.text, tt.args.flags); got != tt.want {
				t.Errorf("SpecialChars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_htmlSpecialCharsDecode(t *testing.T) {
	type args struct {
		text  string
		flags int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				text:  "&lt;a href='test'&gt;Test&lt;/a&gt;",
				flags: EntCompat,
			},
			want: "<a href='test'>Test</a>",
		}, {
			name: "t2",
			args: args{
				text:  "&lt;a href=&#039;test&#039;&gt;Test&lt;/a&gt;",
				flags: EntQuotes,
			},
			want: "<a href='test'>Test</a>",
		}, {
			name: "t3",
			args: args{
				text:  "<p>this -&gt; &quot;</p>\n",
				flags: EntNoQuotes,
			},
			want: "<p>this -> &quot;</p>\n",
		}, {
			name: "t4",
			args: args{
				text:  "<p>this -&gt; &quot;</p>\n",
				flags: EntCompat,
			},
			want: "<p>this -> \"</p>\n",
		}, {
			name: "t5",
			args: args{
				text:  "<p>this -&gt;&nbsp;&quot;</p>\n",
				flags: EntCompat | EntSpace,
			},
			want: "<p>this -> \"</p>\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SpecialCharsDecode(tt.args.text, tt.args.flags); got != tt.want {
				t.Errorf("SpecialCharsDecode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStripTags(t *testing.T) {
	type args struct {
		str       string
		allowable string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				str:       "<p>ppppp<span>ffff</span></p><img />",
				allowable: "<p><img>",
			},
			want: "<p>pppppffff</p><img />",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripTags(tt.args.str, tt.args.allowable); got != tt.want {
				t.Errorf("StripTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStripTagsX(t *testing.T) {
	type args struct {
		str       string
		allowable string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{
				str:       "<p>ppppp<span>ffff</span></p><img />",
				allowable: "<p><img>",
			},
			want: "<p>pppppffff</p><img />",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripTagsX(tt.args.str, tt.args.allowable); got != tt.want {
				t.Errorf("StripTagsX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkStripTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StripTags(`<p>ppppp<span>ffff</span></p><img />`, "<p><img>")
	}
}
func BenchmarkStripTagsX(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StripTagsX(`<p>ppppp<span>ffff</span></p><img />`, "<p><img>")
	}
}

func TestCloseHtmlTag(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "t1",
			args: args{str: `<pre class="wp-block-preformatted">GRANT privileges ON databasename.tablename TO 'username'@'h...<p class="read-more"><a href="/p/305">继续阅读</a></p>`},
			want: "</pre>",
		},
		{
			name: "t2",
			args: args{str: `<pre><div>`},
			want: "</div></pre>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CloseTag(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloseTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clearTag(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "t1",
			args: args{s: []string{"<pre>", "<p>", "<span>", "</span>"}},
			want: []string{"<pre>", "<p>"},
		},
		{
			name: "t2",
			args: args{s: []string{"<pre>", "</pre>", "<div>", "<span>", "</span>"}},
			want: []string{"<div>"},
		},
		{
			name: "t3",
			args: args{s: []string{"<pre>", "</pre>"}},
			want: []string{},
		},
		{
			name: "t4",
			args: args{s: []string{"<pre>", "<p>"}},
			want: []string{"<pre>", "<p>"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnClosedTag(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnClosedTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderedHtml(t *testing.T) {
	type args struct {
		t    *template.Template
		data map[string]any
	}
	tests := []struct {
		name    string
		args    args
		wantR   string
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				t: func() *template.Template {
					tt, err := template.ParseFiles("./a.gohtml")
					if err != nil {
						panic(err)
					}
					return tt
				}(),
				data: map[string]any{
					"xx": "oo",
				},
			},
			wantR:   "oo",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := RenderedHtml(tt.args.t, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderedHtml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("RenderedHtml() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
