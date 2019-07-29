package MiniWeb

import (
	"testing"
)

func TestMain(t *testing.T) {
	Config.EnableLogs = false
	Config.ReConfig()

	w := NewMiniWeb()
	w.GetWebCode("http://www.baidu.com/robots.txt")
	t.Log("HTTP.Recv: ", len(w.ResponseText()))

	w.GetWebCode("https://www.baidu.com/robots.txt")
	t.Log("HTTPS.Recv: ", len(w.ResponseText()))
}
