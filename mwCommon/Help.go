package Common

import (
	"strings"
)

func ReadHeaderVCM(Text string) (string, string, string) {
	var (
		ProtoVer string
		MsgCode  string
		MsgInfo  string
	)
	readOne := func(Code string) (string, string) {
		Index := strings.Index(Code, " ")
		if Index == -1 {
			return "", Code
		} else {
			sByte := []byte(Code)
			lStr := string(sByte[:Index])
			rStr := string(sByte[Index+1:])
			return lStr, rStr
		}
	}
	ProtoVer, Text = readOne(Text)
	MsgCode, Text = readOne(Text)
	MsgInfo = Text
	return ProtoVer, MsgCode, MsgInfo
}
