package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var Verbose bool
var Server bool
var To string
var BitwiseMode bool

const VERSION = "1.1.0 / Build 2"

func main() {
	var port, interfaceAddress string
	{ // Parse arguments
		flag.BoolVar(&Server, "server", false, "Pass this argument to run as server application")
		flag.BoolVar(&BitwiseMode, "bitwise", false, "Pass this argument to enable bitwise mode; Otherwise addition mode is used)")
		flag.BoolVar(&Verbose, "verbose", false, "More logs")
		showVersion := flag.Bool("version", false, "Show version")
		flag.StringVar(&port, "port", "", "If this is server, the port that proxy listens on it; It this is client, the port that accepts the data")
		flag.StringVar(&interfaceAddress, "interface", "", "Binding address. Server's default is 0.0.0.0 and client's is localhost")
		flag.StringVar(&To, "to", "", "If this is server, the address that the data will be forwarded; If this is client, the server address")
		flag.Parse()
		if *showVersion {
			fmt.Println("Shift Forward Version", VERSION)
			fmt.Println("Source https://github.com/HirbodBehnam/Shift-Forwarder")
			os.Exit(0)
		}
		if port == "" || To == "" {
			fmt.Println("Please enter `port` and `to` values as argument. Pass --help to see help")
			os.Exit(0)
		}
		if interfaceAddress == "" && !Server {
			interfaceAddress = "localhost"
		}
	}
	if Verbose {
		fmt.Println("Verbose mode on")
		fmt.Println("Server mode:", Server)
		fmt.Println("Bitwise mode:", BitwiseMode)
		fmt.Println("Listening on " + interfaceAddress + ":" + port)
	}
	ln, err := net.Listen("tcp", interfaceAddress+":"+port) // start listening for connections
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept() // accept incoming connections
		if err != nil {
			log.Println("Error on accepting new connection:", err)
			continue
		}
		if Verbose {
			log.Println("Accepting new connection from", conn.RemoteAddr())
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	proxy, err := net.Dial("tcp", To) // dial the other side; If this is server, dial the address that the traffic is going to be forwarded. If this is client, the server will be dialed
	if err != nil {
		log.Println("Error on dialing:", err)
		_ = conn.Close()
		return
	}

	go copyIO(conn, proxy)
	go copyIO(proxy, conn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	var err error
	if BitwiseMode {
		err = CopyBitwise(src, dest)
	} else {
		if Server { // TBH, server -> client connection is data[i]-- and not data[i]++ :) (Who cares?)
			err = ServerCopyAddition(src, dest)
		} else {
			err = ClientCopyAddition(src, dest)
		}
	}
	if Verbose {
		if err != nil {
			log.Println("Error on forward:", err) // this will actually throw errors when the copying is done :| (Use of closed connection)
		}
	}
}

func ClientCopyAddition(src, dst net.Conn) (err error) {
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			for i := 0; i < nr; i++ {
				buf[i]++
			}
			nw, ew := dst.Write(buf[0:nr])
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
	return err
}

func ServerCopyAddition(src, dst net.Conn) (err error) {
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			for i := 0; i < nr; i++ {
				buf[i]--
			}
			nw, ew := dst.Write(buf[0:nr])
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
	return err
}

func CopyBitwise(src, dst net.Conn) (err error) {
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			for i := 0; i < nr; i++ {
				buf[i] = ^buf[i]
			}
			nw, ew := dst.Write(buf[0:nr])
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
	return err
}
