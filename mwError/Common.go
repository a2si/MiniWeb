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
	}
	return Obj
}
