package Proxy

import (
	mwConst "github.com/a2si/MiniWeb/mwConst"
)

func (self *TProxy) SetProxyType(pType int) {
	self.prv_ProxyMode = pType
}

func (self *TProxy) GetProxyType() int {
	return self.prv_ProxyMode
}

func (self *TProxy) Enable() bool {
	return self.prv_ProxyMode != mwConst.PROXY_TYPE_NONE
}

func (self *TProxy) CancelProxy() {
	self.prv_ProxyMode = mwConst.PROXY_TYPE_NONE
}

func (self *TProxy) SetRealAddr(IP string, Port string) string {
	//URLAddr 转换到 IPAddr
	self.prv_Real_IP = IP
	self.prv_Real_Port = Port

	switch self.prv_ProxyMode {
	case mwConst.PROXY_TYPE_HTTP:
		return self.prv_IP + ":" + self.prv_Port
	case mwConst.PROXY_TYPE_HTTPS:
		return self.prv_IP + ":" + self.prv_Port
	case mwConst.PROXY_TYPE_SOCKS4:
	case mwConst.PROXY_TYPE_SOCKS4A:
	case mwConst.PROXY_TYPE_SOCKS5:
	}
	return self.prv_Real_IP + ":" + self.prv_Real_Port
}

func (self *TProxy) SetProxyIP(IP string, Port string) {
	self.prv_IP = IP
	self.prv_Port = Port
}

func (self *TProxy) GetProxyIP() (string, string) {
	return self.prv_IP, self.prv_Port
}

func (self *TProxy) SetProxyUserPwd(User string, Password string) {
	self.prv_UserName = User
	self.prv_Password = Password
}

func (self *TProxy) GetProxyUserPwd() (string, string) {
	return self.prv_UserName, self.prv_Password
}
