package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var forwarder = flag.String("forwarder", "", "")

func main() {
	flag.Parse()
	go worker()
	go worker()
	go worker()
	worker()
}

func worker() {
	ch := make(chan int)
	for {
		go gateway(ch)
		<-ch
	}
}

func gateway(ch chan int) {
	conn, err := net.Dial("tcp", *forwarder)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Println("dial done")
	var buf [256]byte
	_, err = io.ReadFull(conn, buf[:1])
	if err != nil {
		log.Println(err)
		return
	}
	ch <- 1
	size := buf[0]
	log.Println("size", size)
	_, err = io.ReadFull(conn, buf[:size])
	if err != nil {
		log.Println(err)
		return
	}
	host := string(buf[:size])
	log.Println("dial", host)
	upstream, err := net.Dial("tcp", host)
	if err != nil {
		log.Println(err)
		return
	}
	defer upstream.Close()
	wait := make(chan int)
	go func() {
		io.Copy(upstream, conn)
		wait <- 1
	}()
	io.Copy(conn, upstream)
	<-wait
}
