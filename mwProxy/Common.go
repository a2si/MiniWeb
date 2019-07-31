package Proxy

import (
	mwError "github.com/a2si/MiniWeb/mwError"
)

const (
	PROXY_TYPE_NONE = iota
	PROXY_TYPE_HTTP
	PROXY_TYPE_HTTPS
	PROXY_TYPE_SOCKS4
	PROXY_TYPE_SOCKS4A
	PROXY_TYPE_SOCKS5
)

type TProxy struct {
	prv_ProxyMode int
	prv_Real_IP   string
	prv_Real_Port string
	prv_IP        string
	prv_Port      string
	prv_UserName  string
	prv_Password  string
	ObjError      *mwError.TError // Error Object
}

func NewProxy(errObj *mwError.TError) *TProxy {
	Obj := &TProxy{
		prv_ProxyMode: PROXY_TYPE_NONE,
		ObjError:      errObj,
	}
	return Obj
}
