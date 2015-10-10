package main

import (
	"net"
	// "ftp"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

// login
// file transfer
// logout

func main() {
	// give certain host and port
	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s -http=`host:port` \n", os.Args[0])
		os.Exit(1)
	}

	addr := flag.String("http", "", "given the certain host and port")
	flag.Parse()
	l, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Println("err when listen ", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("accept error ", err.Error())
			continue
		}
		go handleConn(conn)
	}

}

// ftp protocol
// first byte:
// 1 : dir	2:pwd	3:cd <dir>	4:get <file>
// end line : \r\n\r\n

func handleConn(conn net.Conn) {
	defer conn.Close()
	// read from conn
	bf := bufio.NewReader(conn)
	for {
		b := handleReadErr(bf.ReadSlice)('\n')
		if b == nil {
			break
		}
		fmt.Println("b : ", b)
		handleWriteErr(conn.Write)(b)
	}
	// send
}

type ReadHandler func(byte) []byte

func handleReadErr(ReadSlice func(byte) ([]byte, error)) ReadHandler {
	return func(b byte) []byte {

		bytes, err := ReadSlice(b)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Println("error when ReadSlice: ", err.Error())
			return []byte("Error when Read " + err.Error())
		}
		return bytes
	}

}

type WriteHandler func([]byte)

func handleWriteErr(write func([]byte) (int, error)) WriteHandler {
	return func(bytes []byte) {

		n, err := write(bytes)
		if err != nil {
			fmt.Println("error when write: ", err.Error())
			write([]byte("Error when write " + err.Error()))
		}

		if n != len(bytes) {
			fmt.Println("error when write , not write all data ")
			write([]byte("Error when write Not write all data "))
		}

	}

}
