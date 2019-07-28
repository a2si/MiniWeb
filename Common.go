package MiniWeb

import (
	"time"

	mwCookie "github.com/MiniWeb/Cookie"
	mwCore "github.com/MiniWeb/Core"
	DevLogs "github.com/MiniWeb/DevLogs"
	mwHeader "github.com/MiniWeb/Header"
	mwProxy "github.com/MiniWeb/Proxy"
	mwURL "github.com/MiniWeb/UrlExtend"
	mwUserAgent "github.com/MiniWeb/UserAgent"
	mwError "github.com/MiniWeb/mwError"
)

const (
	prv_MINI_WEB_VERSION    string        = "3.1.3" // MiniWeb 版本
	prv_MW_TIME_OUT         time.Duration = 30      // 传输超时
	prv_MW_TIME_OUT_CONNECT time.Duration = 50      // 连接超时
	prv_MINI_WEB_USERAGENT  string        = "MiniWeb V" + prv_MINI_WEB_VERSION
)

type miniWeb struct {
	prv_Core  *mwCore.WebCore
	prv_Byte  []byte
	prv_Error *mwError.TError
}

type tsConfig struct {
}

var (
	Config           *tsConfig
	cfgRandUserAgent bool = true
	cfgLogsEnable    bool = true
)

func init() {
	Config = &tsConfig{}
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
	if cfgRandUserAgent {
		Obj.prv_Core.UserAgent = mwUserAgent.NewUserAgent().Random()
	}
	Obj.initMiniWebClient()
	return Obj
}
