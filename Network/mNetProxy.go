package Network

import (
	"net"
	"strconv"

	DevLogs "github.com/MiniWeb/DevLogs"
)

func (self *TNet) InitProxyHttps(conn net.Conn, Host string, Port string) {
	DevLogs.Debug("TNet.InitProxyHttps")
	self.Conn = conn
	self.netIO2BufferIO()
	Header := self.genHttpProxyReqHeader(Host, Port)
	self.SendPacket([]byte(Header))
	if self.ObjError.IsError() {
		return
	}
	var tmpCode, MsgInfo string = self.httpProxyRspHeader()
	MsgCode, _ := strconv.Atoi(tmpCode)

	switch MsgCode {
	case 200:
		return
	case 503: // 503 Too many open connections
		self.ObjError.ProxyError503()
	case 401:
	}
	//fmt.Println(MsgCode, MsgInfo)
}

func (self *TNet) InitProxySocks4(conn net.Conn, Host string, Port string) {
	DevLogs.Debug("TNet.InitProxySocks4")
	self.Conn = conn
	self.netIO2BufferIO()
	Packet := self.genSocksV4Connect(Host, Port)
	self.SendPacket(Packet)
	if self.ObjError.IsError() {
		return
	}
	dwResult := self.getSocksV4ConnectResult()
	if dwResult != 0x5A {
		self.ObjError.ProxyNoSupport()
	}
	return
	switch dwResult {
	case 0x5A: // 允许转发
	case 0x5B: // 拒绝转发，一般性失败
	case 0x5C: // 拒绝转发，SOCKS 4 Server无法连接到SOCS 4 Client所在主机的 IDENT服务
	case 0x5D: // 拒绝转发，请求报文中的USERID与IDENT服务返回值不相符
	}
}

func (self *TNet) InitProxySocks4a(conn net.Conn, Host string, Port string) {
	DevLogs.Debug("TNet.InitProxySocks4a")
	self.Conn = conn
	self.netIO2BufferIO()
	Packet := self.genSocksV4aConnect(Host, Port)
	self.SendPacket(Packet)
	if self.ObjError.IsError() {
		return
	}
	dwResult := self.getSocksV4ConnectResult()
	if dwResult != 0x5A {
		self.ObjError.ProxyNoSupport()
	}
	return
	switch dwResult {
	case 0x5A: // 允许转发
	case 0x5B: // 拒绝转发，一般性失败
	case 0x5C: // 拒绝转发，SOCKS 4 Server无法连接到SOCS 4 Client所在主机的 IDENT服务
	case 0x5D: // 拒绝转发，请求报文中的USERID与IDENT服务返回值不相符
	}
}

func (self *TNet) InitProxySocks5(conn net.Conn, Host string, Port string) {
	DevLogs.Debug("TNet.InitProxySocks5")
	self.Conn = conn
	self.netIO2BufferIO()

	Packet := self.genSocksV5Req()
	self.SendPacket(Packet)
	if self.ObjError.IsError() {
		return
	}
	recvPacket := self.ReadBytes(2)
	if self.ObjError.IsError() || recvPacket[0] != 5 {
		return
	}

	//fmt.Println(hex.EncodeToString(recvPacket))
	var dwStatus byte = recvPacket[1]
	if dwStatus == 1 { // GSSAPI
		self.ObjError.ProxyNoSupport()
		return
	}
	if dwStatus == 2 { // UserName/PassWord Authentication
		Packet = self.genSocksV5Login()
		self.SendPacket(Packet)
		if self.ObjError.IsError() {
			return
		}
		recvPacket = self.ReadBytes(2)
		dwStatus = recvPacket[1]
		// 状态 0 表示认证成功, 否则应关闭连接
		if dwStatus != 0 {
			self.ObjError.ProxyAuthenticationFail()
			return
		}
		dwStatus = 0
	}
	if dwStatus != 0 { // 如果非0, 则表示未知标志
		self.ObjError.ProxyOtherError()
		return
	}

	// Status==0 || Login OK
	Packet = self.genSocksV5Connect(Host, Port)
	self.SendPacket(Packet)
	if self.ObjError.IsError() {
		return
	}
	recvPacket = self.ReadBytes(len(Packet))
	/*
	   o  X'00' succeeded
	   o  X'01' general SOCKS server failure
	   o  X'02' connection not allowed by ruleset
	   o  X'03' Network unreachable
	   o  X'04' Host unreachable
	   o  X'05' Connection refused
	   o  X'06' TTL expired
	   o  X'07' Command not supported
	   o  X'08' Address type not supported
	   o  X'09' to X'FF' unassigned
	*/
	dwStatus = recvPacket[1]
	if dwStatus != 0 {
		self.ObjError.ProxyOtherError()
		return
	}
	//fmt.Println(hex.EncodeToString(recvPacket))
	// SocksV5协议已通过认证, 可以使用了
}
