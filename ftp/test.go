package main

import (
	"bytes"
	"fmt"
)

func main() {
	bf := bytes.NewBuffer([]byte{3})

	fmt.Println("buf len : ", bf.Len())
	b, err := bf.ReadByte()
	if err != nil {
		fmt.Println("error when read : ", err.Error())
		return
	}

	fmt.Println("b is :", b,b & 0x02)
}
