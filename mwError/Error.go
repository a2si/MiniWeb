package mwError

type TError struct {
	prv_ErrCode int
	prv_ErrMsg  string
	prv_ErrObj  error
}

func NewError() *TError {
	Obj := &TError{
		prv_ErrCode: 0,
		prv_ErrMsg:  "",
		prv_ErrObj:  nil,
	}
	return Obj
}

func (self *TError) IsError() bool {
	return self.prv_ErrCode != ERROR_NO_ERROR || self.prv_ErrObj != nil
}

func (self *TError) Clear() {
	self.prv_ErrMsg = ""
	self.prv_ErrCode = 0
	self.prv_ErrObj = nil
}

func (self *TError) SetErrorMsg(Msg string) {
	self.prv_ErrMsg = Msg
}

func (self *TError) GetErrorMsg() string {
	defer func() {
		self.prv_ErrMsg = ""
	}()
	return self.prv_ErrMsg
}

func (self *TError) SetErrorCode(Code int) {
	self.prv_ErrCode = Code
}

func (self *TError) GetErrorCode() int {
	defer func() {
		self.prv_ErrCode = 0
	}()
	return self.prv_ErrCode
}

func (self *TError) SetErrorObj(Obj error) {
	self.prv_ErrObj = Obj
}

func (self *TError) GetErrorObj() error {
	defer func() {
		self.prv_ErrObj = nil
	}()
	return self.prv_ErrObj
}
