`syscall` with syscall lib

```go
package main

import (
	"log"
	"syscall"
	"net"
)

const (
	host = "127.0.0.1"
	port = 8888
	message = "HTTP/1.1 200 OK\r\n"  +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"Content-Length: 25\r\n" +
		"\r\n" +
		"Server with syscall"
)

func startServer(host string, port int) (int, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if  err != nil {
		log.Fatal("error socket()", err)

	}

	srv := &syscall.SockaddrInet4{Port: port}
	addrs, err := net.LookupHost(host)
	if err != nil {
		log.Fatal("error lookup()", err)
	}

	for _, addr := range addrs {
		ip := net.ParseIP(addr).To4()
		copy(srv.Addr[:], ip)
		if err = syscall.Bind(fd, srv); err != nil {
			log.Fatal("error bind()", err)
		}
	}

	if err = syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		log.Fatal("error listend()", err)
	} else {
		log.Println("listen() on ", host, ":", port)
	}

	if err != nil {
		log.Fatal("error start seerver", err)
	}
	return fd, nil
}

func main() {
	fd, err := startServer(host, port)
	if err != nil {
		log.Fatal("error startServer()", err)
	}

	for {
		clientSock, clientAddr, err  := syscall.Accept(fd)
		// thing like &{37884 [127 0 0 1] {0 0 [0 0 0 0] [0 0 0 0 0 0 0 0]}} how to parse it?
		log.Println("conn from", clientAddr)

		if err != nil {
			log.Fatal("error accept()", err)
		}

		go func(clientSocket int, clientAddress syscall.Sockaddr) {
			err := syscall.Sendmsg(clientSocket, []byte(message), []byte{}, clientAddress, 0)

			if err != nil {
				log.Fatal("error send()", err)
			}
		}(clientSock, clientAddr)

	}

}

```

with `strace` and `nc localhost 8888`
```bash

>>> strace --trace=network ./main                                                                                                                 19:42.13 Mon Mar 20 2023 >>> 
socket(AF_INET, SOCK_STREAM, IPPROTO_IP) = 3
bind(3, {sa_family=AF_INET, sin_port=htons(8888), sin_addr=inet_addr("127.0.0.1")}, 16) = 0
listen(3, 128)                          = 0
2023/03/20 19:42:15 listen() on  127.0.0.1 : 8888
accept4(3, {sa_family=AF_INET, sin_port=htons(55576), sin_addr=inet_addr("127.0.0.1")}, [112 => 16], 0) = 4
2023/03/20 19:42:25 conn from &{55576 [127 0 0 1] {0 0 [0 0 0 0] [0 0 0 0 0 0 0 0]}}
accept4(3, 


```

