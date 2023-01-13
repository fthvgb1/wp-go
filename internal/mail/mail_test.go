package mail

import (
	"github/fthvgb1/wp-go/internal/pkg/config"
	"testing"
)

func TestSendMail(t *testing.T) {
	err := config.InitConfig("../config.yaml")
	if err != nil {
		panic(err)
	}
	type args struct {
		mailTo  []string
		subject string
		body    string
		files   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				mailTo:  []string{"fthvgb1@163.com"},
				subject: "测试发邮件",
				body:    "测试发邮件",
				files:   []string{"/home/xing/Downloads/favicon.ico"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendMail(tt.args.mailTo, tt.args.subject, tt.args.body, tt.args.files...); (err != nil) != tt.wantErr {
				t.Errorf("SendMail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
