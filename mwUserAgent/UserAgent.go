package UserAgent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
)

type UserAgent struct {
	UA map[string][]string
}

func NewUserAgent() *UserAgent {
	Obj := &UserAgent{
		UA: make(map[string][]string),
	}
	Obj.loadUAJson()
	return Obj
}

func (self *UserAgent) getUserAgentJsonFile() string {
	Path, _ := os.Getwd()
	Spr := string(os.PathSeparator)
	FullPath := fmt.Sprintf("%s%sConfig%sUserAgent.json", Path, Spr, Spr)
	return FullPath
}

func (self *UserAgent) loadUAJson() {

	info, err := ioutil.ReadFile(self.getUserAgentJsonFile())
	if err == nil {
		json.Unmarshal(info, &self.UA)
	} else {
		json.Unmarshal(jsonInfo, &self.UA)
	}
}

func (self *UserAgent) IE() string {
	return randUA(self.UA["ie"])
}

func (self *UserAgent) Opera() string {
	return randUA(self.UA["opera"])
}

func (self *UserAgent) Firefox() string {
	return randUA(self.UA["firefox"])
}

func (self *UserAgent) Chrome() string {
	return randUA(self.UA["chrome"])
}

func (self *UserAgent) Safari() string {
	return randUA(self.UA["safari"])
}

func (self *UserAgent) Random() string {
	str := []string{"ie", "opera", "firefox", "chrome", "safari"}
	iCount := len(str)
	iCount = rand.Intn(iCount)
	usStr := str[iCount]
	return randUA(self.UA[usStr])
}

func randUA(UAList []string) string {
	iCount := len(UAList)
	if iCount == 0 {
		return "MiniWeb/3.1.0"
	}
	iCount = rand.Intn(iCount)
	return UAList[iCount]
}

/*
   def __getattr__(self, Name):
       Cmd = Name.lower()
       if Cmd == "ie":
           return self.__Random(Cmd)
       elif Cmd == "opera":
           return self.__Random(Cmd)
       elif Cmd == "firefox":
           return self.__Random(Cmd)
       elif Cmd == "ff":
           return self.__Random("firefox")
       elif Cmd == "chrome":
           return self.__Random(Cmd)
       elif Cmd == "safari":
           return self.__Random(Cmd)
       else:
           Keys = ["ie", "opera", "firefox", "chrome", "safari"]
           Len = len(Keys)
           Len -= 1
           Pos = random.randint(0, Len)
           return self.__Random(Keys[Pos])

   def __Random(self, Name):
       Arr = CONFIG_BROWSERS[Name]
       Len = len(Arr)
       Len -= 1
       Pos = random.randint(0, Len)
       return Arr[Pos]
*/
