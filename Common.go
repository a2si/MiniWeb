package MiniWeb

import (
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

const (
	prv_MINI_WEB_VERSION    string        = "3.1.3" // MiniWeb 版本
	prv_MW_TIME_OUT         time.Duration = 30      // 传输超时
	prv_MW_TIME_OUT_CONNECT time.Duration = 50      // 连接超时
	prv_MINI_WEB_USERAGENT  string        = "MiniWeb V" + prv_MINI_WEB_VERSION
)

type miniWeb struct {
	prv_Core  *mwCore.WebCore
	prv_Error *mwError.TError
	prv_Byte  []byte
}

var (
	Config *mwConfig.TConfig
)

func init() {
	Config = mwConfig.NewConfig()
}

func NewMiniWeb() *miniWeb {
	DevLogs.Debug("Package.NewMiniWeb")
	Obj := &miniWeb{
		prv_Error: mwError.NewError(),
	}
	Obj.prv_Core = &mwCore.WebCore{
		ObjError:       Obj.prv_Error,
		Method:         "GET",
		Referer:        "",
		UserAgent:      prv_MINI_WEB_USERAGENT,
		Redirect:       false,
		TimeOut:        prv_MW_TIME_OUT,
		TimeOutConnect: prv_MW_TIME_OUT_CONNECT,
		PostData:       make(map[string]string),
		URL:            mwURL.NewUrl(Obj.prv_Error),
		Cookie:         mwCookie.NewCookie(Obj.prv_Error),
		ReqHeader:      mwHeader.NewHeader(Obj.prv_Error),
		RspHeader:      mwHeader.NewHeader(Obj.prv_Error),
		Proxy:          mwProxy.NewProxy(Obj.prv_Error),
	}
	if Config.RandUserAgent {
		Obj.prv_Core.UserAgent = mwUserAgent.NewUserAgent().Random()
	}
	Obj.initMiniWebClient()
	return Obj
}
