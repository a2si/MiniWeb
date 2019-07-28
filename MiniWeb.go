package MiniWeb

import (
	"net/url"
	"time"

	mwCookie "github.com/MiniWeb/Cookie"
	DevLogs "github.com/MiniWeb/DevLogs"
	mwHeader "github.com/MiniWeb/Header"
	mwProxy "github.com/MiniWeb/Proxy"
)

func (self *miniWeb) initMiniWebClient() {
	DevLogs.Debug("MiniWeb.initMiniWebClient")
	self.prv_Core.Cookie.Clear()
	self.prv_Core.ReqHeader.ClearHeader()
	self.prv_Core.RspHeader.ClearHeader()
	self.prv_Core.InitHeader()

}

// 暴露接口
func (self *miniWeb) Cookie() *mwCookie.Cookie {
	return self.prv_Core.Cookie
}

// 暴露接口
func (self *miniWeb) ReqHeader() *mwHeader.Header {
	return self.prv_Core.ReqHeader
}

// 暴露接口
func (self *miniWeb) RspHeader() *mwHeader.Header {
	return self.prv_Core.RspHeader
}

// 暴露接口
func (self *miniWeb) Proxy() *mwProxy.TProxy {
	return self.prv_Core.Proxy
}

// 正常模式
func (self *miniWeb) IsClient() {
	DevLogs.Debug("MiniWeb.IsClient")
	self.prv_Core.ReqHeader.RemoveHeader("X-Requested-With")
	self.prv_Core.ReqHeader.SetHeader("User-Agent", self.prv_Core.UserAgent)
}

// XML HTTP 模式
func (self *miniWeb) IsXMLHttp() {
	DevLogs.Debug("MiniWeb.IsXMLHttp")
	self.prv_Core.ReqHeader.SetHeader("X-Requested-With", "XMLHttpRequest")
}

// 微信客户端模块
func (self *miniWeb) IsWeiXin() {
	DevLogs.Debug("MiniWeb.IsWeiXin")
	self.prv_Core.ReqHeader.RemoveHeader("X-Requested-With")
	self.prv_Core.ReqHeader.SetHeader("User-Agent", "MicroMessenger")
}

func (self *miniWeb) SetURL(URL string) {
	DevLogs.Debug("MiniWeb.SetURL")
	if len(URL) > 0 {
		self.prv_Core.URL.SetUrl(URL)
	}
}

func (self *miniWeb) SetHttpMethod(Method string) {
	DevLogs.Debug("MiniWeb.SetHttpMethod")
	self.prv_Core.SetMethod(Method)
}

func (self *miniWeb) SetReferer(Referer string) {
	DevLogs.Debug("MiniWeb.SetReferer Referer=" + Referer)
	self.prv_Core.Referer = Referer
}

func (self *miniWeb) GetReferer() string {
	DevLogs.Debug("MiniWeb.GetReferer")
	return self.prv_Core.Referer
}

func (self *miniWeb) SetUserAgent(UserAgent string) {
	DevLogs.Debug("MiniWeb.SetUserAgent: UserAgent=" + UserAgent)
	self.prv_Core.UserAgent = UserAgent
}

func (self *miniWeb) GetUserAgent() string {
	DevLogs.Debug("MiniWeb.GetUserAgent")
	return self.prv_Core.UserAgent
}

func (self *miniWeb) SetRedirect(Redirect bool) {
	DevLogs.Debug("MiniWeb.SetRedirect")
	self.prv_Core.Redirect = Redirect
}

func (self *miniWeb) GetRedirect() bool {
	DevLogs.Debug("MiniWeb.GetRedirect")
	return self.prv_Core.Redirect
}

func (self *miniWeb) SetTimeOut(TimeOut time.Duration) {
	DevLogs.Debug("MiniWeb.SetTimeOut")
	self.prv_Core.TimeOut = TimeOut
}

func (self *miniWeb) GetTimeOut() time.Duration {
	DevLogs.Debug("MiniWeb.GetTimeOut")
	return self.prv_Core.TimeOut
}

func (self *miniWeb) SetTimeOutConnect(TimeOutConnect time.Duration) {
	DevLogs.Debug("MiniWeb.SetTimeOutConnect")
	self.prv_Core.TimeOutConnect = TimeOutConnect
}

func (self *miniWeb) GetTimeOutConnect() time.Duration {
	DevLogs.Debug("MiniWeb.GetTimeOutConnect")
	return self.prv_Core.TimeOutConnect
}

// MAP 方法, 数据顺序会改变
func (self *miniWeb) SetPOST(data map[string]string) {
	DevLogs.Debug("MiniWeb.SetPOST")
	self.SetHttpMethod("POST")
	for k, v := range data {
		Value, err := url.QueryUnescape(v)
		if err == nil {
			self.prv_Core.AddPost(k, Value)
		} else {
			DevLogs.Warn("QueryUnescape.Error: " + err.Error())
		}
	}
}

func (self *miniWeb) SetPOSTFile(Name string, FileName string) {
	DevLogs.Debug("MiniWeb.SetPOSTFile")
	self.SetHttpMethod("POST")
	self.prv_Core.AddPost("@"+Name, FileName)
}

func (self *miniWeb) SetErrorMsg(Msg string) {
	self.prv_Error.SetErrorMsg(Msg)
}

func (self *miniWeb) GetErrorMsg() string {
	return self.prv_Error.GetErrorMsg()
}

func (self *miniWeb) SetErrorCode(Code int) {
	self.prv_Error.SetErrorCode(Code)
}

func (self *miniWeb) GetErrorCode() int {
	return self.prv_Error.GetErrorCode()
}

func (self *miniWeb) GetStatusCode() int {
	return self.prv_Core.StatusCode
}

func (self *miniWeb) GetWebCode(URL string) int {
	DevLogs.Debug("MiniWeb.GetWebCode")

	if len(URL) > 0 {
		self.prv_Core.URL.SetUrl(URL)
	}
	dwCode := self.prv_Core.SendRequest()
	self.SetHttpMethod("GET")
	return dwCode
}

func (self *miniWeb) ResponseText() string {
	DevLogs.Debug("MiniWeb.ResponseText")
	return string(self.prv_Core.Result)
}

func (self *miniWeb) ResponseByte() []byte {
	DevLogs.Debug("MiniWeb.ResponseByte")
	return self.prv_Core.Result

}
