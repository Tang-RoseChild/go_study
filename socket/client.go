package main

import (
	"socket/protocol"

	"flag"
	"fmt"
	"net"
	"os"

	"strconv"
	"time"
)

func main() {
	addr := flag.String("http", ":20003", "host address : `-http=:20003`")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)

	CheckErr(err)
	var count int
	defer conn.Close()
	for i := 0; i < 10; i++ {
		fs := `{"ID":"%s","Session":"%s","Meta":"golang","Content":"message"}`
		msg := fmt.Sprintf(fs, strconv.Itoa(i), getSession())
		count += len(msg)
		conn.Write(protocol.Pack([]byte(msg)))
		CheckErr(err)
		fmt.Println("msg : ", msg)
		// time.Sleep(1 * time.Second)
	}
	fmt.Println("total count is : ", count)
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("error : ", err.Error())
		os.Exit(1)

	}
}

func getSession() string {
	s := time.Now().Unix()
	return strconv.FormatInt(s, 10)
}
