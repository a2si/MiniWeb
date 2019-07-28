package Common

import (
	"testing"
)

const (
	NNNCount = 1000
)

func TestMain(t *testing.T) {

}

// 性能测试
func BenchmarkCMS(b *testing.B) {
	var (
		sByte   []byte = []byte("b.N会根据函数的运行时间取一个合适的值")
		tmpByte []byte
	)
	b.N = NNNCount
	for i := 0; i < NNNCount; i++ {
		tmpByte = CopyMergeSlice(sByte, sByte)
	}
	sByte = tmpByte
	//sByte = append(sByte, self.ReadBytes(int(n))...)
	//sByte = mwCommon.CopyMergeSlice(sByte, self.ReadBytes(int(n)))
}

func BenchmarkAppend(b *testing.B) {
	var (
		sByte   []byte = []byte("b.N会根据函数的运行时间取一个合适的值")
		tmpByte []byte
	)
	b.N = NNNCount
	for i := 0; i < NNNCount; i++ {
		tmpByte = append(sByte, sByte...)
	}
	sByte = tmpByte
	//sByte = append(sByte, self.ReadBytes(int(n))...)
	//sByte = mwCommon.CopyMergeSlice(sByte, self.ReadBytes(int(n)))
}
