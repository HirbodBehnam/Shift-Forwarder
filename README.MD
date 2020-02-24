# Shift Forwarder
A simple tool to shift all bytes and forward traffic
## Why this?
I personally use this with an HTTP Proxy server to bypass censorship. Iran's DPI system is stupid enough to let the traffic get through. You can also use a SOCKS server behind the forwarder.

No shit guys, I almost went WTF when this worked.
## Install
Use the executables from [releases](https://github.com/HirbodBehnam/Shift-Forwarder/releases) and download one for your server and client.
### Server
```
./sf -server -to 127.0.0.1:8888 -interface 127.0.0.1 -port 1080
```
Run this command. It will get all traffics from port 1080, shift the bytes back and forwards them to 127.0.0.1:8888.
### Client
```
 ./sf -to 1.1.1.1:1080 -port 8080
```
This will get all incoming traffics from port 8080 and forwards them to 1.1.1.1:1080.
### Building
Just as easy as
```bash
go build main.go
```
## Arguments
```
Usage of sf.exe:
  -bitwise
        Pass this argument to enable bitwise mode; Otherwise addition mode is used
  -interface string
        Binding address. Server's default is 0.0.0.0 and client's is localhost
  -port string
        If this is server, the port that proxy listens on it; It this is client, the port that accepts the data
  -server
        Pass this argument to run as server application
  -to string
        If this is server, the address that the data will be forwarded; If this is client, the server address
  -verbose
        More logs
  -version
        Show version
```
### Bitwise Mode
Bitwise mode, bitwises all bytes before sending them to server or client. So server -> client and client -> server connections are bitwised.

If you enable bitwise mode on server, it must be enabled on server too.
### How non bitwise mode works
When client forwards something to server, it adds 1 to all of the bytes. On server, the all the bytes will be subtracted by one. After the subtract, server will forward the traffic to the destination. The server->client traffic is also shifted like this.
