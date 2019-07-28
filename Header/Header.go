package Header

import (
	"fmt"
	"strings"

	mwError "github.com/a2si/MiniWeb/mwError"
)

type Header struct {
	header   map[string]string // HEADER 信息
	ObjError *mwError.TError   // Error Object
}

func NewHeader(errObj *mwError.TError) *Header {
	Obj := &Header{
		header:   make(map[string]string),
		ObjError: errObj,
	}
	return Obj
}

func (self *Header) HeaderExists(Key string) bool {
	if _, ok := self.header[Key]; ok {
		return true
	}
	LowerKey := strings.ToLower(Key)
	for k, _ := range self.header {
		if strings.ToLower(k) == LowerKey {
			return true
		}
	}
	return false
}

func (self *Header) SetHeader(Name string, Value string) {
	LowerKey := strings.ToLower(Name)
	for k, _ := range self.header {
		if strings.ToLower(k) == LowerKey {
			self.header[k] = Value
			return
		}
	}
	self.header[Name] = Value
}

func (self *Header) GetHeader(Name string) string {
	LowerKey := strings.ToLower(Name)
	for k, v := range self.header {
		if strings.ToLower(k) == LowerKey {
			return v
		}
	}
	return ""
}

func (self *Header) RemoveHeader(Name string) {
	if self.HeaderExists(Name) {
		delete(self.header, Name)
	}
}

func (self *Header) ClearHeader() {
	self.header = make(map[string]string)
}

func (self *Header) GetAllHeader() string {
	var dwRet string
	for k, v := range self.header {
		dwRet = fmt.Sprintf("%s\r\n%s: %s", dwRet, k, v)
	}
	return dwRet
}
