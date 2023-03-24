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
		log.Fatal("error start server", err)
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


## read ELF info

```bash
>> readelf -hd  ./a.out                                                                                                                         18:45.20 Fri Mar 24 2023 >>> 
ELF Header:
  Magic:   7f 45 4c 46 02 01 01 00 00 00 00 00 00 00 00 00 
  Class:                             ELF64
  Data:                              2's complement, little endian
  Version:                           1 (current)
  OS/ABI:                            UNIX - System V
  ABI Version:                       0
  Type:                              EXEC (Executable file)
  Machine:                           Advanced Micro Devices X86-64
  Version:                           0x1
  Entry point address:               0x401060
  Start of program headers:          64 (bytes into file)
  Start of section headers:          23936 (bytes into file)
  Flags:                             0x0
  Size of this header:               64 (bytes)
  Size of program headers:           56 (bytes)
  Number of program headers:         13
  Size of section headers:           64 (bytes)
  Number of section headers:         31
  Section header string table index: 30

Dynamic section at offset 0x2e00 contains 25 entries:
  Tag        Type                         Name/Value
 0x0000000000000001 (NEEDED)             Shared library: [libpcap.so.1]
 0x0000000000000001 (NEEDED)             Shared library: [libc.so.6]
 0x000000000000000c (INIT)               0x401000
 0x000000000000000d (FINI)               0x4011bc
 0x0000000000000019 (INIT_ARRAY)         0x403df0
 0x000000000000001b (INIT_ARRAYSZ)       8 (bytes)
 0x000000000000001a (FINI_ARRAY)         0x403df8
 0x000000000000001c (FINI_ARRAYSZ)       8 (bytes)
 0x000000006ffffef5 (GNU_HASH)           0x4003a0
 0x0000000000000005 (STRTAB)             0x4004a0
 0x0000000000000006 (SYMTAB)             0x4003c8
 0x000000000000000a (STRSZ)              164 (bytes)
 0x000000000000000b (SYMENT)             24 (bytes)
 0x0000000000000015 (DEBUG)              0x0
 0x0000000000000003 (PLTGOT)             0x404000
 0x0000000000000002 (PLTRELSZ)           72 (bytes)
 0x0000000000000014 (PLTREL)             RELA
 0x0000000000000017 (JMPREL)             0x400600
 0x0000000000000007 (RELA)               0x400588
 0x0000000000000008 (RELASZ)             120 (bytes)
 0x0000000000000009 (RELAENT)            24 (bytes)
 0x000000006ffffffe (VERNEED)            0x400558
 0x000000006fffffff (VERNEEDNUM)         1
 0x000000006ffffff0 (VERSYM)             0x400544
 0x0000000000000000 (NULL)               0x0
<<< 
```

### show info with `goplay`

code of `goplay`

```golang
package main

import (
        "debug/elf"
        "flag"
        "fmt"
        "os"
        "strings"
        "syscall"
)

func dump_dynamic_str(file *elf.File) {
        fmt.Printf("dynamic libarry Strings:\n")
        dynamic_strs, _ := file.DynString(elf.DT_NEEDED)

        for _, e := range dynamic_strs {
                fmt.Printf("\t%s\n", e)
        }

        dynamic_strs, _ = file.DynString(elf.DT_SONAME)
        for _, e := range dynamic_strs {
                fmt.Printf("\t%s\n", e)
        }

        dynamic_strs, _ = file.DynString(elf.DT_RPATH)
        for _, e := range dynamic_strs {
                fmt.Printf("\t%s\n", e)
        }

        dynamic_strs, _ = file.DynString(elf.DT_RUNPATH)
        for _, e := range dynamic_strs {
                fmt.Printf("\t%s\n", e)
        }
}

func dump_symbols(file *elf.File) {
        fmt.Printf("Symbols:\n")
        symbols, _ := file.Symbols()
        for _, e := range symbols {
                if !strings.EqualFold(e.Name, "") {
                        fmt.Printf("\t%s\n", e.Name)
                }
        }
}

func dump_elf(filename string) int {
        file, err := elf.Open(filename)
        if err != nil {
                fmt.Printf("Couldn’t open file : \"%s\" as an ELF.\n")
                return 2
        }
        dump_dynamic_str(file)
        dump_symbols(file)
        return 0
}

func init_debug(filename string) int {
        attr := &os.ProcAttr{Sys: &syscall.SysProcAttr{Ptrace: true}}
        if proc, err := os.StartProcess(filename, []string{"/"}, attr); err == nil {
                _, w_err := proc.Wait()
                if w_err != nil {
                        fmt.Printf("Wait process error")
                        return 2
                }
                foo := syscall.PtraceAttach(proc.Pid)
                fmt.Printf("Started New Process: %v.\n", proc.Pid)
                fmt.Printf("PtraceAttach res: %v.\n", foo)
                return 0
        }
        return 2
}

func main() {
        len_args := len(os.Args)
        if len_args > 1 {
                filename := flag.String("filename", "", "A binary ELF file.")
                action := flag.String("action", "", "Action to make: {dump|debug}.")
                flag.Parse()
                if *filename == "" || *action == "" {
                        goto Error
                }

                file, err := os.Stat(*filename)
                if os.IsNotExist(err) {
                        fmt.Printf("No such file or directory: %s.\n", *filename)
                        goto Error
                } else if mode := file.Mode(); mode.IsDir() {
                        fmt.Printf("Parameter must be a file, not a " +
                                "directory.\n")
                        goto Error
                }
                f, err := os.Open(*filename)
                if err != nil {
                        fmt.Printf("Couldn’t open file: \"%s\".\n", *filename)
                        goto Error
                }
                err = f.Close()
                if err != nil {
                        fmt.Printf("Close file error")
                        goto Error
                }
                fmt.Printf("Tracing program : \"%s\".\n", *filename)
                fmt.Printf("Action : \"%s\".\n", *action)

                switch *action {
                case "debug":
                        os.Exit(init_debug(*filename))
                case "dump":
                        os.Exit(dump_elf(*filename))
                }
        } else {
                goto Usage
        }

Usage:
        fmt.Printf("Usage of ./main:\n" +
                "  -action=\"{dump|debug}\": Action to make.\n" +
                "  -filename=\"file\": A binary ELF file.\n")
        goto Error
Error:
        os.Exit(1)
}

```

