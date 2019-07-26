package Core

import (
	mwCommon "MiniWeb/Common"
	DevLogs "MiniWeb/DevLogs"
	mwNet "MiniWeb/Network"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"strings"
)

func (self *WebCore) InitHeader() {
	DevLogs.Debug("WebCore.InitHeader")
	self.ReqHeader.SetHeader("User-Agent", self.UserAgent)
	self.ReqHeader.SetHeader("Accept", "*/*")
	self.ReqHeader.SetHeader("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	self.ReqHeader.SetHeader("Accept", "ext/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	self.ReqHeader.SetHeader("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	self.ReqHeader.SetHeader("Accept-Encoding", "gzip, deflate")
	self.ReqHeader.SetHeader("Cache-Control", "max-age=0")
}

func (self *WebCore) SetMethod(Method string) {
	DevLogs.Debug("WebCore.SetMethod")
	self.Method = strings.ToUpper(Method)
}

func (self *WebCore) AddPost(Name string, Value string) {
	DevLogs.Debug("WebCore.AddPost")
	self.PostData[Name] = Value
}

func (self *WebCore) SendRequest() int {
	DevLogs.Debug("WebCore.SendRequest")
	var (
		NetWork mwNet.TNet
		Host    string = self.URL.GetHost()
		Port    string = self.URL.GetPort()
	)
	self.ReqHeader.SetHeader("Host", Host)
	self.RspHeader.ClearHeader()

	if self.URL.IsTls() {
		NetWork = &mwNet.NetTls{}
	} else {
		NetWork = &mwNet.NetTCP{}
	}

	/*	代理这里添加
			拦截 Host, Port 重定向到本地
			本地启动一个代理连接进行代理加载工作
		Proxy.SetRealAddr(Host, Port)
		Proxy.SetProxyType(HTTP/HTTPS/SOCKS/Other)
		Host = "0.0.0.0"
		Port = "Rand.Port"
	*/
	if err := NetWork.Init(Host, Port, self.TimeOutConnect); err != nil {
		DevLogs.Error("WebCore.SendRequest.ConnectError Error=" + err.Error())
		return ERROR_CODE_TOC
	}

	defer NetWork.Close()
	// 设置超时
	NetWork.SetTimeOut(self.TimeOut)
	tmpBody := self.genReqBody()
	tmpHead := self.genReqHeader()
	NetWork.Send(tmpHead)
	NetWork.Send(tmpBody)

	self.readRspHeader(NetWork) // Response Header
	if getErrorCode() != ERROR_CODE_NO_ERROR {
		return getErrorCode()
	}
	self.Result = self.readRspBody(NetWork) // Response Body
	if getErrorCode() != ERROR_CODE_NO_ERROR {
		return getErrorCode()
	}

	return self.StatusCode
}

func (self *WebCore) readRspBody(NetWork mwNet.TNet) []byte {
	DevLogs.Debug("WebCore.readRspBody")
	var (
		sByte            []byte
		TransferEncoding string
		ContentEncoding  string
		ContentLength    int
	)

	if self.RspHeader.HeaderExists("Content-Length") {
		ContentLength, _ = strconv.Atoi(self.RspHeader.GetHeader("Content-Length"))
	}
	if self.RspHeader.HeaderExists("Transfer-Encoding") {
		TransferEncoding = self.RspHeader.GetHeader("Transfer-Encoding")
		TransferEncoding = strings.ToLower(TransferEncoding)
	}
	if self.RspHeader.HeaderExists("Content-Encoding") {
		ContentEncoding = self.RspHeader.GetHeader("Content-Encoding")
		ContentEncoding = strings.ToLower(ContentEncoding)
	}
	/*
		chunked 		== 不定长度
		ContentLength 	== 指定长度
		代码暂时观察, 如果有异常, 再修改逻辑
	*/
	if ContentLength != 0 {
		sByte = NetWork.ReadBytes(ContentLength)
	}
	if strings.Contains(TransferEncoding, "chunked") {
		sByte = NetWork.ReadChunk()
	}

	if strings.Contains(ContentEncoding, "gzip") {
		// 一个缓存区压缩的内容
		buf := bytes.NewBuffer(sByte)

		// 解压刚压缩的内容
		flateReader, _ := gzip.NewReader(buf)
		defer flateReader.Close()

		sByte, _ = ioutil.ReadAll(flateReader)

	}
	//fmt.Println(TransferEncoding)
	//fmt.Println(ContentEncoding)
	return sByte
}

func (self *WebCore) readRspHeader(NetWork mwNet.TNet) {
	DevLogs.Debug("WebCore.readRspHeader")
	for {
		Text, err := NetWork.ReadLine()
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") {
				setError(ERROR_CODE_READ_TIME, ERROR_MSG_READ_TIME)
				return
			}
			fmt.Println(err)
			return
		}
		if Text == "\r\n" || Text == "" {
			break
		}
		self.rspParserHeaderLine(Text)
	}
}

func (self *WebCore) genReqBody() []byte {
	DevLogs.Debug("WebCore.genReqBody")
	var (
		tempBody []byte
		tmpStr   string
		k        string
		v        string
	)
	if self.havePostFile() == false {
		for k, v = range self.PostData {
			if len(tmpStr) == 0 {
				tmpStr = fmt.Sprintf("%s=%s", k, v)
			} else {
				tmpStr = fmt.Sprintf("%s&%s=%s", tmpStr, k, v)
			}
		}
		tempBody = []byte(tmpStr)
	} else {
		Round := fmt.Sprintf("%x", md5.Sum([]byte(self.UserAgent)))

		mimeByte := new(bytes.Buffer)
		mime := multipart.NewWriter(mimeByte)
		mime.SetBoundary(Round)

		self.ReqHeader.SetHeader("Content-Type", mime.FormDataContentType())

		tempBody = []byte("")
		for k, v = range self.PostData {
			if string([]byte(k)[:1]) != "@" {
				mime.WriteField(k, v)
			} else {
				if mwCommon.FileExists(v) == true {
					k = string([]byte(k)[1:])
					fileIO, _ := mime.CreateFormFile(k, v)
					data, _ := ioutil.ReadFile(v)
					fileIO.Write(data)
				} else {
					DevLogs.Warn("POST.File Not Exists")
				}
			}
		}
		mime.Close()
		tempBody = mimeByte.Bytes()
	}
	self.PostData = make(map[string]string)

	contentLength := len(tempBody)
	if contentLength > 0 {
		self.ReqHeader.SetHeader("Content-Length", strconv.Itoa(contentLength))
	}
	return tempBody
}

func (self *WebCore) havePostFile() bool {
	DevLogs.Debug("WebCore.havePostFile")
	for k, _ := range self.PostData {
		if string([]byte(k)[:1]) == "@" {
			return true
		}
	}
	return false
}

func (self *WebCore) genReqHeader() []byte {
	DevLogs.Debug("WebCore.genReqHeader")
	Header := self.buildReqHeader()
	return []byte(Header)
}

func (self *WebCore) buildReqHeader() string {
	DevLogs.Debug("WebCore.buildReqHeader")
	var (
		Query string = self.URL.GetEncode()
		dwRet string
	)
	if len(Query) == 0 {
		dwRet = fmt.Sprintf("%s %s HTTP/1.1", self.Method, self.URL.GetPath())
	} else {
		dwRet = fmt.Sprintf("%s %s?%s HTTP/1.1", self.Method, self.URL.GetPath(), Query)
	}
	dwRet = dwRet + self.ReqHeader.GetAllHeader() + "\r\n\r\n"
	return dwRet
}
