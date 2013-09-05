package main

import (
	"flag"
	z "github.com/nutzam/zgo"
	"os"
)

// 这里声明了参数
var contentType = flag.String("type", "text", "Content Type of send, 'text' as default")
var httpUrl = flag.String("http", "", "HTTP Target URL")
var body = flag.String("body", "", "Content to send")
var file = flag.String("f", "", "File as the content to Send")
var out = flag.String("out", "none", "Show Header information")
var cookie = flag.String("cookie", "", "for HTTP only")
var cookief = flag.String("cookief", "", "read cookie from file, for HTTP only")

// 程序的主入口
func main() {
	flag.Parse()
	to := new(TcpObj)
	to.Type = *contentType

	// Cookie 放在文件里
	if len(*cookief) > 0 {
		to.Cookie, _ = z.Utf8f(*cookief)
	} else {
		to.Cookie = *cookie
	}
	to.Body = *body
	if len(*file) > 0 {
		f, err := os.Open(z.Ph(*file))
		if nil != err {
			panic(err)
		}
		to.File = f
	}
	switch *out {
	case "all":
		to.OutputRequest = true
		to.OutputResponse = true
	case "req":
		to.OutputRequest = true
	case "resp":
		to.OutputResponse = true
	}

	if len(*httpUrl) > 0 {
		to.Target = *httpUrl
		to.DoHttp()
	} else {
		flag.PrintDefaults()
	}
}
