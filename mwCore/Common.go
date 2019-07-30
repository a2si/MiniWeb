package mwCore

import (
	"time"

	mwCookie "github.com/a2si/MiniWeb/mwCookie"
	mwError "github.com/a2si/MiniWeb/mwError"
	mwHeader "github.com/a2si/MiniWeb/mwHeader"
	mwNet "github.com/a2si/MiniWeb/mwNet"
	mwProxy "github.com/a2si/MiniWeb/mwProxy"
	mwURL "github.com/a2si/MiniWeb/mwURL"
)

type WebCore struct {
	Method         string            // HTTP
	Referer        string            // HTTP
	UserAgent      string            // HTTP
	Redirect       bool              // 重定向
	TimeOut        time.Duration     // 数据返回超时
	TimeOutConnect time.Duration     // 连接超时
	PostData       map[string]string // POST Data
	Result         []byte            // Web Response
	HttpVersion    string            // HTTP/1.1
	StatusCode     int               // Status Code
	StatusMsg      string            // Status Text
	URL            *mwURL.TUrl       // 访问地址
	Cookie         *mwCookie.Cookie  // Cookie
	ReqHeader      *mwHeader.Header  // 请求 Request Header
	RspHeader      *mwHeader.Header  // 响应 Response Header
	Proxy          *mwProxy.TProxy   // Proxy Module
	ObjError       *mwError.TError   // Error Object
	cbFunc         mwNet.CBHOOK
}
