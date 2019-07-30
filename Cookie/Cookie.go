package Cookie

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	mwCommon "github.com/a2si/MiniWeb/Common"
	mwError "github.com/a2si/MiniWeb/mwError"
)

type Cookie struct {
	prv_BaseURI   string            // URL地址
	bSaveCookie   bool              // 是否保存COOKIE
	prv_CookieDir string            // COOKIE 保存路径
	ckFile        string            // COOKIE 文件完整路径
	prv_Cookie    map[string]string // COOKIE 信息
	ObjError      *mwError.TError   // Error Object
}

func NewCookie(errObj *mwError.TError) *Cookie {
	Obj := &Cookie{
		prv_BaseURI:   "",
		bSaveCookie:   false,
		prv_CookieDir: "",
		ckFile:        "",
		prv_Cookie:    make(map[string]string),
		ObjError:      errObj,
	}
	return Obj
}

func (self *Cookie) SetURL(URL string) {
	if len(URL) == 0 {
		return
	}
	self.prv_BaseURI = URL
	Spr := string(os.PathSeparator) // linux=/ windows=\\
	self.ckFile = fmt.Sprintf("%s%s%s.cookie", self.prv_CookieDir, Spr, URL)
	self.LoadCookie()
}

func (self *Cookie) SetSaveCookie(IsSave bool) {
	self.bSaveCookie = IsSave
}

func (self *Cookie) GetSaveCookie() bool {
	return self.bSaveCookie
}

func (self *Cookie) SetCookieDir(CookieDir string, Automkdir bool) {
	self.prv_CookieDir = CookieDir
	if Automkdir == true && !mwCommon.DirExists(CookieDir) {
		os.MkdirAll(CookieDir, 0777)
	}
}

func (self *Cookie) Clear() {
	//fmt.Println("SetCookie::Clear")
	self.prv_Cookie = make(map[string]string)
}

func (self *Cookie) Count() int {
	return len(self.prv_Cookie)
}

// Cookie Name 区分大小写, 查了一下RFC, 没有明确规定
func (self *Cookie) SetCookie(CookieName string, CookieValue string) {
	//fmt.Println("SetCookie: ", CookieName, CookieValue)
	self.prv_Cookie[CookieName] = CookieValue
}

// COOKIE 不存在则返回 ""
func (self *Cookie) GetCookie(CookieName string) string {
	return self.prv_Cookie[CookieName]
}

func (self *Cookie) GetAllCookie() string {
	var dwRet string
	for k, v := range self.prv_Cookie {
		dwRet = fmt.Sprintf("%s%s=%s; ", dwRet, k, v)
	}
	fmt.Println("GetAllCookie: ", dwRet)
	return dwRet
}

func (self *Cookie) Dump() {
	for k, v := range self.prv_Cookie {
		dwRet := fmt.Sprintf("%s=%s", k, v)
		fmt.Println(dwRet)
	}
}

func (self *Cookie) LoadCookie() {
	if self.bSaveCookie {
		self.Clear()
		if mwCommon.DirExists(self.prv_CookieDir) == false {
			return
		}
		if mwCommon.FileExists(self.ckFile) && !mwCommon.DirExists(self.ckFile) {
			info, err := ioutil.ReadFile(self.ckFile)
			if err == nil {
				json.Unmarshal(info, &self.prv_Cookie)
			}
		}
	}

}

func (self *Cookie) SaveCookie() {
	if self.bSaveCookie {
		if mwCommon.DirExists(self.prv_CookieDir) == false {
			return
		}
		if len(self.ckFile) > 0 && !mwCommon.DirExists(self.ckFile) {
			if mwCommon.FileExists(self.ckFile) {
				os.Remove(self.ckFile)
			}
			info, _ := json.Marshal(self.prv_Cookie)
			ioutil.WriteFile(self.ckFile, info, 0777)
		}
	}
}

func (self *Cookie) RemoveSaveCookie() {
	if mwCommon.DirExists(self.prv_CookieDir) == false {
		return
	}
	self.prv_Cookie = make(map[string]string)
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
