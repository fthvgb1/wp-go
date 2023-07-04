package wp

import (
	"testing"
)

func Test_themeJson(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "t1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			themeJson()
		})
	}
}
