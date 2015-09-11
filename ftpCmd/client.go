package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	addr := flag.String("host", ":20003", "host address like :20003")
	conn, err := net.Dial("tcp", *addr)
	checkErr(err)
	defer conn.Close()
	// read from stdin
	go func(conn net.Conn) {
		bf := bufio.NewReader(os.Stdin)
		for {
			b, err := bf.ReadSlice('\n')
			checkErr(err)
			if string(b[:4]) == "quit" {
				os.Exit(0)
			}
			conn.Write(b)
		}
	}(conn)
	var buf []byte = make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			fmt.Println("server closed")
			return
		}
		fmt.Printf(string(buf[:n]))
	}

}

func checkErr(err error) {
	if err != nil {
		fmt.Println("error : ", err.Error())
	}
}
