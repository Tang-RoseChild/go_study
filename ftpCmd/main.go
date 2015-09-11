package main

import (
	"fmt"
	// "log"
	// "bufio"
	"flag"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

// func def() {

// 	fmt.Println(l)

// }

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage : %s -http=addr\n", os.Args[0])
		return
	}
	addr := flag.String("http", ":20003", "host address like 127.0.0.1:20003")

	l, err := net.Listen("tcp", *addr)
	checkErr(err)
	for {
		conn, err := l.Accept()
		checkErr(err)
		go handleConn(conn)
	}
	// test *********
	// out, _ := exec.Command("/bin/bash", "-c", "cd .. ;ls").Output()
	// out, _ := exec.Command("ls").Output()
	// fmt.Println("cmd out : ", string(out))
	// conn, _ := l.Accept()
	// dirResp(conn)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("error : ", err.Error())
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	// loop to read from conn
	var buf []byte = make([]byte, 512)
	for {
		// parse the string
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		checkErr(err)
		// if pwd, write back pwd path
		if string(buf[:3]) == "pwd" {
			pwdResp(conn)

		} else if string(buf[:2]) == "cd" {
			// get the dir path but without newline char
			cdResp(string(buf[2:n-1]), conn)
		} else if string(buf[:3]) == "dir" {
			dirResp(conn)
		}

	}
}

func pwdResp(conn net.Conn) {
	out, err := exec.Command("pwd").Output()
	if err != nil {
		fmt.Println("pwdResp error : ", err.Error())
		conn.Write([]byte("ERROR\r\n"))
	}
	if n, err := conn.Write(out); err != nil || n != len(out) {
		fmt.Printf("error when write pwdResp %d %s \n ", n, err.Error())
	}
}

func dirResp(conn net.Conn) {
	out, err := exec.Command("ls").Output()
	if err != nil {
		fmt.Println("dirResp error : ", err.Error())
		conn.Write([]byte("ERROR\r\n"))
	}
	if n, err := conn.Write(out); err != nil || n != len(out) {
		fmt.Printf("error when write dirResp %d %s \n ", n, err.Error())
	}
	// fmt.Println("output : ", string(out))
}

func cdResp(dir string, conn net.Conn) {
	/*
		throug the following way,the current path dir hasn't changed, so
		when next time call `pwd`, the result is still the formmer one

		// cmd := "cd " + dir + ";ls"
		// fmt.Println("cmd is : ", cmd)
		// out, err := exec.Command("/bin/bash", "-c", cmd).Output()
	*/
	err := os.Chdir(strings.TrimSpace(dir))

	if err != nil {
		fmt.Println("cdResp error : ", err.Error())
		conn.Write([]byte("ERROR\r\n"))
	}
	// if n, err := conn.Write(out); err != nil || n != len(out) {
	// 	fmt.Printf("error when write cdResp %d %s \n ", n, err.Error())
	// }
	dirResp(conn)

}
