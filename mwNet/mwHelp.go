package mwNet

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"
)

func Host2IP(Host string) string {
	str, _ := net.LookupHost(Host)
	if len(str) > 1 {
		return str[0]
	}
	//fmt.Println(err)
	return ""
}

func helpIPHostToByte(Host string) []byte {
	var (
		dstIP []byte = make([]byte, 4)
	)
	binary.BigEndian.PutUint32(dstIP, helpIPToUInt32(net.ParseIP(Host)))
	return dstIP
}

func helpPortToByte(Port string) []byte {
	var (
		dstPort []byte = make([]byte, 2)
		iPort   int
	)
	iPort, _ = strconv.Atoi(Port)
	binary.BigEndian.PutUint16(dstPort, uint16(iPort))
	return dstPort
}

func helpIPToUInt32(ipnr net.IP) uint32 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}
