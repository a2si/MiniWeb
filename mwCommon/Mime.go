package Common

import (
	"mime"
)

func MimeType(Ext string) string {
	str := mime.TypeByExtension(Ext)
	if len(str) == 0 {
		return "application/octet-stream"
	}
	return str
}

func MimeExt(ExtType string) string {
	str, err := mime.ExtensionsByType(ExtType)
	if err == nil {
		return str[0]
	}
	return "bin"
}
