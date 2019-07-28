package mwError

import (
	mwConst "github.com/MiniWeb/mwConst"
)

func (self *TError) SocketReadTimeout() {
	self.prv_ErrCode = mwConst.ERROR_SOCKET_READ_TIMEOUT
	self.prv_ErrMsg = "read: i/o timeout"
}

//read: connection reset by peer
func (self *TError) SocketReadReset() {
	self.prv_ErrCode = mwConst.ERROR_SOCKET_REMOTE_CLOSE
	self.prv_ErrMsg = "read: connection reset by peer"
}
