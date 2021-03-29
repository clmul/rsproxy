package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/clmul/socks5"
)

var (
	socks5proxy = flag.String("socks5", "", "SOCKS5 proxy listen address")
	listen      = flag.String("listen", "", "client listen address")
)

func main() {
	flag.Parse()

	connCh := make(chan net.Conn, 10)
	go func() {
		ln, err := net.Listen("tcp", *listen)
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatal(err)
			}
			connCh <- conn
		}
	}()
	socks5Server(connCh)
}

func socks5Server(connCh chan net.Conn) {
	dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn := <-connCh
		buf := make([]byte, 0, len(addr)+1)
		buf = append(buf, byte(len(addr)))
		buf = append(buf, addr...)
		_, err := conn.Write(buf)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return conn, nil
	}

	server, err := socks5.New(&socks5.Config{Dial: dial})
	if err != nil {
		log.Println(err)
		return
	}
	log.Fatal(server.ListenAndServe("tcp", *socks5proxy))
}
