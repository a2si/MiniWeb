# MiniWeb
- golang 实现的一个 HTTP/HTTPS 访问客户端.
- 写这个轮子的原因是以前用 libcurl 习惯了, 自己写一个g语言的实现.

# 使用
- DevLogs 
这个库可以删除相关代码, 在此项目中, 仅为开发中使用

[code]
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
[/code]
# 版本
- 3.1.1
曾经是我个人使用的一个库, 在各种语言上均有实现, golang版本内核采用了原生网络库实现, 非 lbcurl
可能会有一些小问题, 下一步先实现代理功能, 然后扩展各种功能, 比如HTTPS证书设置, 而不是忽略ssl错误.

# 架构
- Cookie
对 Cookie 提供支持, 支持存储到文件
- Header
对 ReqHeader, RspHeader 提供支持
- Proxy
实现代理支持
- Net
对 TCP与SSL 网络功能支持
- Core
实现网页访问的具体功能实现

# 缺陷
- 代理模块暂未实现
- 功能覆盖并不完善, 与 libcurl 相比还是有比较大的差距.
