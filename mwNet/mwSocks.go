package mwNet

import (
	"encoding/hex"
	"fmt"

	DevLogs "github.com/a2si/MiniWeb/DevLogs"
	mwConfig "github.com/a2si/MiniWeb/mwConfig"
)

/*
	Socks4
		http://www.ufasoft.com/doc/socks4_protocol.htm
	Socks4a
		http://www.openssh.com/txt/socks4a.protocol
*/

func (self *TNet) genSocksConnect(Version int, Host string, Port string) []byte {
	const (
		socks_Command_CONNECT = 0x1
		socks_Command_BIND    = 0x2
		socks_Command_UDP     = 0x3 //Only V5
	)
	var (
		sByte []byte
	)
	switch Version {
	case 4:
		//+----+----+----+----+----+----+----+----+----+----+....+----+
		//| VN | CD | DSTPORT |      DSTIP        | USERID       |NULL|
		//+----+----+----+----+----+----+----+----+----+----+....+----+
		//  1    1      2              4           variable       1
		sByte = []byte{0x4, socks_Command_CONNECT}
		sByte = append(sByte, helpPortToByte(Port)...)
		sByte = append(sByte, helpIPHostToByte(Host)...)
		sByte = append(sByte, 0x0)
		return sByte
	case 41:
		//04 01 0050 00000001 00 7777772e62616964752e636f6d 00
		sByte = []byte{0x4, socks_Command_CONNECT}
		sByte = append(sByte, helpPortToByte(Port)...)
		sByte = append(sByte, []byte{0, 0, 0, 1}...)
		sByte = append(sByte, 0x0)
		sByte = append(sByte, []byte(Host)...)
		return sByte
	case 5:
		sByte = []byte{0x5, socks_Command_CONNECT, 0x0}

		// IP V4 address: X'01'
		//sByte = append(sByte, 0x1, helpIPHostToByte(Host)...)
		// DOMAINNAME: X'03'
		sByte = append(sByte, 0x3, byte(len(Host)))
		sByte = append(sByte, []byte(Host)...)
		// IP V6 address: X'04'
		//sByte = append(sByte, 0x4, helpIPV6HostToByte(Host)...)

		sByte = append(sByte, helpPortToByte(Port)...)
		return sByte
	}
	return nil
}

func (self *TNet) genSocksV5SelectAuth() []byte {
	DevLogs.Debug("TNet.genSocksV5Req")
	const (
		PROXY_METHOD_NOAUTH  byte = 0x0
		PROXY_METHOD_GSSAPI  byte = 0x1
		PROXY_METHOD_USERPWD byte = 0x2
	)
	if mwConfig.GetConfig("Proxy.SocksV5.GSSAPI") == true {
		return []byte{0x5, 0x3, PROXY_METHOD_NOAUTH, PROXY_METHOD_GSSAPI, PROXY_METHOD_USERPWD}
	} else {
		return []byte{0x5, 0x2, PROXY_METHOD_NOAUTH, PROXY_METHOD_USERPWD}
	}
}

func (self *TNet) genSocksV5Login() []byte {
	DevLogs.Debug("TNet.genSocksV5Login")
	var (
		sByte       []byte
		pUser, pPwd string = self.Proxy.GetProxyUserPwd()
	)
	sByte = append(sByte, 0x5)
	sByte = append(sByte, byte(len(pUser)))
	sByte = append(sByte, []byte(pUser)...)
	sByte = append(sByte, byte(len(pPwd)))
	sByte = append(sByte, []byte(pPwd)...)
	fmt.Println(hex.EncodeToString(sByte))
	return sByte
}
