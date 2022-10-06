package helper

import "testing"

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
			if got := htmlSpecialChars(tt.args.text, tt.args.flags); got != tt.want {
				t.Errorf("htmlSpecialChars() = %v, want %v", got, tt.want)
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
			if got := htmlSpecialCharsDecode(tt.args.text, tt.args.flags); got != tt.want {
				t.Errorf("htmlSpecialCharsDecode() = %v, want %v", got, tt.want)
			}
		})
	}
}
