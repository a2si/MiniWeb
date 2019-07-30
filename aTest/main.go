package main

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/robertkrimen/otto"

	mwWeb "github.com/a2si/MiniWeb"
	mwConfig "github.com/a2si/MiniWeb/mwConfig"
	mwConst "github.com/a2si/MiniWeb/mwConst"
	mwNet "github.com/a2si/MiniWeb/mwNet"
	JavaScript "github.com/dop251/goja"
)

var (
	_ = JavaScript.AssertFunction
	_ = mwWeb.NewMiniWeb
	_ = mwConst.PROXY_TYPE_HTTP
)

func main() {
	fmt.Println("Hello World!")

	mwConfig.SetConfig("Logs.Enable", false)
	w := mwWeb.NewMiniWeb()
	w.SetUserAgent("Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:68.0) Gecko/20100101 Firefox/68.0")

	w.SetTimeOutConnect(100)
	var (
		wsNet    *mwNet.TNet = nil
		isClosed bool        = false
	)

	go func() {
		StatusCode := w.ConnectWebSocket("wss://echo.websocket.org/", func(Event int, ObjNet *mwNet.TNet, data []byte) {
			switch Event {
			case mwNet.EVENT_OBJECT:

				wsNet = ObjNet

			case mwNet.EVENT_RECV_TEXT:
				fmt.Println("Recv.Text: ", string(data))
			case mwNet.EVENT_RECV_BINARY:
				fmt.Println("Recv.Binary: ", data)
			case mwNet.EVENT_CLOSE:
				fmt.Println("Recv.Close")
				if len(data) >= 2 {
					Code := binary.BigEndian.Uint16(data[:2])
					Msg := string(data[2:])
					fmt.Println("Close.Code: ", Code)
					fmt.Println("Close.Msg: ", Msg)
				}
				isClosed = true
			}
		})
		isClosed = true
		fmt.Println("StatusCode: ", StatusCode)
		fmt.Println(w.GetErrorCode())
		fmt.Println(w.GetErrorMsg())
	}()
	for {
		time.Sleep(2 * time.Second)
		if isClosed == true {
			break
		}
		if wsNet != nil {
			fmt.Println("===============")
			wsNet.WebSocketSendText("Hello Test Server")
			wsNet.WebSocketSendPing([]byte("ping msg"))
		}
	}
}

func testDecJs() {

	/*
		info, err := ioutil.ReadFile("./test.js")
		if err != nil {
			fmt.Println(err)
			return
		}
		jsCode := string(info)
		fmt.Println("=============================================")
		ottoJsCode := decJavaScript(jsCode)
		fmt.Println("OTTO::Result: ", ottoJsCode)
		fmt.Println("=============================================")
		gojaJsCode := decJavaScript2(jsCode)
		fmt.Println("GOJA::Result: ", gojaJsCode)
		return
	*/

	mwConfig.SetConfig("Logs.Enable", false)
	w := mwWeb.NewMiniWeb()
	w.SetUserAgent("Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:68.0) Gecko/20100101 Firefox/68.0")

	for {
		start := time.Now().UnixNano()
		StatusCode := w.GetWebCode("https://*******")
		fmt.Println(StatusCode)
		recvStr := w.ResponseText()
		fmt.Printf("Recv.Size: %d\n", len(recvStr))
		fmt.Println(recvStr)
		if StatusCode == 521 {
			Index := strings.Index(recvStr, ">")
			recvStr = string([]byte(recvStr[Index+1:]))
			Index = strings.LastIndex(recvStr, "<")
			recvStr = string([]byte(recvStr[:Index]))
			fmt.Println(recvStr)
			recvStr = decJavaScript(recvStr)
			fmt.Println(recvStr)

			SymIndex := strings.Index(recvStr, "=")
			Name := string([]byte(recvStr)[:SymIndex])
			Value := string([]byte(recvStr)[SymIndex+1:])
			fmt.Println(Name, Value)
			w.Cookie().SetCookie(Name, Value)
		} else {
			break
		}
		end := time.Now().UnixNano()
		fmt.Printf("执行消耗的时间为:%v\n", end-start)
	}
}

