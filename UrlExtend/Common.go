package UrlExtend

import (
	"net/url"

	mwError "github.com/a2si/MiniWeb/mwError"
)

type TUrl struct {
	URL      *url.URL        // 访问地址
	ObjError *mwError.TError // Error Object
}

func NewUrl(errObj *mwError.TError) *TUrl {
	Obj := &TUrl{
		URL:      &url.URL{},
		ObjError: errObj,
	}
	return Obj
}
