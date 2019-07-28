package Proxy

import (
	"math/rand"
	"net"
	"strconv"
	"time"
)

func getUnUsedPort() int {
	port := 20000
	for {
		seed := rand.New(rand.NewSource(time.Now().UnixNano()))
		randNum := seed.Intn(10000)
		port += randNum
		//fmt.Printf("used port num is:%d", port)

		dwRet := func() bool {
			defer func() {
				if ok := recover(); ok != nil {
					//fmt.Printf(" -> %s\n", "fail")
					port = 20000
				}
			}()
			conn, err := net.Listen("tcp", ":"+strconv.Itoa(port))
			if err != nil {
				panic("unavailable")
				return false
			} else {
				conn.Close()
				return true
			}
		}()
		if dwRet == true {
			//fmt.Printf("\n%d\n", port)
			return port
			break
		}
	}
	return 0
}
