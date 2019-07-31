package mwNet

import (
	"crypto/tls"
	"fmt"
	"net"

	DevLogs "github.com/a2si/MiniWeb/DevLogs"
)

func (self *TNet) InitTLS(conn net.Conn, ServerName string) {
	DevLogs.Debug("TNet.InitTLS")
	TlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	c := TlsConfig.Clone()
	c.ServerName = ServerName
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

func (self *TNet) ReadLine() string {
	return self.ioRead.ReadLine()
}

func (self *TNet) ReadBytes(Length int) []byte {
	return self.ioRead.ReadBytes(Length)
}

func (self *TNet) ReadToEOF() []byte {
	return self.ioRead.ReadToEOF()
}

func (self *TNet) ReadChunk() []byte {
	return self.ioRead.ReadChunk()
}
