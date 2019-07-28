package mwError

import (
	mwConst "github.com/MiniWeb/mwConst"
)

func (self *TError) SocketRemoteClose() {
	self.prv_ErrCode = mwConst.ERROR_SOCKET_REMOTE_CLOSE
	self.prv_ErrMsg = "connect: connection refused"
}

func (self *TError) SocketConnectTimeout() {
	self.prv_ErrCode = mwConst.ERROR_SOCKET_CONNECT_TIMEOUT
	self.prv_ErrMsg = "connect: timeout"
}

//write: broken pipe
