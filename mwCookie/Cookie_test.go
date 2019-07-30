package Cookie

import (
	"testing"
)

func TestMain(t *testing.T) {

	ck := NewCookie(nil)
	ck.SetCookie("A", "1")
	ck.SetCookie("a", "2")

	t.Log("Cookie.Read: A=", ck.GetCookie("A"))
	t.Log("Cookie.Read: a=", ck.GetCookie("a"))
	t.Log("Cookie.Read: nil=", ck.GetCookie("nil"))
}
