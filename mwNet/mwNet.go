package mwNet

import (
	"net"
	"strconv"
	"time"

	DevLogs "github.com/a2si/MiniWeb/DevLogs"
	mwProxy "github.com/a2si/MiniWeb/Proxy"
	mwConfig "github.com/a2si/MiniWeb/mwConfig"
	mwError "github.com/a2si/MiniWeb/mwError"
)

type TNet struct {
	ObjError  *mwError.TError // Error Object
	Proxy     *mwProxy.TProxy // Proxy Module
	Conn      net.Conn
	ioRead    *TIOBuffer
	timeOutRW time.Duration
}

func NewNet(errObj *mwError.TError, p *mwProxy.TProxy) *TNet {
	Obj := &TNet{
		ObjError:  errObj,
		Proxy:     p,
		timeOutRW: 0,
	}
	return Obj
}

func (self *TNet) SetTimeOut(TimeOut time.Duration) {
	self.timeOutRW = TimeOut
	self.Conn.SetDeadline(time.Now().Add(TimeOut * time.Second))
}

func (self *TNet) SendPacket(Packet []byte) {
	var (
		needSend int = len(Packet)
		nowSend  int = 0
		err      error
	)
	// 发送空数据, 直接返回
	if needSend == 0 {
		return
	}
	// 检查网络状态
	if self.Conn == nil {
		self.ObjError.SocketNetWorkNoConnect()
		return
	}
	//self.Conn.SetWriteDeadline(time.Now().Add(self.timeOutRW))
	nowSend, err = self.Conn.Write(Packet)
	// 出错, 设置错误, 返回
	if err != nil {
		self.ObjError.SocketIOError(err)
		return
	}
	// 正常发送, 无异常
	if nowSend == needSend {
		return
	}
	// 发送不全, 补充发送
	if nowSend != needSend {
		needSend = needSend - nowSend
		self.SendPacket(Packet[needSend:])
		return
	}
	//发送出现未知错误
	if needSend < needSend {
		self.ObjError.SetErrorCode(-1)
		self.ObjError.SetErrorMsg("Net.SendPacket.Send: 已发送 > needSend")
		return
	}

}

func (self *TNet) Close() {
	self.Conn.Close()
}

func (self *TNet) IsIOClosed() bool {
	return self.ioRead.IsClose()
}

func (self *TNet) netIO2BufferIO() {
	self.ioRead = NewIOBuffer(self.ObjError, self.Conn)
}

func (self *TNet) StartNetwork(d net.Dialer, Host string, Port string) net.Conn {
	DevLogs.Debug("TNet.StartNetwork")
	var (
		RemoteAddr   string = Host + ":" + Port
		pHost, pPort string = self.Proxy.GetProxyIP()
		ProxyAddr    string = pHost + ":" + pPort
		ProxyType    int    = self.Proxy.GetProxyType()
		RawConn      net.Conn
	)

	// 网络层 startNetwork
	// self.Proxy.Enable()
	if ProxyType != mwProxy.PROXY_TYPE_NONE {
		if len(ProxyAddr) < 3 { // 代理未设置
			self.ObjError.ProxyNoSettings()
			return nil
		}
		RemoteAddr = ProxyAddr
	}
	// 连接底层网络模型 TCP/QUIC(HTTP V3 ?)
	if mwConfig.GetConfig("NetBase.QUIC") == false {
		RawConn = self.netStartTCP(d, RemoteAddr)
	} else {
		self.ObjError.SocketNotSupport()
	}
	if self.ObjError.IsError() || RawConn == nil {
		// RawConn == nil 要么ObjSet, 要么IsError. 除非内存异常
		return nil
	}

	self.Conn = RawConn
	self.netIO2BufferIO()

	// 代理层 基于 RawConn
	if ProxyType != mwProxy.PROXY_TYPE_NONE { // self.Proxy.Enable()
		switch ProxyType {
		case mwProxy.PROXY_TYPE_HTTP:
		/*
			WebCore.buildReqHeader 检测/修改 HTTP 请求
				GET /PATH  ==> GET HTTP://WWW.XXX.COM/PATH
		*/
		case mwProxy.PROXY_TYPE_HTTPS:
			/*
				CONNECT 连接后, 发送正文或TLS认证
				Host IP/HostName
			*/
			self.initProxyHttps(Host2IP(Host), Port)
		case mwProxy.PROXY_TYPE_SOCKS4:
			/*
				Socks4 协议
				Host IP
			*/
			self.initProxySocks4(Host2IP(Host), Port)
		case mwProxy.PROXY_TYPE_SOCKS4A:
			/*
				Socks4a 协议
				Host IP/HostName
			*/
			self.initProxySocks4a(Host, Port)
		case mwProxy.PROXY_TYPE_SOCKS5:
			/*
				Socks5 协议
				Host IP/HostName
			*/
			self.initProxySocks5(Host, Port)
			//self.initProxySocks5(mwNet.Host2IP(Host), Port)
		}
	}
	if self.ObjError.IsError() {
		return nil
	}

	return self.Conn
}

