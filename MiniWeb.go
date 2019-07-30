package MiniWeb

import (
	"net/url"
	"time"

	mwCookie "github.com/a2si/MiniWeb/Cookie"
	mwCore "github.com/a2si/MiniWeb/Core"
	DevLogs "github.com/a2si/MiniWeb/DevLogs"
	mwHeader "github.com/a2si/MiniWeb/Header"
	mwProxy "github.com/a2si/MiniWeb/Proxy"
	mwURL "github.com/a2si/MiniWeb/UrlExtend"
	mwUserAgent "github.com/a2si/MiniWeb/UserAgent"
	mwConfig "github.com/a2si/MiniWeb/mwConfig"
	mwError "github.com/a2si/MiniWeb/mwError"
)

type (
	TMiniWeb struct {
		prv_Core  *mwCore.WebCore
		prv_Error *mwError.TError
		prv_Byte  []byte
	}
	iMiniWeb interface {
		Cookie() *mwCookie.Cookie                       // Cookie 模块
		ReqHeader() *mwHeader.Header                    // 请求HTTP头部
		RspHeader() *mwHeader.Header                    // 响应HTTP头部
		Proxy() *mwProxy.TProxy                         // 代理模块
		IsClient()                                      // 正常模式
		IsXMLHttp()                                     // 模拟XML UA
		IsWeiXin()                                      // 模拟微UA
		SetURL(URL string)                              // 设置访问地址
		SetHttpMethod(Method string)                    // 设置请求协议
		SetReferer(Referer string)                      // 设置请求来源
		GetReferer() string                             // 取得设定访问来源
		SetUserAgent(UserAgent string)                  // 设置UA
		GetUserAgent() string                           // 获取当前UA
		SetRedirect(Redirect bool)                      // 设置是否自动转向 如302
		GetRedirect() bool                              // 读取转向设置
		SetTimeOut(TimeOut time.Duration)               // 设置R/W超时
		GetTimeOut() time.Duration                      // 取得R/W超时设置
		SetTimeOutConnect(TimeOutConnect time.Duration) // 设置连接超时
		GetTimeOutConnect() time.Duration               // 获取连接超时设定
		SetPOST(data map[string]string)                 // POST 上传数据
		SetPOSTFile(Name string, FileName string)       // POST 上传文件
		SetErrorMsg(Msg string)                         // 设置错误信息
		GetErrorMsg() string                            // 读取错误信息
		SetErrorCode(Code int)                          // 设置错误代码
		GetErrorCode() int                              // 读取错误代码
		GetStatusCode() int                             // 读取状态码
		GetWebCode(URL string) int                      // 访问网页
		SendRequest(URL string) int                     // 访问网页 同上
		ResponseText() string                           // 返回页面字符串数据
		ResponseByte() []byte                           // 返回页面字节流数据
	}
)

const (
	prv_MINI_WEB_VERSION        string        = "3.1.3" // MiniWeb 版本
	prv_MW_TIME_OUT             time.Duration = 30      // 传输超时
	prv_MW_TIME_OUT_CONNECT     time.Duration = 50      // 连接超时
	prv_MINI_WEB_USERAGENT      string        = "MiniWeb V" + prv_MINI_WEB_VERSION
	prv_MW_COOKIE_DIR           string        = "Cookie" // Cookie 存储路径
	prv_MW_COOKIE_DIR_AUTO_MAKE bool          = false    // Cookie 存储路径自动创建
	prv_MW_COOKIE_SAVE          bool          = false    // Cookie 自动保存
)

func init() {
	mwConfig.ConfigInit()
}

func NewMiniWeb() *TMiniWeb {
	DevLogs.Debug("Package.NewMiniWeb")
	self := &TMiniWeb{
		prv_Error: mwError.NewError(),
	}
	self.prv_Core = &mwCore.WebCore{
		ObjError:       self.prv_Error,
		Method:         "GET",
		Referer:        "",
		UserAgent:      prv_MINI_WEB_USERAGENT,
		Redirect:       false,
		TimeOut:        prv_MW_TIME_OUT,
		TimeOutConnect: prv_MW_TIME_OUT_CONNECT,
		PostData:       make(map[string]string),
		URL:            mwURL.NewUrl(self.prv_Error),
		Cookie:         mwCookie.NewCookie(self.prv_Error),
		ReqHeader:      mwHeader.NewHeader(self.prv_Error),
		RspHeader:      mwHeader.NewHeader(self.prv_Error),
		Proxy:          mwProxy.NewProxy(self.prv_Error),
	}

	self.prv_Core.Cookie.SetSaveCookie(prv_MW_COOKIE_SAVE)
	self.prv_Core.Cookie.SetCookieDir(prv_MW_COOKIE_DIR, prv_MW_COOKIE_DIR_AUTO_MAKE)
	if mwConfig.GetConfig("UserAgent.Rand") {
		self.prv_Core.UserAgent = mwUserAgent.NewUserAgent().Random()
	}
	self.initMiniWebClient()
	return self
}

func (self *TMiniWeb) initMiniWebClient() {
	DevLogs.Debug("MiniWeb.initMiniWebClient")
	self.prv_Core.Cookie.Clear()
	self.prv_Core.ReqHeader.ClearHeader()
	self.prv_Core.RspHeader.ClearHeader()
	self.prv_Core.InitHeader()

}

// 暴露接口
func (self *TMiniWeb) Cookie() *mwCookie.Cookie {
	return self.prv_Core.Cookie
}

// 暴露接口
func (self *TMiniWeb) ReqHeader() *mwHeader.Header {
	return self.prv_Core.ReqHeader
}

