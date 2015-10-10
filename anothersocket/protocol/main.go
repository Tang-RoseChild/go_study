package protocol

import (
	"encoding/binary"
	// "bufio"
	"bytes"
	"fmt"
)

const (
	HEAD1       = 011
	HEAD2       = 022 // normal msg
	HEART       = 033 // heart beat
	LSIZE       = 4   // msg length size like int32 is 4 byte
	TSIZE       = 4   // msg topic length size,like LSIZE
	HEADER_SIZE = 6
)

func Pack(msg []byte) []byte {
	// get msg leng
	lMsg := len(msg)
	return append(append([]byte{HEAD1, HEAD2}, Itob(lMsg)...), msg...)
}

// func Unpack : return the lefer data
func Unpack(packer []byte, msgReceiver chan []byte) []byte {
	var i int
	lMsg := len(packer)
	// fmt.Printf("%o %o before loop %v", packer[0], packer[1], []byte{HEAD1, HEAD2})
	for i = 0; i < lMsg; i++ {
		if i+HEADER_SIZE > lMsg {
			// fmt.Println("break ")
			break
		}
		// head wrong
		if packer[i] != HEAD1 || packer[i+1] != HEAD2 {
			// fmt.Printf(" %d break in header %o %o : %o %o\n", i, packer[i], packer[i+1], HEAD1, HEAD2)
			continue
		}

		// msg content not enough
		cSize := Btoi(packer[i+2 : i+2+LSIZE])
		// fmt.Println("cSize ", cSize, " ", packer[i+HEADER_SIZE:i+HEADER_SIZE+cSize])
		if (lMsg - i) < (HEADER_SIZE + cSize) {
			// fmt.Println("break after cSize")
			break

		}

		// send msg to receiver
		msgReceiver <- packer[i+HEADER_SIZE : i+HEADER_SIZE+cSize]
		// fmt.Println("send msg to receiver : ", packer[i+HEADER_SIZE:i+HEADER_SIZE+cSize])

	}

	// packer is just ok
	if i == lMsg {
		return make([]byte, 0)
	}
	// fmt.Println("i : ", i, string(packer[i:]))
	// more data left
	return packer[i:]
}

func Itob(i int) []byte {
	i32 := int32(i)
	bf := new(bytes.Buffer)

	if err := binary.Write(bf, binary.BigEndian, i32); err != nil {
		fmt.Println("error when write in Itob : ", err.Error())
		return make([]byte, 0)
	}
	return bf.Bytes()
}
func Btoi(b []byte) int {
	var i32 int32
	bf := bytes.NewBuffer(b)
	if err := binary.Read(bf, binary.BigEndian, &i32); err != nil {
		fmt.Println("error when write in Itob : ", err.Error())
		return 0
	}
	return int(i32)
}
