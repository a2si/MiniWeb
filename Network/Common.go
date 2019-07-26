package Network

import "time"

type TNet interface {
	Init(Host string, Port string, TimeOut time.Duration) error
	Close()
	SetTimeOut(TimeOut time.Duration)
	Send([]byte) error

	ReadLine() (string, error)
	ReadBytes(Length int) []byte
	ReadChunk() []byte
}
