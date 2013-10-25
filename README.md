# TCP 测试工具

本项目采用 go 语言编写，提供一个 tcp 客户端，专门用作各种 socket 测试，HTTP 测试。尤其是对 AJAX 请求很便利的支持。

# 安装

**检查依赖库**

本项目依赖

* [zgo](https://github.com/nutzam/zgo) : 需要它提供的一些便利函数
* [iconv-go](https://github.com/djimenez/iconv-go) : 需要用它来转换各种编码的 HTTP 页面

它们的安装很简单，请到它们的主页 查看其说明文档。我反正是这么装的:

	# 安装 zgo
	go get github.com/nutzam/zgo
	
	# 安装 iconv-go
	go get github.com/djimenez/iconv-go
	
确保这个代码库已经被装到你机器上了

**自动安装**

	go get github.com/nutzam/ztcp
	
**手动安装**

自己手动从 github 下载代码后，放置在你的 $GOPATH 的 src/github.com/nutzam/ztcp 目录下

	go install github.com/nutzam/ztcp
	
**安装成功的标志**

请检查你的 $GOPATH 是不是

	$GOPATH
		[bin]
			ztcp       # <- 这个是编译好的可执行文件
		[src]
			[github.com]
				[nutzam]
					[ztcp]           # <- 这里是下载下来的源码
						REAME.md
						tcp.go
						tcp_http.go
						...

# 作为 HTTP 客户端

本程序最主要的功能之一是作为一个 HTTP 客户端，下面是针对不通请求方式的具体用法。

## 普通 GET 请求

	# 发送普通 GET 请求到 localhost
	ztcp -http=:8080/app/doit?nm=f8f9
	
	# 发送普通 GET 请求到 localhost:80
	ztcp -http=/app/index.html
	
	# 发送普通 GET 到 www.google.com
	ztcp -http=www.google.com



## 普通 POST 请求
	
	# 发送简单的参数
	ztcp -http=:8080/app/doit -type=form -body="a=10&b=hello"
	
	# 可以把请求内容记录到文件里
	ztcp -http=:8080/app/doit -type=form -f=/home/xiaobai/form.txt
	--------------------------- form.txt 文件的内容就是 :
	a=10&b=hello
	
## JSON 请求
	
	# 发送 JSON 请求
	ztcp -type=json -http=localhost:8080/app/doit -body="{nm:'zozoh'}"
	
	# 采用文件的方式发送 JSON 请求
	ztcp -type=json -http=localhost:8080/app/doit \
	    -f=/home/xiaobai/test.json

## 控制 HTTP 返回的开关
	
	# 是否显示请求，响应，的头部信息，默认为 "none"
	tcp … -out=all,req,resp,none …
	
	    
## 支持的 -type

在 `ztcp` 的参数  `-type` 中，我们根据你的参数会生成对应的 HTTP `Content-Type`。
下面是一个对照表:

	form : application/x-www-form-urlencoded
	text : text/plain
	json : application/json
	css  : text/css
	html : text/html
	png  : image/png
    file : application/octet-stream

## 指定 Header

在 `-http` 模式下，你可以给定一个 JSON 字符串或者文件来定制你要发送的 HTTP HEADER

    # 发送普通 cookie 字符串
    ztcp -http=/app/doit -cookie="CNZZDATA1291011=cnzz_eid;"

    # 发送文件里的 cookie 字符串
    ztcp -http=/app/doit -cookief="mycookie.txt"	
			
## 指定 Cookie

在 `-http` 模式下，为了更加便利的指定 cookie，提供了 cookie 选项。
它的优先级比通过 `-header` 要高。

    # 发送普通 cookie 字符串
    ztcp -http=/app/doit -cookie="CNZZDATA1291011=cnzz_eid;"

    # 发送文件里的 cookie 字符串
    ztcp -http=/app/doit -cookief="mycookie.txt"

