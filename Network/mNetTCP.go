package Network

import (
	"fmt"
	"net"
	"strings"

	DevLogs "github.com/a2si/MiniWeb/DevLogs"
	mwConst "github.com/a2si/MiniWeb/mwConst"
)

func (self *TNet) StartNetwork(d net.Dialer, Host string, Port string) net.Conn {
	DevLogs.Debug("TNet.StartNetwork")
	var (
		RemoteAddr      string = Host + ":" + Port
		pHost, pPort    string = self.Proxy.GetProxyIP()
		RemoteAddrProxy string = pHost + ":" + pPort
	)

	switch self.Proxy.GetProxyType() {
	case mwConst.PROXY_TYPE_HTTPS:
		RemoteAddr = RemoteAddrProxy
	case mwConst.PROXY_TYPE_SOCKS4:
		RemoteAddr = RemoteAddrProxy
	case mwConst.PROXY_TYPE_SOCKS4A:
		RemoteAddr = RemoteAddrProxy
	case mwConst.PROXY_TYPE_SOCKS5:
		RemoteAddr = RemoteAddrProxy
	}

	conn, err := d.Dial("tcp", RemoteAddr)
	if err != nil {
		str := err.Error()
		//dial tcp 117.90.1.61:9000: connect: connection refused
		if strings.Contains(str, "connect: connection refused") {
			self.ObjError.SocketRemoteClose()
			return nil
		}
		//dial tcp 219.159.38.202:56210: i/o timeout
		if strings.Contains(str, "i/o timeout") {
			self.ObjError.SocketConnectTimeout()
			return nil
		}
		fmt.Println(err)
	}
	return conn
}

func (self *TNet) Close() {
	self.Conn.Close()
}
