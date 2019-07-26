package Network

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"time"

	mwCommon "github.com/MiniWeb/Common"
)

type NetTls struct {
	Conn   *tls.Conn
	ioRead *bufio.Reader
}

func (self *NetTls) Init(Host string, Port string, TimeOut time.Duration) error {
	var err error
	TlsTime := &net.Dialer{Timeout: TimeOut * time.Second}
	TlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	self.Conn, err = tls.DialWithDialer(TlsTime, "tcp", Host+":"+Port, TlsConfig)
	if err == nil {
		self.ioRead = bufio.NewReader(self.Conn)
	}
	return err
}

func (self *NetTls) Close() {
	self.Conn.Close()
}

func (self *NetTls) SetTimeOut(TimeOut time.Duration) {
	self.Conn.SetDeadline(time.Now().Add(TimeOut * time.Second))
}

func (self *NetTls) Send(sByte []byte) error {
	self.Conn.Write(sByte)
	return nil
}

func (self *NetTls) ReadLine() (string, error) {
	line, isprefix, err := self.ioRead.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = self.ioRead.ReadLine()
		line = append(line, bs...)
	}
	return string(line), err
}

func (self *NetTls) ReadBytes(Length int) []byte {
	var (
		sByte []byte
	)
	for {
		b, e := self.ioRead.ReadByte()
		if e != nil {
			fmt.Println(e)
			return sByte
		}
		sByte = append(sByte, b)
		if len(sByte) == Length {
			return sByte
		}
	}
}

func (self *NetTls) ReadChunk() []byte {
	var (
		sByte []byte
		Text  string
		n     int64
		err   error
	)
	for {
		Text, err = self.ReadLine()
		if err != nil {
			fmt.Println(err)
			return sByte
		}
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
				sByte = mwCommon.CopyMergeSlice(sByte, self.ReadBytes(int(n)))
			}
		}
	}
	return sByte
}
