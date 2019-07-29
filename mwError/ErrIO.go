package mwError

import (
	"fmt"
	"strings"
)

func (self *TError) IOReadByNegative() {
	self.prv_ErrCode = ERR_IO_READ_BY_NEGATIVE
	self.prv_ErrMsg = MsgIOReadByNegative
}

func (self *TError) SetIOError(err error) {
	str := err.Error()
	fmt.Println(err)
	if strings.Contains(str, "read: connection refused") {
		return
	}
}
