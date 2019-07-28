# MiniWeb
golang 实现的一个 HTTP/HTTPS 访问客户端.
写这个轮子的原因是以前用 libcurl 习惯了, 自己写一个g语言的实现.

# 版本
- 3.1.1 初次上传版本, 已可正常使用
- 3.1.2 修改网络架构, 采用分层模型
	1. TCP基础上, 增加代理层
	2. 代理层基础上, 增加TLS支持HTTPS
- 3.1.3 微调框架
	1. 修改GO错误处理方式为全局ERRCODE, ERRMSG模式
	2. 所有可公开常量定义定义到 mwConsts

# 使用
```
package main

import (
	mwWeb "MiniWeb"
	"fmt"
)

func main() {
	w := mwWeb.NewMiniWeb()
	w.GetWebCode("https://www.baidu.com/")
	fmt.Println(w.ResponseText())
}
```

- DevLogs
这个库仅作为开发时使用的日之库, 功能并不完善, 建议使用时删除相关代码

# 架构
- Cookie
对 Cookie 提供支持, 支持存储到文件
- Header
对 ReqHeader, RspHeader 提供支持
- Proxy
	- HTTP
	- HTTPS
	- SOCKS4
	- SOCKS4a
	- SOCKS5
	- socksV5中GSSAPI认证暂未实现, 其他功能均可使用, 但是都在测试状态, 尤其HTTP/HTTPS各种异常很闹人
- Net
网络层 -> 代理层 -> 传输层 -> 网络通讯 -> 网络完毕
	- 网络层
		- TCP 连接
	- 代理层
		- 如果使用代理, 这里则与代理通讯
	- 传输层
		- HTTP默认无动作, HTTPS 则进行TLS通讯
	- 网络通讯
		- HTTP协议通讯
	- 网络完毕
		- TCP 关闭
- Core
实现网页访问的具体功能实现
