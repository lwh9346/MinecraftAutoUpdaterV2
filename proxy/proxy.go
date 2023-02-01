package proxy

import (
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/golang/snappy"
)

// SetUp 建立一个代理
func SetUp(listenAddr, targetAddr string, isServer bool) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handleConn(conn, targetAddr, isServer)
	}
}

func handleConn(source net.Conn, targetAddr string, isServer bool) {
	log.Printf("new connection %s->%s", source.RemoteAddr(), source.LocalAddr())
	defer source.Close()
	target, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Println(err)
		return
	}
	var compressed, uncompressed net.Conn
	if isServer {
		compressed = source
		uncompressed = target
	} else {
		compressed = target
		uncompressed = source
	}
	reader := snappy.NewReader(compressed)
	writer := snappy.NewBufferedWriter(compressed)
	go io.Copy(uncompressed, reader)
	copyWithTimeout(writer, uncompressed, 50*time.Millisecond)
}

func copyWithTimeout(dst *snappy.Writer, src net.Conn, timeout time.Duration) {
	defer dst.Close()
	buf := make([]byte, 65536)
	nextFlush := time.Now().Add(timeout)
	for {
		src.SetReadDeadline(nextFlush)
		n, err := src.Read(buf)
		now := time.Now()
		if err != nil {
			e, ok := err.(net.Error)
			if !ok || !e.Timeout() {
				if !errors.Is(err, io.EOF) {
					log.Println(err)
				}
				return
			}
		}
		_, err = dst.Write(buf[:n])
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println(err)
			}
			return
		}
		if now.After(nextFlush) {
			dst.Flush()
			nextFlush = now.Add(timeout)
		}
	}
}
