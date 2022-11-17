package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/config"
	"github/fthvgb1/wp-go/config/wpconfig"
	"github/fthvgb1/wp-go/logs"
	"github/fthvgb1/wp-go/mail"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"runtime"
	"strings"
	"time"
)

func RecoverAndSendMail(w io.Writer) func(ctx *gin.Context) {
	return gin.CustomRecoveryWithWriter(w, func(ctx *gin.Context, err any) {
		c := ctx.Copy()
		stack := stack(4)
		go func() {
			httpRequest, _ := httputil.DumpRequest(c.Request, true)
			headers := strings.Split(string(httpRequest), "\r\n")
			for idx, header := range headers {
				current := strings.Split(header, ":")
				if current[0] == "Authorization" {
					headers[idx] = current[0] + ": *"
				}
			}
			headersToStr := strings.Join(headers, "<br/>")
			content := `<dl><dt>err:</dt><dd>%v</dd><hr/>
<dt>stack: </dt><dd>%v</dd><hr/>
<dt>headers:  </dt><dd>%s</dd></dl>
`
			content = fmt.Sprintf(content,
				err,
				formatStack(string(stack)),
				headersToStr,
			)

			er := mail.SendMail(
				[]string{config.Conf.Load().Mail.User},
				fmt.Sprintf("%s%s %s 发生错误", fmt.Sprintf(wpconfig.Options.Value("siteurl")), c.FullPath(), time.Now().Format(time.RFC1123Z)), content)

			if er != nil {
				logs.ErrPrintln(er, "recover send mail fail", fmt.Sprintf("%v", err))
			}
		}()
		ctx.AbortWithStatus(http.StatusInternalServerError)
	})
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func formatStack(s string) (r string) {
	ss := strings.Builder{}
	t := strings.Split(s, "\n")
	for i, line := range t {
		if i%2 == 0 {
			ss.WriteString("<dt>")
			ss.WriteString(line)
			ss.WriteString("</dt>")
		} else {
			ss.WriteString("<dd>")
			ss.WriteString(strings.Trim(line, "\t"))
			ss.WriteString("</dd>")
		}
	}
	r = ss.String()
	return
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contain dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}
