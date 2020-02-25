package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var (
	BitwiseMode bool
	Server      bool
	Verbose     bool
	To          string
)

const VERSION = "1.2.0 / Build 3"

func main() {
	var port, interfaceAddress string
	if os.Getenv("SS_LOCAL_HOST") != "" { // check if the program is running as shadowsocks plugin
		pluginOptions := strings.Split(os.Getenv("SS_PLUGIN_OPTIONS"), ";")
		Server = ArrayContains(pluginOptions, "server")
		BitwiseMode = ArrayContains(pluginOptions, "bitwise")
		Verbose = ArrayContains(pluginOptions, "verbose")
		if Server { // server mode
			To = os.Getenv("SS_LOCAL_HOST") + ":" + os.Getenv("SS_LOCAL_PORT")
			interfaceAddress = os.Getenv("SS_REMOTE_HOST")
			port = os.Getenv("SS_REMOTE_PORT")
		} else { // client mode
			To = os.Getenv("SS_REMOTE_HOST") + ":" + os.Getenv("SS_REMOTE_PORT")
			interfaceAddress = os.Getenv("SS_LOCAL_HOST")
			port = os.Getenv("SS_LOCAL_PORT")
		}
	} else { // Parse arguments
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
		fmt.Println("Listening on", interfaceAddress+":"+port)
		fmt.Println("Forwarding to", To)
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
	copyIO(proxy, conn)
	if Verbose {
		log.Println("Closing a connection form", conn.RemoteAddr())
	}
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
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") { // when a connection is closed, the other pipe raises use of closed network connection error
		log.Println("Error on forward:", err)
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

// checks if an array contains a specific element
func ArrayContains(ary []string, check string) bool {
	for _, k := range ary {
		if k == check {
			return true
		}
	}
	return false
}
