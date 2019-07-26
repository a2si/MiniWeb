package Core

const (
	ERROR_CODE_NO_ERROR = iota
	ERROR_CODE_TOC      = iota
	ERROR_CODE_READ_TIME
	ERROR_CODE_CONN_REFUSED
	ERROR_CODE_CONN_NO_ROUTE
)
const (
	ERROR_MSG_NO_ERROR      = ""
	ERROR_MSG_TOC           = "Connect Time Out"
	ERROR_MSG_READ_TIME     = "Read Time Out"
	ERROR_MSG_CONN_REFUSED  = "connection refused"
	ERROR_MSG_CONN_NO_ROUTE = "no route to host"
)

var (
	ErrCode int
	ErrMsg  string
)

func setError(eCode int, Msg string) {
	ErrCode = eCode
	ErrMsg = Msg
}

func setErrorCode(eCode int) {
	ErrCode = eCode
}

func setErrorMsg(Msg string) {
	ErrMsg = Msg
}

func getErrorCode() int {
	return ErrCode
}

func getErrorMsg() string {
	return ErrMsg
}
