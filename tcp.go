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

var header = flag.String("header", "", "HTTP header info, JSON string")
var headerf = flag.String("headerf", "", "HTTP header info in file, JSON string")

var cookie = flag.String("cookie", "", "HTTP cookie info content")
var cookief = flag.String("cookief", "", "HTTP cookie info in file")

// 程序的主入口
func main() {
	flag.Parse()
	to := new(TcpObj)
	to.Type = *contentType
	to.Header = make(map[string]string)

	// 分析 Header
	switch {
	// 从文件读取
	case !z.IsBlank(*headerf):
		str, _ := z.Utf8f(z.Ph(*headerf))
		z.JsonFromString(str, &to.Header)
	// 来自参数
	case !z.IsBlank(*header):
		z.JsonFromString(*header, &to.Header)
	}

	// 分析 cookie
	switch {
	// 从文件读取
	case !z.IsBlank(*cookief):
		to.Cookie, _ = z.Utf8f(*cookief)
	// 来自参数
	default:
		to.Cookie = *cookie
	}

	// 请求体的内容，考虑 -body 或者 -f
	to.Body = *body
	if len(*file) > 0 {
		f, err := os.Open(z.Ph(*file))
		if nil != err {
			panic(err)
		}
		to.File = f
	}

	// 分析输出方式
	switch *out {
	case "all":
		to.OutputRequest = true
		to.OutputResponse = true
	case "req":
		to.OutputRequest = true
	case "resp":
		to.OutputResponse = true
	}

	// 执行
	switch {
	// 发送 HTTP 请求
	case !z.IsBlank(*httpUrl):
		to.Target = *httpUrl
		to.DoHttp()
	// 默认，打印参数说明
	default:
		flag.PrintDefaults()
	}

}
