package phpass

import "testing"

func TestPasswordHash_getRandomBytes(t *testing.T) {
	type fields struct {
		itoa64             string
		iterationCountLog2 int
		portableHashes     bool
		randomState        string
	}
	type args struct {
		count int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantR   string
		wantErr bool
	}{
		{
			name: "t1",
			fields: fields{
				itoa64:             "",
				iterationCountLog2: 0,
				portableHashes:     false,
				randomState:        "",
			},
			args:    args{5},
			wantR:   "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PasswordHash{
				itoa64:             tt.fields.itoa64,
				iterationCountLog2: tt.fields.iterationCountLog2,
				portableHashes:     tt.fields.portableHashes,
				randomState:        tt.fields.randomState,
			}
			_, err := p.getRandomBytes(tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRandomBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
