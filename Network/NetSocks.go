package Network

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	DevLogs "github.com/MiniWeb/DevLogs"
)

func (self *TNet) genSocksV4Connect(Host string, Port string) []byte {
	/*
		+----+----+----+----+----+----+----+----+----+----+...+----+
		| VN | CD | DSTPORT |      DSTIP        | USERID      |NULL|
		+----+----+----+----+----+----+----+----+----+----+...+----+
		   1    1      2              4           variable       1
		CD      SOCKS命令，可取如下值:
		        0x01    CONNECT
		        0x02    BIND
	*/
	DevLogs.Debug("TNet.genSocksV4Connect")
	var (
		sByte []byte
		//pUser, pPwd string = self.Proxy.GetProxyUserPwd()
	)
	sByte = append(sByte, 0x4)
	sByte = append(sByte, 0x1)
	sByte = append(sByte, helpPortToByte(Port)...)
	sByte = append(sByte, helpIPHostToByte(Host)...)
	sByte = append(sByte, 0x0) // len(pUser)
	//sByte = append(sByte, []byte(pUser))
	//fmt.Println(hex.EncodeToString(sByte))
	return sByte
}

func (self *TNet) getSocksV4ConnectResult() int {
	/*
		+----+----+----+----+----+----+----+----+
		| VN | CD | DSTPORT |      DSTIP        |
		+----+----+----+----+----+----+----+----+
		   1    1      2              4
		VN      响应应该为0x00而不是0x04
		CD      可取如下值:
		        0x5A    允许转发
		        0x5B    拒绝转发，一般性失败
		        0x5C    拒绝转发，SOCKS 4 Server无法连接到SOCS 4 Client所在主机的 IDENT服务
		        0x5D    拒绝转发，请求报文中的USERID与IDENT服务返回值不相符
		DSTPORT CD相关的端口信息
		DSTIP   CD相关的地址信息
	*/
	DevLogs.Debug("TNet.getSocksV4ConnectResult")
	var (
		sByte  []byte = self.ReadBytes(8)
		Result int    = int(sByte[1])
		//pUser, pPwd string = self.Proxy.GetProxyUserPwd()
	)
	//fmt.Println(hex.EncodeToString(sByte))
	return Result
}

func (self *TNet) genSocksV4aConnect(Host string, Port string) []byte {
	DevLogs.Debug("TNet.genSocksV4aConnect")
	var (
		sByte []byte
		//pUser, pPwd string = self.Proxy.GetProxyUserPwd()
	)
	sByte = append(sByte, 0x4)
	sByte = append(sByte, 0x1)
	sByte = append(sByte, helpPortToByte(Port)...)

	sByte = append(sByte, 0x0)
	sByte = append(sByte, 0x0)
	sByte = append(sByte, 0x0)
	sByte = append(sByte, 0x1)

	sByte = append(sByte, 0x0) // len(pUser)
	//sByte = append(sByte, []byte(pUser))
	sByte = append(sByte, byte(len(Host)))
	sByte = append(sByte, []byte(Host)...)
	fmt.Println(hex.EncodeToString(sByte))
	return sByte
}

func (self *TNet) genSocksV5Req() []byte {
	DevLogs.Debug("TNet.genSocksV5Req")
	var (
		sByte []byte
	)
	sByte = append(sByte, 0x5) // VER : 协议版本号，固定取值 0x05
	sByte = append(sByte, 0x2) // NMETHODS : 客户端支持的认证机制数目
	sByte = append(sByte, 0x0) // NO AUTHENTICATION REQUIRED
	//sByte = append(sByte, 0x1) // GSSAPI
	sByte = append(sByte, 0x2) // USERNAME/PASSWORD
	return sByte
}

func (self *TNet) genSocksV5Connect(Host string, Port string) []byte {
	DevLogs.Debug("TNet.genSocksV5Connect")
	var (
		sByte []byte
		//pUser, pPwd string = self.Proxy.GetProxyUserPwd()
	)
	sByte = append(sByte, 0x5) // VER : 协议版本号
	sByte = append(sByte, 0x1) // CMD : 请求类型  - 0x01 : CONNECT - 0x02 : BIND - 0x03 : UDP ASSOCIATE
	sByte = append(sByte, 0x0) // RSV : 保留字段，固定取值 0x00

	/*
		var dstIP   []byte = make([]byte, 4)
		sByte = append(sByte, 0x1) // ATYP : DST.ADDR地址类型 - 0x01 : IPV4 - 0x03 : 域名 - 0x04 : IPV6
		sByte = append(sByte, helpIPHostToByte(Host)...)
	*/
	sByte = append(sByte, 0x3) // ATYP : DST.ADDR地址类型 - 0x01 : IPV4 - 0x03 : 域名 - 0x04 : IPV6
	sByte = append(sByte, byte(len(Host)))
	sByte = append(sByte, []byte(Host)...)

	sByte = append(sByte, helpPortToByte(Port)...)

	fmt.Println(hex.EncodeToString(sByte))
	return sByte
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

func helpIPHostToByte(Host string) []byte {
	var (
		dstIP []byte = make([]byte, 4)
	)
	binary.BigEndian.PutUint32(dstIP, helpIPToUInt32(net.ParseIP(Host)))
	return dstIP
}

func helpPortToByte(Port string) []byte {
	var (
		dstPort []byte = make([]byte, 2)
		iPort   int
	)
	iPort, _ = strconv.Atoi(Port)
	binary.BigEndian.PutUint16(dstPort, uint16(iPort))
	return dstPort
}

func helpIPToUInt32(ipnr net.IP) uint32 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}
