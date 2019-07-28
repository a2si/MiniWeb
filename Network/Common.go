package Network

import (
	"bufio"
	"net"
	"time"

	mwProxy "github.com/a2si/MiniWeb/Proxy"
	mwError "github.com/a2si/MiniWeb/mwError"
)

type TNet struct {
	ObjError *mwError.TError // Error Object
	Proxy    *mwProxy.TProxy // Proxy Module
	Conn     net.Conn
	ioRead   *bufio.Reader
	timeOut  time.Duration
	isClosed bool
}

func NewNet(errObj *mwError.TError, p *mwProxy.TProxy) *TNet {
	Obj := &TNet{
		ObjError: errObj,
		Proxy:    p,
		timeOut:  0,
		isClosed: false,
	}
	return Obj
}

func Host2IP(Host string) string {
	str, _ := net.LookupHost(Host)
	if len(str) > 1 {
		return str[0]
	}
	//fmt.Println(err)
	return ""
}
