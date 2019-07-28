package Network

import (
	"crypto/tls"
	"fmt"
	"net"

	DevLogs "github.com/a2si/MiniWeb/DevLogs"
)

func (self *TNet) InitTCP(conn net.Conn) {
	DevLogs.Debug("TNet.InitTCP")
	self.Conn = conn
	self.netIO2BufferIO()
}

func (self *TNet) InitTLS(conn net.Conn) {
	DevLogs.Debug("TNet.InitTLS")
	TlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	c := TlsConfig.Clone()
	c.ServerName = "www.baidu.com"
	TlsConfig = c

	TlsConn := tls.Client(conn, TlsConfig)
	err := TlsConn.Handshake()
	if err != nil {
		fmt.Println("TNet.InitTLS.Handshake: ", err)
		conn.Close()
		return
	}
	self.Conn = TlsConn
	self.netIO2BufferIO()
}
