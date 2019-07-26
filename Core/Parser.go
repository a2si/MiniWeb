package Core

import (
	"strconv"
	"strings"
)

func (self *WebCore) rspParserHeaderLine(l string) {
	if len(l) > 0 {
		if strings.Contains(l, ":") == true {
			Index := strings.Index(l, ":")
			sByte := []byte(l)
			Name := string(sByte[:Index])
			Index += 1
			Value := string(sByte[Index:])
			if strings.ToLower(Name) == "set-cookie" {
				self.Cookie.ParserCookie(Value)
			} else {
				self.RspHeader.SetHeader(Name, strings.Trim(Value, " "))
			}
		} else {
			arr := strings.Split(l, " ")
			self.HttpVersion = arr[0]
			self.StatusCode, _ = strconv.Atoi(arr[1])
			if len(arr) >= 3 {
				self.StatusMsg = arr[2]
			}
		}
	}
}
