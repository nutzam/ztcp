package main

import (
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	z "github.com/nutzam/zgo"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// HTTP 客户端对象
var client http.Client

// 各种 MIME 类型的映射
var mimes map[string]string

// 浏览器 UserAgent 的映射
var user_agents map[string]string

func init() {
	client = http.Client{}
	//...................................................
	mimes = make(map[string]string)
	mimes["form"] = "application/x-www-form-urlencoded"
	mimes["html"] = "text/html"
	mimes["text"] = "text/plain"
	mimes["json"] = "application/json"
	mimes["css"] = "text/css"
	mimes["png"] = "image/png"
	mimes["file"] = "application/octet-stream"
	//...................................................
	user_agents = make(map[string]string)
	user_agents["chrome"] = fmt.Sprint("Mozilla/5.0 ",
		"(Macintosh; Intel Mac OS X 10_8_4) ",
		"AppleWebKit/537.36 (KHTML, like Gecko) ",
		"Chrome/28.0.1500.95 Safari/537.36")

	user_agents["firefox"] = fmt.Sprint("Mozilla/5.0 ",
		"(Macintosh; Intel Mac OS X 10.8; rv:23.0) ",
		"Gecko/20100101 Firefox/23.0")
}

// 按照 HTTP 的方式理解 TcpObj 并执行，将返回打印到标准输出中
func (to *TcpObj) DoHttp() {
	// 如果默认是访问本地的，那么会 : 开头，或者 / 开头
	if strings.HasPrefix(to.Target, ":") ||
		strings.HasPrefix(to.Target, "/") {
		to.Target = "localhost" + to.Target
	}

	req := to.createRequestHeader()

	// 如果要发送 POST 请求，特殊的设置
	switch {
	// 根据 body 参数发送
	case len(to.Body) > 0:
		to.setupRequestContentType(req)
		req.Header.Set("Content-Length", strconv.Itoa(len(to.Body)))
		req.Body = ioutil.NopCloser(strings.NewReader(to.Body))
	// 根据 -f 指定的文件发送
	case nil != to.File:
		to.setupRequestContentType(req)
		req.Header.Set("Content-Length", fmt.Sprint(z.Fszf(to.File)))
		req.Body = to.File
	}

	// 打印头部
	if to.OutputRequest {
		fmt.Println(req.Proto)
		fmt.Println(req.Method, to.Target)
		fmt.Println(sep("-", 80), ":req.HEADER")
		for key, _ := range req.Header {
			fmt.Printf("%20s : %s\n", key, req.Header.Get(key))
		}
		// 打印分隔行
		fmt.Println(sep("-", 80), ":req.BODY")
		// 打印请求体
		switch {
		// 打印请求参数
		case len(to.Body) > 0:
			fmt.Println(to.Body)
		// 打印文件内容
		case nil != to.File:
			fi := z.Fif(to.File)
			fmt.Printf("$>: %s %dbytes (%s) %s\n",
				fi.Mode().String(),
				fi.Size(),
				fi.ModTime().Format(time.StampMilli),
				fi.Name())
		}
		// 打印与响应的分隔符
		fmt.Println()
		fmt.Println(strings.Repeat(">", 80))
		fmt.Println()
	}

	// 执行发送
	r, err := client.Do(req)
	to.printResponse(r, err)
}

func (to *TcpObj) setupRequestContentType(req *http.Request) {
	mime := "text/plain"
	if len(to.Type) > 0 {
		mime = mimes[to.Type]
	}

	req.Header.Set("Content-Type", mime)
}

// 生成 HTTP 请求对象，主要是填充头部信息
func (to *TcpObj) createRequestHeader() *http.Request {
	var req *http.Request
	if len(to.Body) > 0 || nil != to.File {
		req, _ = http.NewRequest("POST", "http://"+to.Target, nil)
	} else {
		req, _ = http.NewRequest("GET", "http://"+to.Target, nil)
	}
	// 计算 Host
	pos := strings.Index(to.Target, "/")
	var host string
	if pos > 0 {
		host = string(to.Target[:pos])
	} else {
		host = to.Target
	}
	pos = strings.Index(host, ":")
	if pos > 0 {
		host = string(host[:pos])
	}
	// 填充其他信息
	req.Header.Set("Accept", "text/plain,text/html,application/json,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", host)
	req.Header.Set("User-Agent", user_agents["chrome"])

	// 填充自定义头部信息
	if len(to.Header) > 0 {
		for key, val := range to.Header {
			req.Header.Set(key, val)
		}
	}

	// 填充 Cookie
	if !z.IsBlank(to.Cookie) {
		req.Header.Set("Cookie", to.Cookie)
	}

	return req
}

// 打印响应
func (to *TcpObj) printResponse(r *http.Response, err error) {
	if nil != err {
		fmt.Println("Error Happend!\n", err)
		return
	}

	// 打印头部
	if to.OutputResponse {
		fmt.Println(r.Proto, r.Status)
		fmt.Println(sep("-", 80), ":resp.HEADER")
		for key, _ := range r.Header {
			fmt.Printf("%20s : %s\n", key, r.Header.Get(key))
		}
		// 打印分隔行
		fmt.Println(sep("-", 80), ":resp.BODY")
	}

	// 分析编码和内容
	rContentType := strings.TrimSpace(r.Header.Get("Content-Type"))
	reg := regexp.MustCompile(`^([a-z/]+)(;[ ]*)(charset=)(.*)$`)
	grps := reg.FindAllStringSubmatch(rContentType, -1)
	var charset string
	//var contentType string
	if nil != grps {
		//contentType = grps[0][1]
		charset = grps[0][4]
	} else {
		charset = "utf8"
		//contentType = "text/html"
	}

	// 打印 body
	bs, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	str, _ := iconv.ConvertString(string(bs), charset, "utf8")

	fmt.Println(str)

}

func sep(s string, n int) string {
	return strings.Repeat(s, n)
}
