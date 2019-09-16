package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)
var Server bool
var To string

const VERSION = "1.0.0 / Build 1"
func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h"{
		fmt.Println("Shift Forward Version",VERSION)
		fmt.Println("Use this app like this:")
		fmt.Println("Server side:")
		fmt.Println("./sf s <port> <to>")
		fmt.Println("Example: ./sf s 1080 127.0.0.1:8888	This command accepts connections on port 1080 and forwards them to 127.0.0.1:8888")
		fmt.Println()
		fmt.Println("Client side:")
		fmt.Println("./sf c <port> <to>")
		fmt.Println("Example: ./sf c 8080 1.1.1.1:1080		This command accepts outgoing connections on port 8080 and forwards them to 1.1.1.1:1080")
		os.Exit(0)
	}
	Server = os.Args[1] == "s"
	fmt.Println("Server mode:",Server)
	To = os.Args[3]
	fmt.Println("To:",To)

	ln, err := net.Listen("tcp", ":" + os.Args[2])
	fmt.Println("Listen on " + ":" + os.Args[2])
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error:",err)
			continue
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	proxy, err := net.Dial("tcp", To)
	if err != nil {
		log.Println("Error:",err)
		return
	}

	go copyIO(conn, proxy)
	go copyIO(proxy, conn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	if Server{
		_, _ = copyS(src, dest)
	}else {
		_, _ = copyC(src, dest)
	}
}
func copyC(src,dst net.Conn) (written int64, err error){
	buf := make([]byte, 32 * 1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			for i := 0; i<nr ; i++ {
				buf[i]++
			}
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
func copyS(src,dst net.Conn) (written int64, err error){
	buf := make([]byte, 32 * 1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			for i := 0; i<nr ; i++ {
				buf[i]--
			}
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}