try `goplay`:

```bash
>> ./main -filename=a.out -action=dump                                                                                                          19:15.42 Fri Mar 24 2023 >>> 
Tracing program : "a.out".
Action : "dump".
dynamic libarry Strings:
        libpcap.so.1
        libc.so.6
Symbols:
        /usr/lib/gcc/x86_64-redhat-linux/11/../../../../lib64/crt1.o
        .annobin_abi_note.c
        .annobin_abi_note.c_end
        .annobin_abi_note.c.hot
        .annobin_abi_note.c_end.hot
        .annobin_abi_note.c.unlikely
        .annobin_abi_note.c_end.unlikely
        .annobin_abi_note.c.startup
        .annobin_abi_note.c_end.startup
        .annobin_abi_note.c.exit
        .annobin_abi_note.c_end.exit
        __abi_tag
        .annobin_init.c
        .annobin_init.c_end
        .annobin_init.c.hot
        .annobin_init.c_end.hot
        .annobin_init.c.unlikely
        .annobin_init.c_end.unlikely
        .annobin_init.c.startup
        .annobin_init.c_end.startup
        .annobin_init.c.exit
        .annobin_init.c_end.exit
        .annobin_static_reloc.c
        .annobin_static_reloc.c_end
        .annobin_static_reloc.c.hot
        .annobin_static_reloc.c_end.hot
        .annobin_static_reloc.c.unlikely
        .annobin_static_reloc.c_end.unlikely
        .annobin_static_reloc.c.startup
        .annobin_static_reloc.c_end.startup
        .annobin_static_reloc.c.exit
        .annobin_static_reloc.c_end.exit
        .annobin__dl_relocate_static_pie.start
        .annobin__dl_relocate_static_pie.end
        crtstuff.c
        deregister_tm_clones
        register_tm_clones
        __do_global_dtors_aux
        completed.0
        __do_global_dtors_aux_fini_array_entry
        frame_dummy
        __frame_dummy_init_array_entry
        a.c
        crtstuff.c
        __FRAME_END__
        _DYNAMIC
        __GNU_EH_FRAME_HDR
        _GLOBAL_OFFSET_TABLE_
        __libc_start_main@GLIBC_2.34
        _ITM_deregisterTMCloneTable
        data_start
        _edata
        _fini
        printf@GLIBC_2.2.5
        __data_start
        fprintf@GLIBC_2.2.5
        __gmon_start__
        __dso_handle
        _IO_stdin_used
        _end
        _dl_relocate_static_pie
        _start
        pcap_lookupdev
        __bss_start
        main
        __TMC_END__
        _ITM_registerTMCloneTable
        _init
        stderr@GLIBC_2.2.5

```

after `strip` file, no symbols:

```bash
>> strip a.out                                                                                                                                  19:16.31 Fri Mar 24 2023 >>> 
<<< 
<<< root@gcrun~/go/work/temp
>>> ./main -filename=a.out -action=dump                                                                                                          19:16.39 Fri Mar 24 2023 >>> 
Tracing program : "a.out".
Action : "dump".
dynamic libarry Strings:
        libpcap.so.1
        libc.so.6
Symbols:
```
