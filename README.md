# MiniWeb
纯 go 实现的 go 中的 libcurl

# 版本
- 3.1.1 初次上传版本, 已可正常使用
- 3.1.2 修改网络架构, 采用分层模型
	1. TCP基础上, 增加代理层
	2. 代理层基础上, 增加TLS支持HTTPS
- 3.1.3 微调框架
	1. 修改GO错误处理方式为全局ERRCODE, ERRMSG模式
	2. 所有可公开常量定义定义到 mwConsts
- 3.1.4 整体修复, 优化调整
	1. 调整, 检查, 优化现有代码, 逻辑
	2. 精简网络逻辑, 优化错误处理
	
# 使用
```
go get github.com/a2si/MiniWeb
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
	- socksV5中GSSAPI认证暂未实现, 其他功能均已正常
- Net
基础网络层 -> 代理层 -> SSL传输层 -> 网络通讯 -> 网络完毕
	- 基础网络层
		- TCP
		- QUIC 暂未支持, 看其是否成为普遍现象
	- 代理层
		- 如果使用代理, 这里则与代理通讯
	- SSL传输层
		- 如果是 HTTPS 需要加入SSL通讯
	- 网络通讯
		- HTTP协议通讯 GET/POST/...
	- 网络完毕
		- TCP 关闭
- Core
实现网页访问的具体功能实现
