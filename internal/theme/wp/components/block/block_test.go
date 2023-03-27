package block

import "testing"

func TestParseBlock(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "t1",
			args: args{
				s: `<!-- wp:categories {"showPostCounts":true,"showEmpty":true} /-->`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ParseBlock(tt.args.s)
		})
	}
}