// 暴露接口
func (self *TMiniWeb) RspHeader() *mwHeader.Header {
	return self.prv_Core.RspHeader
}

// 暴露接口
func (self *TMiniWeb) Proxy() *mwProxy.TProxy {
	return self.prv_Core.Proxy
}

// 正常模式
func (self *TMiniWeb) IsClient() {
	DevLogs.Debug("MiniWeb.IsClient")
	self.prv_Core.ReqHeader.RemoveHeader("X-Requested-With")
	self.prv_Core.ReqHeader.SetHeader("User-Agent", self.prv_Core.UserAgent)
}

// XML HTTP 模式
func (self *TMiniWeb) IsXMLHttp() {
	DevLogs.Debug("MiniWeb.IsXMLHttp")
	self.prv_Core.ReqHeader.SetHeader("X-Requested-With", "XMLHttpRequest")
}

// 微信客户端模块
func (self *TMiniWeb) IsWeiXin() {
	DevLogs.Debug("MiniWeb.IsWeiXin")
	self.prv_Core.ReqHeader.RemoveHeader("X-Requested-With")
	self.prv_Core.ReqHeader.SetHeader("User-Agent", "MicroMessenger")
}

func (self *TMiniWeb) SetURL(URL string) {
	DevLogs.Debug("MiniWeb.SetURL")
	if len(URL) > 0 {
		self.prv_Core.URL.SetUrl(URL)
		self.Cookie().SetURL(self.prv_Core.URL.GetHost())
	}
}

func (self *TMiniWeb) SetHttpMethod(Method string) {
	DevLogs.Debug("MiniWeb.SetHttpMethod")
	self.prv_Core.SetMethod(Method)
}

func (self *TMiniWeb) SetReferer(Referer string) {
	DevLogs.Debug("MiniWeb.SetReferer Referer=" + Referer)
	self.prv_Core.Referer = Referer
	self.prv_Core.ReqHeader.SetHeader("Referer", Referer)
}

func (self *TMiniWeb) GetReferer() string {
	DevLogs.Debug("MiniWeb.GetReferer")
	return self.prv_Core.Referer
}

func (self *TMiniWeb) SetUserAgent(UserAgent string) {
	DevLogs.Debug("MiniWeb.SetUserAgent: UserAgent=" + UserAgent)
	self.prv_Core.UserAgent = UserAgent
	self.prv_Core.ReqHeader.SetHeader("User-Agent", UserAgent)

}

func (self *TMiniWeb) GetUserAgent() string {
	DevLogs.Debug("MiniWeb.GetUserAgent")
	return self.prv_Core.UserAgent
}

func (self *TMiniWeb) SetRedirect(Redirect bool) {
	DevLogs.Debug("MiniWeb.SetRedirect")
	self.prv_Core.Redirect = Redirect
}

func (self *TMiniWeb) GetRedirect() bool {
	DevLogs.Debug("MiniWeb.GetRedirect")
	return self.prv_Core.Redirect
}

func (self *TMiniWeb) SetTimeOut(TimeOut time.Duration) {
	DevLogs.Debug("MiniWeb.SetTimeOut")
	self.prv_Core.TimeOut = TimeOut
}

func (self *TMiniWeb) GetTimeOut() time.Duration {
	DevLogs.Debug("MiniWeb.GetTimeOut")
	return self.prv_Core.TimeOut
}

func (self *TMiniWeb) SetTimeOutConnect(TimeOutConnect time.Duration) {
	DevLogs.Debug("MiniWeb.SetTimeOutConnect")
	self.prv_Core.TimeOutConnect = TimeOutConnect
}

func (self *TMiniWeb) GetTimeOutConnect() time.Duration {
	DevLogs.Debug("MiniWeb.GetTimeOutConnect")
	return self.prv_Core.TimeOutConnect
}

// MAP 方法, 数据顺序会改变
func (self *TMiniWeb) SetPOST(data map[string]string) {
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

func (self *TMiniWeb) SetPOSTFile(Name string, FileName string) {
	DevLogs.Debug("MiniWeb.SetPOSTFile")
	self.SetHttpMethod("POST")
	self.prv_Core.AddPost("@"+Name, FileName)
}

func (self *TMiniWeb) SetErrorMsg(Msg string) {
	self.prv_Error.SetErrorMsg(Msg)
}

func (self *TMiniWeb) GetErrorMsg() string {
	return self.prv_Error.GetErrorMsg()
}

func (self *TMiniWeb) SetErrorCode(Code int) {
	self.prv_Error.SetErrorCode(Code)
}

func (self *TMiniWeb) GetErrorCode() int {
	return self.prv_Error.GetErrorCode()
}

func (self *TMiniWeb) GetStatusCode() int {
	return self.prv_Core.StatusCode
}

func (self *TMiniWeb) GetWebCode(URL string) int {
	DevLogs.Debug("MiniWeb.GetWebCode")
	return self.SendRequest(URL)
}

func (self *TMiniWeb) SendRequest(URL string) int {
	DevLogs.Debug("MiniWeb.SendRequest")

	self.SetURL(URL)
	dwCode := self.prv_Core.SendRequest()
	self.SetHttpMethod("GET")
	return dwCode
}

func (self *TMiniWeb) ResponseText() string {
	DevLogs.Debug("MiniWeb.ResponseText")
	return string(self.prv_Core.Result)
}

func (self *TMiniWeb) ResponseByte() []byte {
	DevLogs.Debug("MiniWeb.ResponseByte")
	return self.prv_Core.Result

}
