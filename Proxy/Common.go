package Proxy

import (
	mwConst "github.com/MiniWeb/mwConst"
	mwError "github.com/MiniWeb/mwError"
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
		prv_ProxyMode: mwConst.PROXY_TYPE_NONE,
		ObjError:      errObj,
	}
	return Obj
}
