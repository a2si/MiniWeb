package Common

import (
	"os"
)

func CopyMergeSlice(s1 []byte, s2 []byte) []byte {
	slice := make([]byte, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}

func IsExists(Path string) bool {
	f, _ := os.Stat(Path)
	if f == nil {
		return false
	}
	return false
}

func FileExists(FileName string) bool {
	f, _ := os.Stat(FileName)
	if f == nil {
		return false
	}
	return !f.IsDir()
}

func DirExists(DirName string) bool {
	f, _ := os.Stat(DirName)
	if f == nil {
		return false
	}
	return f.IsDir()
}
