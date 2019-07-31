package mwNet

import (
	"strconv"
	"sync"

	"fmt"
	"io"
	"strings"

	mwError "github.com/a2si/MiniWeb/mwError"
)

type (
	Reader interface {
		Read(p []byte) (n int, err error)
	}

	TIOBuffer struct {
		ioObjError   *mwError.TError // Error Object
		ioBufferSize int
		ioBuffer     []byte
		ioRead       int
		ioWrite      int
		ioLock       *sync.RWMutex
		ioObject     Reader
		ioClosed     bool
	}
)

const (
	CONFIG_IO_BUFFER_SIZE = 4096
)

func NewIOBuffer(errObj *mwError.TError, ioObject Reader) *TIOBuffer {
	Obj := &TIOBuffer{
		ioObjError:   errObj,
		ioBufferSize: CONFIG_IO_BUFFER_SIZE,
		ioBuffer:     make([]byte, CONFIG_IO_BUFFER_SIZE),
		ioRead:       0,
		ioWrite:      0,
		ioLock:       new(sync.RWMutex),
		ioObject:     ioObject,
		ioClosed:     false,
	}
	return Obj
}

func (self *TIOBuffer) IsClose() bool {
	return self.ioClosed
}

func (self *TIOBuffer) Buffered() int {
	self.ioLock.RLock()
	defer self.ioLock.RUnlock()
	return self.ioWrite - self.ioRead
}

func (self *TIOBuffer) FlushBuffer() {
	if self.ioObject == nil {
		return
	}
	self.ioLock.Lock()
	defer self.ioLock.Unlock()
	var (
		ioSize  int = CONFIG_IO_BUFFER_SIZE
		ioWrite int
		err     error
	)
	// 刷新缓冲读指针
	if self.ioRead > 0 {
		copy(self.ioBuffer, self.ioBuffer[self.ioRead:self.ioWrite])
		self.ioWrite -= self.ioRead
		self.ioRead = 0
	}
	// 缓冲已满
	if self.ioWrite >= ioSize {
		fmt.Println("缓冲已满")
		return
	}

	ioWrite, err = self.ioObject.Read(self.ioBuffer[self.ioWrite:])
	if ioWrite < 0 {
		self.ioObjError.IOReadByNegative()
		return
	}
	self.ioWrite += ioWrite
	if err != nil {
		if err == io.EOF {
			self.ioClosed = true
			return
		}
		str := err.Error()
		if strings.Contains(str, "read: connection refused") {
			self.ioClosed = true
			return
		}
		self.ioClosed = true
		self.ioObjError.SetIOError(err)
	}
}

func (self *TIOBuffer) protectRead(Length int) []byte {
	if self.Buffered() < Length {
		Length = self.Buffered()
	}
	if Length == 0 {
		self.FlushBuffer()
		return nil
	}
	var (
		sByte    []byte = make([]byte, Length)
		copySize int    = 0
	)
	self.ioLock.Lock()
	copySize = copy(sByte, self.ioBuffer[self.ioRead:self.ioWrite])
	self.ioRead += copySize
	self.ioLock.Unlock()
	return sByte
}

func (self *TIOBuffer) ReadByte() byte {
	// Check It Object
	if self.ioObject == nil {
		return 0
	}

	for self.ioRead == self.ioWrite {
		if self.ioClosed || self.ioObjError.IsError() {
			return 0
		}
		self.FlushBuffer()
	}
	self.ioLock.Lock()
	c := self.ioBuffer[self.ioRead]
	self.ioRead++
	self.ioLock.Unlock()
	return c
}

func (self *TIOBuffer) ReadLine() string {
	if self.ioObject == nil {
		return ""
	}
	var (
		sByte []byte
		sBit  byte
	)
	for {
		if self.ioClosed || self.ioObjError.IsError() {
			return string(sByte)
		}
		if self.Buffered() >= 1 {
			sBit = self.ReadByte()
			if sBit == '\n' {
				iCount := len(sByte)
				iCount--
				if sByte[iCount] == '\r' {
					sByte = sByte[:iCount]
				}
				return string(sByte)
			}
			sByte = append(sByte, sBit)
		} else {
			self.FlushBuffer()
		}
		if self.ioClosed || self.ioObjError.IsError() {
			return string(sByte)
		}
	}
}

func (self *TIOBuffer) ReadBytes(Length int) []byte {
	// Check It Object
	if self.ioObject == nil || Length == 0 {
		return nil
	}
	var (
		sByte     []byte
		bReadSize int
	)
	bReadSize = self.Buffered()
	// if the buffer size > need to read
	if bReadSize >= Length {
		return self.protectRead(Length)
	}
	if bReadSize > 0 {
		Length -= bReadSize
		sByte = self.protectRead(bReadSize)
	}

	for {
		if Length == 0 || self.ioClosed || self.ioObjError.IsError() {
			return sByte
		}
		self.FlushBuffer()
		bReadSize = self.Buffered()
		if bReadSize > Length {
			bReadSize = Length
		}
		if bReadSize >= 1 {
			tempByte := self.protectRead(bReadSize)
			Length -= len(tempByte)
			sByte = append(sByte, tempByte...)
		}
	}
}

func (self *TIOBuffer) ReadToCBHook(Event int, fn CBHOOK) {
	if self.ioObject == nil {
		return
	}
	for {
		if self.ioClosed || self.ioObjError.IsError() {
			return
		}
		self.FlushBuffer()
		bReadSize := self.Buffered()
		fn(Event, nil, self.protectRead(bReadSize))
	}
}

func (self *TIOBuffer) ReadToEOF() []byte {
	if self.ioObject == nil {
		return nil
	}
	var (
		sByte     []byte
		bReadSize int
	)
	for {
		if self.ioClosed || self.ioObjError.IsError() {
			return sByte
		}
		bReadSize = self.Buffered()
		if bReadSize >= 1 {
			tempByte := self.protectRead(bReadSize)
			sByte = append(sByte, tempByte...)
		}
		self.FlushBuffer()
	}
	return sByte
}

func (self *TIOBuffer) ReadChunk() []byte {
	if self.ioObject == nil {
		return nil
	}
	var (
		sByte []byte
		Text  string
		n     int64
		err   error
	)
	for {
		if self.ioClosed || self.ioObjError.IsError() {
			return nil
		}
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
			}
		}
	}
	return sByte
}
