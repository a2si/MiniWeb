package mwError

const (
	ERROR_NO_ERROR = iota

	ERR_NOW_NO_SUPPORT

	ERR_NETWORK_CONNECT_FAIL
	ERR_NETWORK_TIMEOUT
	ERR_NETWORK_NOT_CONNECT

	ERR_PROXY_NOT_SETTINGS
	ERR_PROXY_MANY_CONNECTIONS
	ERR_PROXY_SOCKS_IDENTD
	ERR_PROXY_REFUSED_FAIL
	ERR_PROXY_ACCOUNT_AUTH_FAIL

	ERR_IO_READ_BY_NEGATIVE
	ERR_IO_TIMEOUT
)

var (
	MsgProxyNoSettings      string = "代理参数设置不完整"
	MsgNotSupport           string = "访问了暂时不支持的功能"
	MsgConnectFail          string = "连接失败(协议不匹配或服务未开启)"
	MsgTimeOut              string = "出现超时错误"
	MsgProxyManyConnections string = "代理暂时不能提供服务: 代理服务器有太多的连接"
	MsgProxyErrID           string = "客户端标识不可用或无法验证"
	MsgProxyRefusedFail     string = "请求被拒绝或失败"
	MsgAccountAuthFail      string = "认证失败, 帐号或密码错误"
	MsgNetworkNotConnect    string = "未连接到服务器"
	MsgIOReadByNegative     string = "IO读入返回负数"
)
