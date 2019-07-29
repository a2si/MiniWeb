package mwConfig

import (
	DevLogs "github.com/a2si/MiniWeb/DevLogs"
)

type TConfig struct {
}

var (
	cfgMap      map[string]bool
	bConfigInit bool = false
)

func ConfigInit() {
	if bConfigInit == true {
		return
	}
	bConfigInit = true

	cfgMap = make(map[string]bool)
	SetConfig("UserAgent.Rand", true)        // UA 随机
	SetConfig("Logs.Enable", true)           // 日志启用
	SetConfig("NetBase.QUIC", false)         // QUIC 是否支持
	SetConfig("Proxy.SocksV5.GSSAPI", false) // QUIC 是否支持

}

func bindConfigChnage(Name string, Enable bool) {
	switch Name {
	case "Logs.Enable":
		DevLogs.LogsEnable(Enable)
	}
}

func SetConfig(Name string, Enable bool) {
	cfgMap[Name] = Enable
	bindConfigChnage(Name, Enable)
}

func GetConfig(Name string) bool {
	return cfgMap[Name]
}
