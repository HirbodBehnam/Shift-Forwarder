# Shift Forwarder
A simple tool to shift all bytes and forward traffic
## Why this?
I personally use this with an HTTP Proxy server to bypass censorship. Iran's DPI system is stupid enough to let the traffic get through. You can also use a SOCKS server behind the forwarder.

No shit guys, I almost went WTF when this worked.
## Install
Use the executables from [releases](https://github.com/HirbodBehnam/Shift-Forwarder/releases) and download one for your server and client.
### Server
```
./sf s 1080 127.0.0.1:8888
```
Run this command. It will get all traffics from port 1080, shift the bytes back and forwards them to 127.0.0.1:8888.
### Client
```
 ./sf c 8080 1.1.1.1:1080
```
This will get all incoming traffics from port 8080 and forwards them to port 1.1.1.1:1080.
### Building
Just as easy as
```bash
go build main.go
```
## How it works
SO SIMPLE! When client forwards something to server, it adds 1 to all of the bytes. On server, the all the bytes will be subtracted by one. After the subtract, server will forward the traffic to the destination. The server->client traffic is also shifted like this.