func (self *TNet) netStartTCP(d net.Dialer, RemoteAddr string) net.Conn {
	conn, err := d.Dial("tcp", RemoteAddr)
	if err != nil {
		self.ObjError.SocketStartError(err)
	}
	return conn
}

func (self *TNet) initProxyHttps(Host string, Port string) {
	DevLogs.Debug("TNet.InitProxyHttps")
	Header := self.genHttpProxyReqHeader(Host, Port)
	self.SendPacket([]byte(Header))
	if self.ObjError.IsError() {
		return
	}
	var tmpCode, MsgInfo string = self.httpProxyRspHeader()
	MsgCode, _ := strconv.Atoi(tmpCode)
	if MsgCode == 200 {
		return
	}
	self.ObjError.ProxyCheckError(MsgCode, MsgInfo)
}

func (self *TNet) initProxySocks4(Host string, Port string) {
	DevLogs.Debug("TNet.InitProxySocks4")
	Packet := self.genSocksConnect(4, Host, Port)
	self.SendPacket(Packet)
	if self.ObjError.IsError() {
		return
	}
	sByte := self.ReadBytes(8)
	if self.ObjError.IsError() {
		return
	}
	dwResult := sByte[1]
	//dwResult := self.getSocksV4ConnectResult()
	switch dwResult {
	case 0x5A: // 允许转发
		return
	case 0x5B: // 请求被拒绝或失败
		self.ObjError.ProxyRefusedFail()
	case 0x5C: // 拒绝转发，SOCKS 4 Server无法连接到SOCS 4 Client所在主机的 IDENT服务
		self.ObjError.ProxyIDEntdError()
	case 0x5D: // 拒绝转发，请求报文中的USERID与IDENT服务返回值不相符
		self.ObjError.ProxyIDEntdError()
	}
}

func (self *TNet) initProxySocks4a(Host string, Port string) {
	DevLogs.Debug("TNet.InitProxySocks4a")
	Packet := self.genSocksConnect(41, Host, Port)
	self.SendPacket(Packet)
	if self.ObjError.IsError() {
		return
	}
	sByte := self.ReadBytes(8)
	if self.ObjError.IsError() {
		return
	}
	dwResult := sByte[1]
	//dwResult := self.getSocksV4ConnectResult()
	switch dwResult {
	case 0x5A: // 允许转发
		return
	case 0x5B: // 请求被拒绝或失败
		self.ObjError.ProxyRefusedFail()
	case 0x5C: // 拒绝转发，SOCKS 4 Server无法连接到SOCS 4 Client所在主机的 IDENT服务
		self.ObjError.ProxyIDEntdError()
	case 0x5D: // 拒绝转发，请求报文中的USERID与IDENT服务返回值不相符
		self.ObjError.ProxyIDEntdError()
	}
}

func (self *TNet) initProxySocks5(Host string, Port string) {
	DevLogs.Debug("TNet.InitProxySocks5")
	var (
		Auth       bool = self.proxySocksV5Authentication()
		Packet     []byte
		packetSize int
	)

	if Auth == true {
		Packet = self.genSocksConnect(5, Host, Port)
		self.SendPacket(Packet)
		if self.ObjError.IsError() {
			return
		}
		packetSize = len(Packet)
		Packet = self.ReadBytes(packetSize)
		if len(Packet) != packetSize {
			return
		}
		switch Packet[1] {
		case 0: //X'00' succeeded
			return
		case 1: //X'01' general SOCKS server failure（普通的SOCKS服务器请求失败）
		case 2: //X'02' connection not allowed by ruleset（现有的规则不允许的连接）
		case 3: //X'03' Network unreachable（网络不可达）
		case 4: //X'04' Host unreachable（主机不可达）
		case 5: //X'05' Connection refused
		case 6: //X'06' TTL expired（TTL超时）
		case 7: //X'07' Command not supported（不支持的命令）
		case 8: //X'08' Address type not supported（不支持的地址类型）
		case 9: //X'09' to X'FF' unassigned（未定义）
		default:
			self.ObjError.SocketNotSupport()
		}
	}
}

func (self *TNet) proxySocksV5Authentication() bool {
	var (
		Packet []byte = self.genSocksV5SelectAuth()
	)
	self.SendPacket(Packet)
	if self.ObjError.IsError() {
		return false
	}
	// Recv Auth Method
	Packet = self.ReadBytes(2) // |VER | METHOD |
	if self.ObjError.IsError() || Packet[0] != 5 {
		return false
	}
	AuthMethod := Packet[1]
	switch AuthMethod {
	case 0: // NO AUTHENTICATION REQUIRED

		return true

	case 1: // GSSAPI

		self.ObjError.SocketNotSupport()
		return false

	case 2: // USERNAME/PASSWORD

		Packet = self.genSocksV5Login()
		self.SendPacket(Packet)
		if self.ObjError.IsError() {
			return false
		}

		Packet = self.genSocksV5Login()
		self.SendPacket(Packet)
		if self.ObjError.IsError() {
			return false
		}
		Packet = self.ReadBytes(2)
		if Packet[1] == 0 {
			return true
		}
		self.ObjError.ProxyAuthenticationFail()
		return false

	default:

		self.ObjError.SocketNotSupport()
		return false
	}
}
