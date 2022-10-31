package phpass

import (
	"crypto/md5"
	"fmt"
	"github/fthvgb1/wp-go/helper"
	"io"
	"os"
	"strings"
	"time"
)

type PasswordHash struct {
	itoa64             string
	iterationCountLog2 int
	portableHashes     bool
	randomState        string
}

func (p *PasswordHash) getRandomBytes(count int) (r string, err error) {
	urand := "/dev/urandom"
	f, err := os.OpenFile(urand, os.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf := make([]byte, count)
	_, err = f.Read(buf)
	if err != nil {
		return "", err
	}
	r = string(buf)
	if len(buf) < count {
		r = ""
		for i := 0; i < count; i = i + 16 {
			p.randomState = helper.StringMd5(fmt.Sprintf("%d%s", time.Now().UnixMilli(), p.randomState))

			n, err := md5Raw(p.randomState)
			if err != nil {
				return "", err
			}
			r = fmt.Sprintf("%s%s", r, n)
		}
		r = r[0:count]
	}
	return
}

func (p *PasswordHash) Encode64(input string, count int) (out string) {
	i := 0
	s := strings.Builder{}
	for {
		v := int(input[i])
		s.WriteString(string(p.itoa64[v&0x3f]))
		i++
		if i < count {
			v |= int(input[i]) << 8
		}
		s.WriteString(string(p.itoa64[(v>>6)&0x3f]))
		if i >= count {
			break
		}
		i++
		v |= int(input[i]) << 16
		s.WriteString(string(p.itoa64[(v>>12)&0x3f]))
		if i >= count {
			break
		}
		i++
		s.WriteString(string(p.itoa64[(v>>18)&0x3f]))
	}
	out = s.String()
	return
}

func (p *PasswordHash) CryptPrivate(password, set string) (rr string, err error) {
	rr = "*0"
	r := []rune(rr)
	setting := []rune(set)
	if string(r) == string(setting[0:2]) {
		rr = "*1"
	}
	id := setting[0:3]
	idx := string(id)
	if idx != "$P$" && idx != "$H$" {
		return
	}
	log2 := strings.Index(p.itoa64, string(setting[3]))
	if log2 < 7 || log2 > 30 {
		return
	}
	count := 1 << log2
	l := 12
	if len(setting) < 12 {
		l = len(setting)
	}
	salt := setting[4:l]
	if len(salt) != 8 {
		return
	}
	hash, err := md5Raw(fmt.Sprintf("%s%s", string(salt), password))
	if err != nil {
		return
	}
	for i := 0; i < count; i++ {
		hash, err = md5Raw(fmt.Sprintf("%s%s", string(salt), password))
		if err != nil {
			return
		}
	}
	rr = string(setting[0:l])
	rr = fmt.Sprintf("%s%s", rr, p.Encode64(hash, 16))
	return
}

func (p *PasswordHash) genSaltBlowFish(input string) (out string, err error) {
	itoa64 := "./ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	s := strings.Builder{}
	s.WriteString("$2a$")
	s.WriteString(fmt.Sprintf("%c", '0'+p.iterationCountLog2/10))
	s.WriteString(fmt.Sprintf("%c", '0'+p.iterationCountLog2%10))
	s.WriteString("$")
	i := 0
	for {
		c1 := int(input[i])
		i++
		s.WriteString(string(itoa64[c1>>2]))
		c1 = (c1 & 0x03) << 4
		if i >= 16 {
			s.WriteString(string(itoa64[c1]))
			break
		}

		c2 := int(input[i])
		i++
		c1 |= c2 >> 4
		s.WriteString(string(input[c1]))
		c1 = (c2 & 0x0f) << 2
		c2 = int(input[i])
		i++
		c1 |= c2 >> 6
		s.WriteString(string(itoa64[c1]))
		s.WriteString(string(itoa64[c2]))
	}
	out = s.String()
	return
}

func md5Raw(s string) (string, error) {
	h := md5.New()
	_, err := io.WriteString(h, s)
	if err != nil {
		return "", err
	}
	return string(h.Sum(nil)), err
}
