package MiniWeb

import (
	mwCookie "MiniWeb/Cookie"
	mwCore "MiniWeb/Core"
	DevLogs "MiniWeb/DevLogs"
	mwHeader "MiniWeb/Header"
	mwURL "MiniWeb/UrlExtend"
	mwUserAgent "MiniWeb/UserAgent"
	"time"
)

const (
	prv_MINI_WEB_VERSION    string        = "3.1.1" // MiniWeb 版本
	prv_MW_TIME_OUT         time.Duration = 30      // 传输超时
	prv_MW_TIME_OUT_CONNECT time.Duration = 50      // 连接超时
	prv_MINI_WEB_USERAGENT  string        = "MiniWeb V" + prv_MINI_WEB_VERSION
)

type miniWeb struct {
	prv_Core *mwCore.WebCore
	prv_Byte []byte
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
		prv_Core: &mwCore.WebCore{
			Method:         "GET",
			Referer:        "",
			UserAgent:      prv_MINI_WEB_USERAGENT,
			Redirect:       false,
			TimeOut:        prv_MW_TIME_OUT,
			TimeOutConnect: prv_MW_TIME_OUT_CONNECT,
			PostData:       make(map[string]string),
			URL:            mwURL.NewUrl(),
			Cookie:         mwCookie.NewCookie(),
			ReqHeader:      mwHeader.NewHeader(),
			RspHeader:      mwHeader.NewHeader(),
		},
	}
	if cfgRandUserAgent {
		Obj.prv_Core.UserAgent = mwUserAgent.NewUserAgent().Random()
	}
	Obj.initMiniWebClient()
	return Obj
}
