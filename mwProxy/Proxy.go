package Proxy

import (
	"encoding/base64"
)

func (self *TProxy) SetProxyType(pType int) {
	self.prv_ProxyMode = pType
}

func (self *TProxy) GetProxyType() int {
	return self.prv_ProxyMode
}

func (self *TProxy) Enable() bool {
	return self.prv_ProxyMode != PROXY_TYPE_NONE
}

func (self *TProxy) CancelProxy() {
	self.prv_ProxyMode = PROXY_TYPE_NONE
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

func (self *TProxy) GetBase64Authorization() string {
	if len(self.prv_UserName) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(self.prv_UserName + ":" + self.prv_Password))
}
