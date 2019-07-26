package Cookie

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	mwCommon "github.com/MiniWeb/Common"
)

type Cookie struct {
	BaseURI     string            // URL地址
	bSaveCookie bool              // 是否保存COOKIE
	CookieDir   string            // COOKIE 保存路径
	ckFile      string            // COOKIE 文件完整路径
	Cookie      map[string]string // COOKIE 信息
}

func NewCookie() *Cookie {
	Obj := &Cookie{
		BaseURI:     "",
		bSaveCookie: false,
		CookieDir:   "",
		ckFile:      "",
		Cookie:      make(map[string]string),
	}
	return Obj
}

func (self *Cookie) SetURL(URL string) {
	self.BaseURI = URL
	Spr := string(os.PathSeparator)
	self.ckFile = fmt.Sprintf("%s%s%s.cookie", self.CookieDir, Spr, URL)
	self.LoadCookie()
}

func (self *Cookie) SetSaveCookie(IsSave bool) {
	//Config.Set("MiniWeb", "SaveCookie", IsSave)
	self.bSaveCookie = IsSave
}

func (self *Cookie) GetSaveCookie() bool {
	return self.bSaveCookie
}

func (self *Cookie) SetCookieDir(CookieDir string) {
	self.CookieDir = CookieDir
	if !mwCommon.DirExists(CookieDir) {
		os.MkdirAll(CookieDir, 0777)
	}
}

func (self *Cookie) Clear() {
	self.Cookie = make(map[string]string)
}

func (self *Cookie) Count() int {
	return len(self.Cookie)
}

func (self *Cookie) SetCookie(CookieName string, CookieValue string) {
	self.Cookie[CookieName] = CookieValue
}

func (self *Cookie) GetCookie(CookieName string) string {
	return self.Cookie[CookieName]
}

func (self *Cookie) GetAllCookie() string {
	var dwRet string
	for k, v := range self.Cookie {
		dwRet = fmt.Sprintf("%s%s=%s; ", dwRet, k, v)
	}
	return dwRet
}

func (self *Cookie) Dump() {
	for k, v := range self.Cookie {
		dwRet := fmt.Sprintf("%s=%s", k, v)
		fmt.Println(dwRet)
	}
}

func (self *Cookie) LoadCookie() {
	self.Clear()
	if self.bSaveCookie {
		if mwCommon.FileExists(self.ckFile) && !mwCommon.DirExists(self.ckFile) {
			info, err := ioutil.ReadFile(self.ckFile)
			if err == nil {
				json.Unmarshal(info, &self.Cookie)
			}
		}
	}

}

func (self *Cookie) SaveCookie() {
	if self.bSaveCookie {
		if len(self.ckFile) > 0 && !mwCommon.DirExists(self.ckFile) {
			if mwCommon.FileExists(self.ckFile) {
				os.Remove(self.ckFile)
			}
			info, _ := json.Marshal(self.Cookie)
			ioutil.WriteFile(self.ckFile, info, 0777)
		}
	}
}

func (self *Cookie) RemoveSaveCookie() {
	self.Cookie = make(map[string]string)
	if mwCommon.FileExists(self.ckFile) && !mwCommon.DirExists(self.ckFile) {
		os.Remove(self.ckFile)
	}

}

func (self *Cookie) gwcInitCookie() {
}

/*
   def gwcInitCookie(self, URL):
       res = urllib_parse.urlparse(URL)
       # print("返回对象：", res)
       # print("域名", res.netloc)
       self.SetURL(res.netloc)
       if not self.GetSaveCookie():
           self.SetOption(pycurl.COOKIESESSION, True)
           self.SetOption(pycurl.FRESH_CONNECT, True)
       else:
           self.SetOption(pycurl.COOKIESESSION, False)
           self.SetOption(pycurl.FRESH_CONNECT, False)
       if self.Count() > 0:
           self.SetOption(pycurl.COOKIE, self.GetAllCookie())
*/
func (self *Cookie) ParserCookie(Cookie string) {
	//BAIDUID=A4D53EE3E753B76B8A049BE9BE339F3E:FG=1; expires=Wed, 15-Jul-20 17:17:06 GMT; max-age=31536000; path=/; domain=.baidu.com; version=1

	//fmt.Println(Cookie)
	if strings.Contains(Cookie, ";") == true {
		arr := strings.Split(Cookie, ";")
		for i := 0; i < len(arr); i++ {
			//fmt.Println(arr[i])
			self.parserItemCookie(strings.Trim(string(arr[i]), " "))
		}
	} else {
		self.parserItemCookie(strings.Trim(Cookie, " "))
	}
}

func (self *Cookie) parserItemCookie(Cookie string) {
	if strings.Contains(Cookie, "=") == true {
		Index := strings.Index(Cookie, "=")
		sByte := []byte(Cookie)
		Name := strings.Trim(string(sByte[:Index]), " ")
		Index += 1
		Value := strings.Trim(string(sByte[Index:]), " ")

		switch strings.ToLower(Name) {
		case "expires":
		case "max-age":
		case "path":
		case "domain":
		case "secure":
		case "httponly":
		case "version":
		default:
			self.SetCookie(Name, Value)
		}
	}
}

/*
   def ParserCookie(self, SetCookie):
       Logs.Debug("MiniWeb", "Cookie.Parser: {}".format(SetCookie))
       s = SetCookie.split("; ")
       for Info in s:
           if Info.find("=") == -1:    # httponly
               continue

           Index = Info.find("=")
           Name = Info[0:Index].strip()
           Value = Info[Index+1:].strip()

           SkipName = [
               'expires', 'max-age', 'path', 'domain', 'secure', 'httponly', 'version'
           ]
           if Name.lower() in SkipName:
               continue
           self.SetCookie(Name, Value)
           Logs.Debug("MiniWeb", "{0}: {1}".format(Name, Value))
*/
