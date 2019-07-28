package Network

import (
	"bufio"
	"time"

	DevLogs "github.com/a2si/MiniWeb/DevLogs"
)

func (self *TNet) netIO2BufferIO() {
	self.ioRead = bufio.NewReader(self.Conn)
	self.isClosed = false
}

func (self *TNet) SetTimeOut(TimeOut time.Duration) {
	self.timeOut = TimeOut
	self.Conn.SetDeadline(time.Now().Add(TimeOut * time.Second))
}

func (self *TNet) Send(sByte []byte) error {
	DevLogs.Debug("TNet.Send")
	if len(sByte) == 0 {
		return nil
	}
	_, err := self.Conn.Write(sByte)
	//fmt.Println("Socket.Send: ", len(sByte), n, err)
	return err
}
