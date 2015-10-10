package main

import (
	"anothersocket/protocol"
	"fmt"
	"net"

	"time"
)

// var msgs chan []byte

// func init() {
// 	msgs = make(chan []byte, 30)
// }

func main() {
	l, err := net.Listen("tcp", ":10005")
	CheckErr(err)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("err when accept : ", err.Error())
			continue
		}
		go handleConn(conn)
	}
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("err : ", err.Error())
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		fmt.Println("closed connect")
		conn.Close()

	}()
	var msgs chan []byte = make(chan []byte, 30)
	var q bool = true
	// start to receive msg
	go receive(conn, msgs, &q)

	buf := make([]byte, 1<<10)
	tmpBuf := make([]byte, 0)

	for q {
		fmt.Println("reading")
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := conn.Read(buf) // no data, this will block until some data comming
		// fmt.Printf("buf :%o %o ", buf[0], buf[1])
		if err != nil {
			if operr, ok := err.(*net.OpError); ok {
				if operr.Timeout() {
					continue
				}
			}
			return
		}
		tmpBuf = protocol.Unpack(append(tmpBuf, buf[:n]...), msgs)
	}

	fmt.Println("return handleConn")
}

func receive(conn net.Conn, msgs chan []byte, q *bool) {
	var count = 3
	for {
		select {
		case msg := <-msgs:
			// fmt.Println("msg from client ", string(msg))
			fmt.Println(string(msg))
			count = 3
		default:
			fmt.Println("send heart", count)
			conn.Write(protocol.Pack([]byte("HEART")))

			time.Sleep(5 * time.Second)
			count--
			if count <= 0 {

				*q = false //quit connection

				return
			}
		}
	}
}
