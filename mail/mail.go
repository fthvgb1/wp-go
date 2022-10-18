package mail

import (
	"crypto/tls"
	"fmt"
	"github/fthvgb1/wp-go/vars"
	"gopkg.in/gomail.v2"
	"mime"
	"strings"
)

type AttacheFile struct {
	Name string
	Path string
}

func (f AttacheFile) GetName() string {
	t := strings.Split(f.Path, ".")
	return fmt.Sprintf("%s.%s", f.Name, t[len(t)-1])
}

func SendMail(mailTo []string, subject string, body string, a ...AttacheFile) error {
	m := gomail.NewMessage(
		gomail.SetEncoding(gomail.Base64),
	)
	m.SetHeader("From",
		m.FormatAddress(vars.Conf.Mail.User,
			vars.Conf.Mail.Alias,
		))
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	for _, files := range a {
		m.Attach(files.Path,
			gomail.Rename(files.Name), //重命名
			gomail.SetHeader(map[string][]string{
				"Content-Disposition": {
					fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", files.GetName())),
				},
			}),
		)
	}

	d := gomail.NewDialer(
		vars.Conf.Mail.Host,
		vars.Conf.Mail.Port,
		vars.Conf.Mail.User,
		vars.Conf.Mail.Pass,
	)
	if vars.Conf.Mail.Ssl {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	err := d.DialAndSend(m)
	return err
}