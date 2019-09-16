#Shift Forwarder
A simple tool to shift all bytes and forward traffic
## Install
Use the executables from releases and download one for your server and client.
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
## Why this?
I personally use this with an HTTP Proxy server to bypass censorship.