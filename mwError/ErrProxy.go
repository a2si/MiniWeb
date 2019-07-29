package mwError

func (self *TError) ProxyError503() {
	self.prv_ErrCode = ERROR_PROXY_503
	self.prv_ErrMsg = "Too many open connections"
}

func (self *TError) ProxyNoSupport() {
	self.prv_ErrCode = ERROR_PROXY_NO_SUPPORT
	self.prv_ErrMsg = "no support client do it"
}

func (self *TError) ProxyAuthenticationFail() {
	self.prv_ErrCode = ERROR_PROXY_AUTH_FAIL
	self.prv_ErrMsg = "authentication fail"
}

func (self *TError) ProxyOtherError() {
	self.prv_ErrCode = ERROR_PROXY_OTHER_ERROR
	self.prv_ErrMsg = "other error"
}
