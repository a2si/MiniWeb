package mwError

import (
	"fmt"
	"strings"
)

func (self *TError) SocketNotSupport() {
	self.prv_ErrCode = ERR_NOW_NO_SUPPORT
	self.prv_ErrMsg = MsgNotSupport
}

func (self *TError) SocketStartError(err error) {
	str := err.Error()
	//dial tcp 117.90.1.61:9000: connect: connection refused
	if strings.Contains(str, "connect: connection refused") {
		self.SocketNetWorkConnectFail()
	}
	//dial tcp 219.159.38.202:56210: i/o timeout
	if strings.Contains(str, "i/o timeout") {
		self.SocketTimeoutError()
	}
	fmt.Println("SocketStartError: ", err)
}

func (self *TError) SocketIOError(err error) {
	str := err.Error()
	if strings.Contains(str, "i/o timeout") {
		self.SocketTimeoutError()
		return
	}
	if strings.Contains(str, "write: broken pipe") {
		self.SocketRemoteClose()
		return
	}
	fmt.Println("SocketIOError: ", err)
}

func (self *TError) SocketRemoteClose() {
	self.prv_ErrCode = ERR_NETWORK_REMOTE_CLOSE
	self.prv_ErrMsg = MsgRemoteClose
}

func (self *TError) SocketNetWorkConnectFail() {
	self.prv_ErrCode = ERR_NETWORK_CONNECT_FAIL
	self.prv_ErrMsg = MsgConnectFail
}

func (self *TError) SocketTimeoutError() {
	self.prv_ErrCode = ERR_NETWORK_TIMEOUT
	self.prv_ErrMsg = MsgTimeOut
}

func (self *TError) SocketNetWorkNoConnect() {
	self.prv_ErrCode = ERR_NETWORK_NOT_CONNECT
	self.prv_ErrMsg = MsgNetworkNotConnect
}

//write: broken pipe
