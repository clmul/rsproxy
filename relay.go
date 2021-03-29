package rsproxy

import (
	"io"
	"log"
	"net"
)

func Relay(remote, local net.Conn) {
	done := make(chan int)
	go relay(remote, local, done)
	relay(local, remote, nil)
	<-done
	remote.Close()
	local.Close()
}

func relay(dst, src net.Conn, done chan int) {
	buffer := make([]byte, 256)
	if done != nil {
		defer func() {
			done <- 1
		}()
	}

	for {
		n, err := src.Read(buffer)
		// read may return EOF with n > 0
		// should always process n > 0 bytes before handling error
		if n > 0 {
			_, err1 := dst.Write(buffer[:n])
			if err1 != nil {
				log.Println(err1)
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("read error:", err)
			break
		}
	}
}
