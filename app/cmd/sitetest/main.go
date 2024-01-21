package main

import (
	"flag"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/httptool"
	strings2 "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"github.com/fthvgb1/wp-go/taskPools"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

var reg = regexp.MustCompile(`<a.*href="([^"]+?)".*?>`)
var m = safety.NewMap[string, bool]()

var mut = sync.Mutex{}

func parseHtml(ss string) {
	r := reg.FindAllStringSubmatch(ss, -1)

	for _, href := range r {
		if href[1] == "/" {
			continue
		}
		if strings.ContainsAny(href[1], ".") {
			continue
		}
		if string([]rune(href[1])[0:3]) == "http" {
			continue
		}
		mut.Lock()
		if _, ok := m.Load(href[1]); !ok {
			m.Store(href[1], false)
		}
		mut.Unlock()
	}
}

func siteFetch(c int, u string) {
	u = strings.TrimRight(u, "/")
	ss, code, err := httptool.GetString(u, nil)
	if err != nil || code != http.StatusOK {
		panic(err)
	}
	parseHtml(ss)

	p := taskPools.NewPools(c)
	for {
		m.Range(func(key string, value bool) bool {
			if value {
				return true
			}
			u := strings2.Join(u, key)
			p.Execute(func() {
				ss, code, err := httptool.GetString(u, nil)
				fmt.Println(u, code)
				if err != nil || code != http.StatusOK {
					panic(err)
					return
				}
				parseHtml(ss)
				m.Store(key, true)
			})
			return true
		})
		var x bool
		m.Range(func(key string, value bool) bool {
			if !value {
				x = true
				return false
			}
			return true
		})
		if !x {
			break
		}
	}
	p.Wait()
	m.Flush()
}

var c int
var u string
var t int

func main() {
	flag.IntVar(&c, "c", 10, "concurrency num")
	flag.StringVar(&u, "url", "http://127.0.0.1:8081", "test url")
	flag.IntVar(&t, "t", 1, "request full site times")
	flag.Parse()
	if u == "" {
		fmt.Println("url can't emtpy")
		os.Exit(2)
	}
	if c < 1 {
		fmt.Println("concurrency num must >= 1")
		os.Exit(2)
	}
	if t < 1 {
		for {
			siteFetch(c, u)
		}
	}
	for i := 0; i < t; i++ {
		siteFetch(c, u)
	}
}
