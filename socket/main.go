package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"socket/protocol"
)

func main() {
	addr := flag.String("http", ":20003", "host address : `-http=:20003`")
	flag.Parse()
	l, err := net.Listen("tcp", *addr)
	CheckErr(err)
	msgs := make(chan []byte, 10)
	go reader(msgs)
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn, msgs)
	}

}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("error : ", err.Error())
		os.Exit(1)

	}
}

func reader(msgs chan []byte) {
	fmt.Println("reader : ")
	for {
		select {
		case m := <-msgs:
			fmt.Println("msg : ", string(m))
		}
	}
}

func handleConn(conn net.Conn, msgs chan []byte) {
	defer conn.Close()
	// fmt.Println("start to read")
	tmpBuf := make([]byte, 0)
	b := make([]byte, 1<<10)
	// fmt.Println("b after init : ", b)
	for {
		n, err := conn.Read(b)
		// fmt.Println("n : ", n)
		if err != nil {
			if err == io.EOF {
				fmt.Println("client closed")
				return
			}
			fmt.Println("read error : ", err.Error())
			return
		}
		tmpBuf = protocol.Unpack(append(tmpBuf, b[:n]...), msgs)
		// fmt.Println("temp buf :: ", string(tmpBuf))
		// fmt.Println("msg is : ", string(msg))
		// msgs <- msg

	}
}