func decJavaScript2(Script string) string {
	vm := otto.New()
	/*
	 */
	vm.Set("eval", func(call otto.FunctionCall) otto.Value {
		fmt.Println("OTTO::jsEval::Enter")
		jsDecCode := call.ArgumentList[0].String()
		fmt.Println("OTTO::jsEval: ", jsDecCode)
		v, e := vm.Run(jsDecCode)
		fmt.Println("OTTO::jsEval: ", v, e)
		fmt.Println("OTTO::jsEval::Leave")
		return otto.Value{}
	})
	vm.Set("Logs", func(call otto.FunctionCall) otto.Value {
		fmt.Println("jsLogs: ", call.ArgumentList[0].String())
		return otto.Value{}
	})

	vm.Set("setTimeout", func(call otto.FunctionCall) otto.Value {
		fmt.Println("jsSetTimeout: ", call.ArgumentList[0].String())
		fmt.Println("jsSetTimeout: ", call.ArgumentList[1].String())
		return otto.Value{}
	})

	vm.Set("document", func(call otto.FunctionCall) otto.Value {
		fmt.Println("document: ", call.ArgumentList[0].String())
		return otto.Value{}
	})
	dom, _ := vm.Get("document")
	dom.Object().Set("attachEvent", func(call otto.FunctionCall) otto.Value {
		jsDecCode := call.ArgumentList[1].String()
		//fmt.Println("jsAttachEvent: ", call.ArgumentList[0].String())
		//fmt.Println("jsAttachEvent: ", jsDecCode)
		//vm.Set("onreadystatechange", jsDecCode)
		fmt.Println("jsAttachEvent::Enter")
		fmt.Println("jsAttachEvent: ", jsDecCode)
		v, e := vm.Run("var a= " + jsDecCode + ";a();")
		//v, e := vm.RunString(jsDecCode)
		fmt.Println(v.String(), e)
		fmt.Println("jsAttachEvent::Leave")
		return otto.Value{}
	})
	dom.Object().Set("cookie", func(call otto.FunctionCall) otto.Value {
		fmt.Println("jsDocumentCookie: ", call.ArgumentList[0].String())
		return otto.Value{}
	})

	vm.Set("window", func(call otto.FunctionCall) otto.Value {
		fmt.Println("window: ", call.ArgumentList[0].String())
		return otto.Value{}
	})
	//vm.Get("window").Object().Set("headless", "undefined")

	v, err := vm.Run(Script) //默认输出最后一个
	if err != nil {
		fmt.Println(err, v)
	}
	zz, _ := vm.Get("document")
	z, _ := zz.Object().Get("cookie")
	str, _ := z.ToString()
	return str
}

func decJavaScript(Script string) string {
	vm := JavaScript.New()
	vm.Set("eval", func(call JavaScript.ConstructorCall) *JavaScript.Object {
		fmt.Println("GOJA::jsEval::Enter")
		jsDecCode := call.Arguments[0].String()
		fmt.Println("GOJA::jsEval: ", jsDecCode)
		v, e := vm.RunString(jsDecCode)
		fmt.Println("GOJA::jsEval: ", v, e)
		fmt.Println("GOJA::jsEval::Leave")
		return nil
	})
	vm.Set("Logs", func(call JavaScript.ConstructorCall) *JavaScript.Object {
		fmt.Println("jsLogs: ", call.Arguments[0].String())
		return nil
	})

	vm.Set("setTimeout", func(call JavaScript.ConstructorCall) *JavaScript.Object {
		fmt.Println("jsSetTimeout: ", call.Arguments[0].String())
		fmt.Println("jsSetTimeout: ", call.Arguments[1].String())
		return nil
	})

	vm.Set("document", func(call JavaScript.ConstructorCall) *JavaScript.Object {
		fmt.Println("document: ", call.Arguments[0].String())
		return nil
	})
	dom := vm.Get("document").ToObject(vm)
	dom.Set("attachEvent", func(call JavaScript.ConstructorCall) *JavaScript.Object {
		jsDecCode := call.Arguments[1].String()
		//fmt.Println("jsAttachEvent: ", call.Arguments[0].String())
		//fmt.Println("jsAttachEvent: ", jsDecCode)
		//vm.Set("onreadystatechange", jsDecCode)
		fmt.Println("jsAttachEvent::Enter")
		fmt.Println("jsAttachEvent: ", jsDecCode)
		v, e := vm.RunString("var a= " + jsDecCode + ";a();")
		//v, e := vm.RunString(jsDecCode)
		fmt.Println(v.String(), e)
		fmt.Println("jsAttachEvent::Leave")
		return nil
	})
	dom.Set("cookie", func(call JavaScript.ConstructorCall) *JavaScript.Object {
		fmt.Println("jsDocumentCookie: ", call.Arguments[0].String())
		return nil
	})

	vm.Set("window", func(call JavaScript.ConstructorCall) *JavaScript.Object {
		fmt.Println("window: ", call.Arguments[0].String())
		return nil
	})
	//vm.Get("window").ToObject(vm).Set("headless", "undefined")

	v, err := vm.RunString(Script + ";") //默认输出最后一个
	if err != nil {
		fmt.Println(err, v)
	}
	//println(v.Export().(string))
	return vm.Get("document").ToObject(vm).Get("cookie").String()
}
