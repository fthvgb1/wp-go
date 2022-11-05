package mail

import (
	"github/fthvgb1/wp-go/config"
	"testing"
)

func TestSendMail(t *testing.T) {
	config.InitConfig("config.yaml")
	type args struct {
		mailTo  []string
		subject string
		body    string
		a       []AttacheFile
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
				a: []AttacheFile{
					{
						Name: "附件",
						Path: "/home/xing/Downloads/favicon.ico",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendMail(tt.args.mailTo, tt.args.subject, tt.args.body, tt.args.a...); (err != nil) != tt.wantErr {
				t.Errorf("SendMail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
