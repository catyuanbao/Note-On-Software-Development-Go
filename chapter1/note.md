`syscall ` with C:

```c
#include<stdio.h>
#include<string.h>
#include<sys/socket.h>
#include<arpa/inet.h>
#include<netdb.h>

int main(int argc, char * argv[]) {
  int socket_desc;
  struct sockaddr_in server;
  char * message, server_reply[2000];
  struct hostent * host;
  const char * hostname = "httpbin.org";
  //Create socket
  socket_desc = socket(AF_INET, SOCK_STREAM, 0);
  if (socket_desc == -1) {
    printf("Could not create socket");
  }

  if ((server.sin_addr.s_addr = inet_addr(hostname)) == 0xffffffff) {
    if ((host = gethostbyname(hostname)) == NULL) {
      return -1;
    }

    memcpy(&server.sin_addr, host -> h_addr, host -> h_length);
  }

  // server.sin_addr.s_addr = inet_addr("54.221.78.73");
  server.sin_family = AF_INET;
  server.sin_port = htons(80);

  //Connect to remote server
  if (connect(socket_desc, (struct sockaddr*)&server, sizeof(server)) < 0) {
    puts("connect error");
    return 1;
  }
  puts("Connected\n");
  //Send some data
  message = "GET / HTTP/1.0\n\n";
  if (send(socket_desc, message, strlen(message), 0) < 0) {
    puts("Send failed");
    return 1;
  }
  puts("Data Send\n");
  //Receive a reply from the server
  if (recv(socket_desc, server_reply, 2000, 0) < 0) {
    puts("recv failed");
  }
  puts("Reply received\n");
  puts(server_reply);
  return 0;
}
```

try with `strace`:

```bash
[root@localhost tmp]# strace --trace=network ./a.out 
socket(AF_INET, SOCK_STREAM, IPPROTO_IP) = 3
socket(AF_UNIX, SOCK_STREAM|SOCK_CLOEXEC|SOCK_NONBLOCK, 0) = 4
connect(4, {sa_family=AF_UNIX, sun_path="/var/run/nscd/socket"}, 110) = -1 ENOENT (No such file or directory)
socket(AF_UNIX, SOCK_STREAM|SOCK_CLOEXEC|SOCK_NONBLOCK, 0) = 4
connect(4, {sa_family=AF_UNIX, sun_path="/var/run/nscd/socket"}, 110) = -1 ENOENT (No such file or directory)
socket(AF_INET, SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, IPPROTO_IP) = 4
setsockopt(4, SOL_IP, IP_RECVERR, [1], 4) = 0
connect(4, {sa_family=AF_INET, sin_port=htons(53), sin_addr=inet_addr("100.125.21.250")}, 16) = 0
sendto(4, "8\306\1\0\0\1\0\0\0\0\0\0\7httpbin\3org\0\0\1\0\1", 29, MSG_NOSIGNAL, NULL, 0) = 29
recvfrom(4, "8\306\201\200\0\1\0\4\0\0\0\0\7httpbin\3org\0\0\1\0\1\300\f\0"..., 1024, 0, {sa_family=AF_INET, sin_port=htons(53), sin_addr=inet_addr("100.125.21.250")}, [28->16]) = 93
connect(3, {sa_family=AF_INET, sin_port=htons(80), sin_addr=inet_addr("107.22.139.22")}, 16) = 0
Connected

sendto(3, "GET / HTTP/1.0\n\n", 16, 0, NULL, 0) = 16
Data Send

recvfrom(3, "HTTP/1.1 200 OK\r\nDate: Sun, 19 M"..., 2000, 0, NULL, NULL) = 2000
Reply received

HTTP/1.1 200 OK
Date: Sun, 19 Mar 2023 13:12:02 GMT
Content-Type: text/html; charset=utf-8
Content-Length: 9593
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>httpbin.org</title>
    <link href="https://fonts.googleapis.com/css?family=Open+Sans:400,700|Source+Code+Pro:300,600|Titillium+Web:400,600,700"
        rel="stylesheet">
    <link rel="stylesheet" type="text/css" href="/flasgger_static/swagger-ui.css">
    <link rel="icon" type="image/png" href="/static/favicon.ico" sizes="64x64 32x32 16x16" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }

        *,
        *:before,
        *:after {
            box-sizing: inherit;
        }

        body {
            margin: 0;
            background: #fafafa;
        }
    </style>
</head>

<body>
    <a href="https://github.com/requests/httpbin" class="github-corner" aria-label="View source on Github">
        <svg width="80" height="80" viewBox="0 0 250 250" style="fill:#151513; color:#fff; position: absolute; top: 0; border: 0; right: 0;"
            aria-hidden="true">
            <path d="M0,0 L115,115 L130,115 L142,142 L250,250 L250,0 Z"></path>
            <path d="M128.3,109.0 C113.8,99.7 119.0,89.6 119.0,89.6 C122.0,82.7 120.5,78.6 120.5,78.6 C119.2,72.0 123.4,76.3 123.4,76.3 C127.3,80.9 125.5,87.3 125.5,87.3 C122.9,97.6 130.6,101.9 134.4,103.2"
                fill="currentColor" style="transform-origin: 130px 106px;" class="octo-arm"></path>
            <path d="M115.0,115.0 C114.9,115.1 118.7,116.5 119.8,115.4 L133.7,101.6 C136.9,99.2 139.9,98.4 142.2,98.6 C133.8,88.0 127.5,74.4 143.8,58.0 C148.5,53.4 154.0,51.2 159.7,51.0 C160.3,49.4 163.2,43.6 171.4,40.1 C171.4,40.1 176.1,42.5 178.8,56.2 C183.1,58.6 187.2,61.8 190.9,65.4 C194.5,69.0 197.7,73.2 200.1
+++ exited with 0 +++

```



`syscall` in Go, need `https://github.com/golang/sys` in $GROOT

```go
package main

import (
	syscall "golang/org/x/sys/unix"
	"log"
)
// for more info: https://pkg.go.dev/golang.org/x/sys/unix
func main() {
	c := make([]byte, 512)

	log.Println("Getpid:", syscall.Getpid())
	log.Println("Getpgrp:", syscall.Getpgrp())
	log.Println("Getppid: ", syscall.Getppid())
	log.Println("Gettid: ", syscall.Gettid())

	_, err := syscall.Getcwd(c)

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("env:", syscall.Environ())

	log.Println(string(c))
}


```



get disk use pct:

```go
package main

import (
	"fmt"
	syscall "golang.org/x/sys/unix"
	"os"
)

const (
	gigabyte = (1024.0 * 1024.0 * 1024.0)
)

func main() {
	var statfs = syscall.Statfs_t{}
	var total uint64
	var used uint64
	var free uint64
	err := syscall.Statfs("/", &statfs)
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err)
		os.Exit(1)
	} else {
		total = statfs.Blocks * uint64(statfs.Bsize)
		free = statfs.Bfree * uint64(statfs.Bsize)
		used = total - free
	}

	fmt.Printf("total Disk Space : %.1f GB\n", float64(total)/gigabyte)
	fmt.Printf("total Disk used  : %.1f GB\n", float64(used)/gigabyte)
	fmt.Printf("total Disk free  : %.1f GB\n", float64(free)/gigabyte)
	fmt.Printf("total Disk used pct  : %.1f %% \n", (float64(used)/float64(total))*100)
}


```



