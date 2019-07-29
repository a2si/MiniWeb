package UrlExtend

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	URL := NewUrl(nil)
	URL.SetUrl("http://www.baidu.com:33/?asdf")
	fmt.Println(URL.GetHost())
}
