package Core

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net"
	"strconv"
	"strings"
	"time"

	mwCommon "github.com/a2si/MiniWeb/Common"
	DevLogs "github.com/a2si/MiniWeb/DevLogs"
	mwNet "github.com/a2si/MiniWeb/Network"
	mwConst "github.com/a2si/MiniWeb/mwConst"
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
		Host    string      = self.URL.GetHost()
		Port    string      = self.URL.GetPort()
		NetWork *mwNet.TNet = mwNet.NewNet(self.ObjError, self.Proxy)
		rawDial net.Dialer  = net.Dialer{Timeout: self.TimeOutConnect * time.Second}
		RawConn net.Conn    = NetWork.StartNetwork(rawDial, Host, Port)
	)

	// 网络层 startNetwork
	if RawConn == nil {
		return self.ObjError.GetErrorCode()
	}

	// 代理层
	if self.Proxy.Enable() {
		/*
			如果代理启用
			startNetwork 时替换连接的IP
		*/
		switch self.Proxy.GetProxyType() {
		case mwConst.PROXY_TYPE_HTTP:
		/*
			buildReqHeader 修改 HTTP 请求
		*/
		case mwConst.PROXY_TYPE_HTTPS:
			/*
				CONNECT 连接后, 发送正文或TLS认证
			*/
			NetWork.InitProxyHttps(RawConn, mwNet.Host2IP(Host), Port)
		case mwConst.PROXY_TYPE_SOCKS4:
			NetWork.InitProxySocks4(RawConn, mwNet.Host2IP(Host), Port)
		case mwConst.PROXY_TYPE_SOCKS4A:
			NetWork.InitProxySocks4a(RawConn, Host, Port)
		case mwConst.PROXY_TYPE_SOCKS5:
			NetWork.InitProxySocks5(RawConn, Host, Port)
			//NetWork.InitProxySocks5(RawConn, mwNet.Host2IP(Host), Port)
		}
	}
	if self.ObjError.IsError() {
		return self.ObjError.GetErrorCode()
	}

	// 传输层
	if self.URL.IsTls() {
		NetWork.InitTLS(RawConn)
	} else {
		NetWork.InitTCP(RawConn)
	}

	// 网络完毕
	defer NetWork.Close()

	// 网络通讯
	self.ReqHeader.SetHeader("Host", Host) // 修正请求信息
	self.RspHeader.ClearHeader()           // 清空返回信息
	NetWork.SetTimeOut(self.TimeOut)       // 设置超时
	tmpBody := self.genReqBody()           // 生成HTTP协议正文
	tmpHead := self.genReqHeader()         // 生成HTTP协议头
	NetWork.Send(tmpHead)                  // 发送 请求头
	NetWork.Send(tmpBody)                  // 发送 请求正文

	self.readRspHeader(NetWork) // Response Header
	if self.ObjError.IsError() {
		return self.ObjError.GetErrorCode()
	}
	self.Result = self.readRspBody(NetWork) // Response Body
	if self.ObjError.IsError() {
		return self.ObjError.GetErrorCode()
	}

	return self.StatusCode
}

func (self *WebCore) readRspBody(NetWork *mwNet.TNet) []byte {
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
	// 如果链接被关闭, 读取所有数据
	if self.RspHeader.HeaderExists("Connection: close") {
		sByte = append(sByte, NetWork.ReadToEOF()...)
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
	//fmt.Println(string(sByte))
	//fmt.Println(fmt.Sprintf("%x", md5.Sum(sByte)))
	return sByte
}

func (self *WebCore) readRspHeader(NetWork *mwNet.TNet) {
	DevLogs.Debug("WebCore.readRspHeader")
	for {
		Text := NetWork.ReadLine()
		//fmt.Println(Text)
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
		Query      string = "?" + self.URL.GetEncode()
		MethodPath string = fmt.Sprintf("%s%s", self.URL.GetPath(), Query)
		mpSize     int    = len(MethodPath) - 1
		dwRet      string
	)
	if []byte(MethodPath)[mpSize] == '?' {
		MethodPath = string([]byte(MethodPath)[:mpSize])
	}
	/*
		如果是HTTP代理, 仅修改 GET Script == Get scheme://host:port/Script
		暂时未考虑 user:pass
	*/
	if self.Proxy.GetProxyType() == mwConst.PROXY_TYPE_HTTP {
		MethodPath = fmt.Sprintf("%s://%s:%s%s", self.URL.GetScheme(), self.URL.GetHost(), self.URL.GetPort(), MethodPath)
	}
	dwRet = fmt.Sprintf("%s %s HTTP/1.1", self.Method, MethodPath)
	dwRet = dwRet + self.ReqHeader.GetAllHeader() + "\r\n\r\n"
	//fmt.Println(dwRet)
	return dwRet
}
