package main

import (
	"anothersocket/protocol"
	"fmt"
	"net"
	"sync"

	"bufio"
	"io"
	"os"
	"strconv"
	"time"
)

var wg sync.WaitGroup
var msgs chan []byte

func main() {
	msgs = make(chan []byte, 10)
	wg = sync.WaitGroup{}
	wg.Add(1)
	conn, err := net.Dial("tcp", ":10005")
	CheckErr(err)
	go send(conn)

	var q bool = true
	go receive(conn, &q)
	b := make([]byte, 1<<10)
	for q {
		n, err := conn.Read(b)
		CheckErr(err)
		if err == io.EOF {
			fmt.Println("server closed : ")
			return
		}
		protocol.Unpack(b[:n], msgs)
	}

	wg.Wait()

}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("err : ", err.Error())
	}
}

func send(conn net.Conn) {
	defer func() {
		wg.Done()
		// conn.Close()
	}()

	fs := `{"id":"%d","Session":"%s","Content":"Testing Msg %d"}`
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf(fs, i, getSession(i), i)

		n, err := conn.Write(protocol.Pack([]byte(msg)))
		fmt.Println(n, ": msg : ", msg)
		// if n != len(msg) {
		// 	fmt.Println("not write all ")
		// 	return
		// }
		if err != nil {
			fmt.Println("err ", err.Error())
			return
		}
	}
	bf := bufio.NewReader(os.Stdin)
	for {
		b, _ := bf.ReadSlice('\n')
		conn.Write(b)
	}

}

func getSession(id int) string {
	t := time.Now().Unix()
	// time.Sleep(5 * time.Millisecond)
	return strconv.FormatInt(t, 10)
}

func receive(conn net.Conn, q *bool) {
	for {
		select {
		case msg := <-msgs:
			// fmt.Println("msg from client ", string(msg))
			if string(msg) == "HEART" {
				// _, err := conn.Write(protocol.Pack([]byte("HEART")))
				// if err != nil {
				// 	*q = false
				// 	return
				// }
				continue
			}
			fmt.Println(string(msg))
		}
	}
}
