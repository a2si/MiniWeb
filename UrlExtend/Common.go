package UrlExtend

import (
	"net/url"
)

type TUrl struct {
	URL *url.URL // 访问地址
}

func NewUrl() *TUrl {
	Obj := &TUrl{
		URL: &url.URL{},
	}
	return Obj
}
