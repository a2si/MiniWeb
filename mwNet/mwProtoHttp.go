package mwNet

import (
	"fmt"
	"strings"

	DevLogs "github.com/a2si/MiniWeb/DevLogs"
	mwCommon "github.com/a2si/MiniWeb/mwCommon"
)

func (self *TNet) genHttpProxyReqHeader(Host string, Port string) string {
	var (
		Header     string = fmt.Sprintf("CONNECT %s:%s HTTP/1.1\r\n", Host, Port)
		AuthBase64 string = self.Proxy.GetBase64Authorization()
	)
	if len(AuthBase64) != 0 {
		Header = fmt.Sprintf("%sProxy-Authorization:Basic %s\r\n", Header, AuthBase64)
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
		if self.ObjError.IsError() {
			return MsgCode, MsgInfo
		}
	}
}
