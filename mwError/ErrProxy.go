package mwError

import (
	"fmt"
	"strings"
)

func (self *TError) ProxyNoSettings() {
	self.prv_ErrCode = ERR_PROXY_NOT_SETTINGS
	self.prv_ErrMsg = MsgProxyNoSettings
}

func (self *TError) ProxyCheckError(MsgCode int, MsgInfo string) {
	switch MsgCode {
	case 503:
		// 503 Too many open connections
		if strings.Contains(MsgInfo, "any open connections") {
			self.prv_ErrCode = ERR_PROXY_MANY_CONNECTIONS
			self.prv_ErrMsg = MsgProxyManyConnections
			return
		}
	case 401:
	}
	fmt.Println("ProxyCheckError: ", MsgCode, MsgInfo)
}

func (self *TError) ProxyIDEntdError() {
	self.prv_ErrCode = ERR_PROXY_SOCKS_IDENTD
	self.prv_ErrMsg = MsgProxyErrID
}

func (self *TError) ProxyRefusedFail() {
	self.prv_ErrCode = ERR_PROXY_REFUSED_FAIL
	self.prv_ErrMsg = MsgProxyRefusedFail
}

func (self *TError) ProxyAuthenticationFail() {
	self.prv_ErrCode = ERR_PROXY_ACCOUNT_AUTH_FAIL
	self.prv_ErrMsg = MsgAccountAuthFail
}
