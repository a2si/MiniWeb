package mwNet

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type (
	CBHOOK func(Event int, ObjNet *TNet, data []byte)
)

const (
	EVENT_OBJECT      = 1
	EVENT_RECV_TEXT   = 2
	EVENT_RECV_BINARY = 3
	EVENT_CLOSE       = 4
	EVENT_PING        = 5
	EVENT_PONG        = 6
)

func (self *TNet) WebSocketSendText(Msg string) {
	sByte := self.wsBuildFrame(0x1, 0x1, []byte(Msg))
	self.SendPacket(sByte)
}

func (self *TNet) WebSocketSendClose() {
	sByte := self.wsBuildFrame(0x1, 0x8, []byte(""))
	self.SendPacket(sByte)
}

func (self *TNet) WebSocketSendPing(Msg []byte) {
	sByte := self.wsBuildFrame(0x1, 0x9, Msg)
	self.SendPacket(sByte)
}

func (self *TNet) WebSocketSendPong(Msg []byte) {
	sByte := self.wsBuildFrame(0x1, 0xA, Msg)
	self.SendPacket(sByte)
}

func (self *TNet) wsBuildFrame(FIN byte, opCode byte, data []byte) []byte {
	var (
		WF1      byte = 0
		WF2      byte = 0
		dataSize int  = len(data)
		sByte    []byte
		Mask     uint32 = 0
	)
	if FIN == 1 {
		WF1 = WF1 | 0x80   // WF_FIN
		WF1 = WF1 | opCode // 帧
		sByte = append(sByte, WF1)
	}
	WF2 = WF2 | 0x80 // MASK
	if dataSize > 126 {
		if dataSize > 65535 {
			WF2 = WF2 | 127
			temp := make([]byte, 8)
			binary.BigEndian.PutUint64(sByte, uint64(dataSize))
			sByte = append(sByte, WF2)
			sByte = append(sByte, temp...)
		} else {
			WF2 = WF2 | 126
			temp := make([]byte, 2)
			binary.BigEndian.PutUint16(sByte, uint16(dataSize))
			sByte = append(sByte, WF2)
			sByte = append(sByte, temp...)
		}
	} else {
		WF2 = WF2 | byte(dataSize)
		sByte = append(sByte, WF2)
	}

	MaskKey := make([]byte, 4)
	binary.BigEndian.PutUint32(MaskKey, Mask)
	sByte = append(sByte, MaskKey...)

	sByte = append(sByte, data...)
	//fmt.Println(hex.EncodeToString(sByte))
	return sByte
}

func (self *TIOBuffer) WebSocketReadToCBHook(Event int, fn CBHOOK) {
	if self.ioObject == nil {
		return
	}
	for {
		if self.ioClosed || self.ioObjError.IsError() {
			return
		}
		self.FlushBuffer()
		var (
			WF1            byte   = self.ReadByte()
			WF2            byte   = self.ReadByte()
			WF_FIN                = (WF1 & 0x80) >> 7
			WF_RSV1               = (WF1 & 0x40) >> 6
			WF_RSV2               = (WF1 & 0x20) >> 5
			WF_RSV3               = (WF1 & 0x10) >> 4
			WF_OPCODE             = WF1 & 0xF
			WF_MASK               = (WF2 & 0x80) >> 7
			WF_PAYLOAD_LEN        = WF2 & 0x7F
			PlayLoadLen    uint64 = 0
			WF_MASK_KEY    []byte
			sByte          []byte
			LastByte       []byte
		)
		if 1 == 2 {
			fmt.Println(WF_FIN, WF_RSV1, WF_RSV2, WF_RSV3, WF_MASK_KEY)
		}
		if WF_PAYLOAD_LEN < 126 {
			PlayLoadLen = uint64(WF_PAYLOAD_LEN)
		}
		if WF_PAYLOAD_LEN == 126 {
			sByte = []byte{self.ReadByte(), self.ReadByte()}
			PlayLoadLen = uint64(binary.BigEndian.Uint16(sByte))
		}
		if WF_PAYLOAD_LEN == 127 {
			sByte = []byte{self.ReadByte(), self.ReadByte(), self.ReadByte(), self.ReadByte()}
			sByte = append(sByte, self.ReadByte(), self.ReadByte(), self.ReadByte(), self.ReadByte())
			PlayLoadLen = binary.BigEndian.Uint64(sByte)
		}
		if WF_MASK == 1 {
			sByte := []byte{self.ReadByte(), self.ReadByte(), self.ReadByte(), self.ReadByte()}
			WF_MASK_KEY = sByte
		}
		sByte = self.ReadBytes(int(PlayLoadLen))
		if WF_MASK == 1 {
			sByte = wsMask(WF_MASK_KEY, sByte)
		}
		fmt.Println("wsRecv: ", WF_FIN, WF_RSV1, WF_RSV2, WF_RSV3, WF_MASK, WF_MASK_KEY, WF_OPCODE, PlayLoadLen)
		fmt.Println(string(sByte))
		fmt.Println(hex.EncodeToString(sByte))
		switch WF_OPCODE {
		case 0x0: // %x0 表示一个持续帧
			LastByte = append(LastByte, sByte...)
			fmt.Println(string(LastByte))
		case 0x1: //​ %x1 表示一个文本帧
			fn(EVENT_RECV_TEXT, nil, sByte)
		case 0x2: //​ ​%x2 表示一个二进制帧
			fn(EVENT_RECV_BINARY, nil, sByte)
		case 0x3: //​ %x3-7 预留给以后的非控制帧
		case 0x8: //​ %x8 表示一个连接关闭包
			fn(EVENT_CLOSE, nil, sByte)
			return
		case 0x9: //​ %x9 表示一个ping包
			fn(EVENT_PING, nil, sByte)
		case 0xA: //​ %xA 表示一个pong包
			fn(EVENT_PONG, nil, sByte)
		case 0xB: //​ %xB-F 预留给以后的控制帧

		}

	}
}

func wsMask(MaskKey []byte, sByte []byte) []byte {
	var (
		bSize int = len(sByte)
		i     int = 0
	)
	fmt.Println("wsMask: ", hex.EncodeToString(sByte))
	for i = 0; i < bSize; i++ {
		sByte[i] ^= MaskKey[i%4]
	}
	fmt.Println("wsMask: ", hex.EncodeToString(sByte))
	return sByte
}
