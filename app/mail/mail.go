package mail

import (
	"crypto/tls"
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/soxfmr/gomail"
	"mime"
	"path"
)

type AttacheFile struct {
	Name string
	Path string
}

func SendMail(mailTo []string, subject string, body string, files ...string) error {
	m := gomail.NewMessage(
		gomail.SetEncoding(gomail.Base64),
	)
	c := config.GetConfig()
	m.SetHeader("From",
		m.FormatAddress(c.Mail.User,
			c.Mail.Alias,
		))
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	for _, file := range files {
		_, f := path.Split(file)
		m.Attach(file,
			gomail.Rename(f), //重命名
			gomail.SetHeader(map[string][]string{
				"Content-Disposition": {
					fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", f)),
				},
			}),
		)
	}

	d := gomail.NewDialer(
		c.Mail.Host,
		c.Mail.Port,
		c.Mail.User,
		c.Mail.Pass,
	)
	if c.Mail.InsecureSkipVerify {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	err := d.DialAndSend(m)
	return err
}
