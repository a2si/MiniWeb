package mwError

import mwConst "github.com/MiniWeb/mwConst"

func (self *TError) IsError() bool {
	return self.prv_ErrCode != mwConst.ERROR_NO_ERROR
}

func (self *TError) SetErrorMsg(Msg string) {
	self.prv_ErrMsg = Msg
}

func (self *TError) GetErrorMsg() string {
	return self.prv_ErrMsg
}

func (self *TError) SetErrorCode(Code int) {
	self.prv_ErrCode = Code
}

func (self *TError) GetErrorCode() int {
	return self.prv_ErrCode
}
