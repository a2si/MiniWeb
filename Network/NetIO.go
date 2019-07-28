package Network

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	/*
		不预读会导致缓冲没有数据
		self.ioRead.ReadByte() 不一定会触发 内置 Reader.fill() 函数
		self.ioRead.Peek() 手动触发缓冲更新 Reader.fill()
		但是预读过大会导致读取超时
	*/
	NET_PEEK_SIZE = 1
)

func (self *TNet) SendPacket(sByte []byte) {
	var (
		sendSize int = len(sByte)
		realSize int = 0
		err      error
	)
	realSize, err = self.Conn.Write(sByte)
	if sendSize != realSize {
		fmt.Println("捕捉到发送数据大小不一致")
		fmt.Println("[]byte size: ", len(sByte))
		fmt.Println("sendSize: ", sendSize)
		fmt.Println("realSize: ", realSize)
		fmt.Println("ErrMsg: ", err)
		panic("特殊情况需要处理")
	}
	if err != nil {
		fmt.Println("[]byte size: ", len(sByte))
		fmt.Println("sendSize: ", sendSize)
		fmt.Println("realSize: ", realSize)
		fmt.Println("ErrMsg: ", err)
	}

}

func (self *TNet) readError(err error) bool {
	if err == nil {
		return false
	}
	if err == io.EOF { // EOF == 网络断开
		self.isClosed = true
		return false
	}
	errMsg := err.Error()
	if strings.Contains(errMsg, "i/o timeout") {
		self.ObjError.SocketReadTimeout()
		return true
	}
	//read: connection reset by peer
	if strings.Contains(errMsg, "read: connection reset by peer") {
		self.ObjError.SocketReadReset()
		return true
	}
	//fmt.Println("readError", err)
	panic(err)
	return true
}

func (self *TNet) ReadLine() string {
	var (
		sByte []byte
		sBit  byte
		err   error
	)
	for {
		if self.ioRead.Buffered() >= 1 {
			sBit, err = self.ioRead.ReadByte()
			if sBit == '\n' {
				iCount := len(sByte)
				iCount--
				if sByte[iCount] == '\r' {
					sByte = sByte[:iCount]
				}
				return string(sByte)
			}
			sByte = append(sByte, sBit)
		}
		if self.isClosed || self.readError(err) {
			return string(sByte)
		}
		self.ioRead.Peek(NET_PEEK_SIZE)
		//time.Sleep(10 * time.Millisecond)
	}
}

func (self *TNet) ReadBytes(Length int) []byte {
	var (
		sByte     []byte
		err       error
		iCount    int
		bReadSize int
	)
	bReadSize = 0
	for {
		iCount++
		if self.ioRead.Buffered() >= 1 {
			Length = Length - bReadSize
			tempByte := make([]byte, Length)
			bReadSize, err = self.ioRead.Read(tempByte)
			tempByte = tempByte[:bReadSize]

			// 对性能影响不大, 非频繁,大量的调用, append 相对更稳定一些
			sByte = append(sByte, tempByte...)
			//sByte = mwCommon.CopyMergeSlice(sByte, tempByte)
			//fmt.Println(hex.EncodeToString(sByte))
		}

		if bReadSize == Length || self.isClosed || self.readError(err) {
			return sByte
		}
		self.ioRead.Peek(Length - bReadSize)
	}
}

func (self *TNet) ReadToEOF() []byte {
	var (
		sByte []byte
		sBit  byte
		err   error
	)
	for {
		if self.ioRead.Buffered() >= 1 {
			sBit, err = self.ioRead.ReadByte()
			sByte = append(sByte, sBit)
		}
		if self.isClosed || self.readError(err) {
			return sByte
		}
		self.ioRead.Peek(NET_PEEK_SIZE)
	}
}

func (self *TNet) ReadChunk() []byte {
	var (
		sByte []byte
		Text  string
		n     int64
		err   error
	)
	for {
		Text = self.ReadLine()
		// Chunk == 0, Next = \r\n
		if Text == "" && n == 0 {
			return sByte
		}
		if len(Text) > 0 {
			n, err = strconv.ParseInt(Text, 16, 32)
			if err != nil {
				fmt.Println(err)
				return sByte
			}
			if n > 0 {
				sByte = append(sByte, self.ReadBytes(int(n))...)
				//sByte = mwCommon.CopyMergeSlice(sByte, self.ReadBytes(int(n)))
			}
		}
	}
	return sByte
}
