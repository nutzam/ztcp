/*
提供了 TCP 工具包
*/
package main

import (
	"os"
)

// 封装了从命令行来的命令参数
// 不同的子协议处理器对这些参数有不同的理解
type TcpObj struct {
	Target         string   // TCP 请求的目标
	Type           string   // TCP 内容补充说明，可以是 json, text, bin 等
	Body           string   // TCP 请求的内容
	File           *os.File // TCP 请求的内容，采用文件，优先级没有 body 高
	Cookie         string   // HTTP 请求所带的 cookie
	OutputRequest  bool     // 是否输出请求的头部
	OutputResponse bool     // 是否输出响应的头部
}

// 判断当前的请求内容是否为 JSON
func (to *TcpObj) IsJson() bool {
	return "json" == to.Type
}
