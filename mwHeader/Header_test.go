package Header

import (
	"testing"
)

func TestMain(t *testing.T) {

	Head := NewHeader()

	Head.SetHeader("test", "Value")
	if Head.GetHeader("test") != "Value" {
		t.Error("Head.GetHeader Error")
	}
}
