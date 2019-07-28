package Network

import (
	"encoding/base64"
	"fmt"
	"strings"

	mwCommon "github.com/MiniWeb/Common"
	DevLogs "github.com/MiniWeb/DevLogs"
)

func (self *TNet) genHttpProxyReqHeader(Host string, Port string) string {
	var (
		Header      string = fmt.Sprintf("CONNECT %s:%s HTTP/1.1\r\n", Host, Port)
		pUser, pPwd string = self.Proxy.GetProxyUserPwd()
	)

	if len(pUser) != 0 {
		str := base64.StdEncoding.EncodeToString([]byte(pUser + ":" + pPwd))
		Header = fmt.Sprintf("%sProxy-Authorization:Basic %s\r\n", Header, str)
	}
	//Header = fmt.Sprintf("%sProxy-Connection: Keep-Alive\r\n", Header)
	Header = fmt.Sprintf("%s\r\n", Header)
	return Header
}

func (self *TNet) httpProxyRspHeader() (string, string) {
	DevLogs.Debug("TNet.httpProxyRspHeader")
	var (
		//ProtoVer string
		MsgCode string
		MsgInfo string
	)
	for {
		Text := self.ReadLine()
		fmt.Println("httpProxyRspHeader: ", Text)
		if strings.Contains(Text, "HTTP/") {
			_, MsgCode, MsgInfo = mwCommon.ReadHeaderVCM(Text)
			if MsgCode == "200" {
				return MsgCode, MsgInfo
			}
		}
		if Text == "\r\n" || Text == "" {
			return MsgCode, MsgInfo
		}
	}
	return MsgCode, MsgInfo
}
