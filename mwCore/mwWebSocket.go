package mwCore

import (
	"encoding/base64"

	mwNet "github.com/a2si/MiniWeb/mwNet"
)

func (self *WebCore) wsInitHeader() {
	self.ReqHeader.SetHeader("Connection", "Upgrade")
	self.ReqHeader.SetHeader("Upgrade", "WebSocket")
	//随机字符串，用于验证协议是否为WebSocket协议而非HTTP协议
	self.ReqHeader.SetHeader("Sec-WebSocket-key", base64.StdEncoding.EncodeToString([]byte(self.URL.GetHost())))
	//self.ReqHeader.SetHeader("Sec-WebSocket-key", "1AauaZL5W473vgMmaIWb4w==")
	if self.URL.IsTls() {
		self.ReqHeader.SetHeader("Origin", "http://"+self.URL.GetHost())
	} else {
		self.ReqHeader.SetHeader("Origin", "http://"+self.URL.GetHost())
	}

	//表示使用WebSocket的哪一个版本
	self.ReqHeader.SetHeader("Sec-WebSocket-Version", "13")
	//根据Sec-WebSocket-Accept和特殊字符串计算。验证协议是否为WebSocket协议
	//self.ReqHeader.SetHeader("Sec-WebSocket-Accept", "")
	//与Host字段对应，表示请求WebSocket协议的地址
	//self.ReqHeader.SetHeader("Sec-WebSocket-Location", "")

}

func (self *WebCore) connectWebSocket(NetWork *mwNet.TNet) int {
	self.wsInitHeader()
	Header := self.buildReqHeader()
	NetWork.SendPacket([]byte(Header))
	if self.ObjError.IsError() {
		return self.ObjError.GetErrorCode()
	}
	self.readRspHeader(NetWork)
	if self.ObjError.IsError() {
		return self.ObjError.GetErrorCode()
	}
	if self.StatusCode == 101 || self.StatusCode == 0 {
		//var wsBuffer []byte
		self.cbFunc(mwNet.EVENT_OBJECT, NetWork, nil)

		NetWork.WebSocketReadToCBHook(-1, func(Event int, ObjNet *mwNet.TNet, data []byte) {
			switch Event {

			case mwNet.EVENT_RECV_TEXT:
				NetWork.WebSocketSendPing([]byte("Helo Server"))
			case mwNet.EVENT_CLOSE:
				NetWork.WebSocketSendClose()
			case mwNet.EVENT_PING:
				NetWork.WebSocketSendPong(data)
			case mwNet.EVENT_PONG:
			}
			self.cbFunc(Event, ObjNet, data)
		})
	}
	return self.StatusCode
}
