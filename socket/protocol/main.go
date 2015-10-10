package protocol

import (
	"bytes"
	"encoding/binary"
	// "fmt"
)

const (
	HEAD      = "HEAD"
	HEAD_SIZE = len(HEAD)
	MSG_SIZE  = 4 // use int32 to store msg content size
)

// Pack : pack msg and return result
func Pack(msg []byte) []byte {
	msgContentSize := len(msg)

	// add head and MSG_SIZE
	msgCSB := Itob(msgContentSize)
	return append(append([]byte(HEAD), msgCSB...), msg...)

}

// Unpack: from packer to get neccesiry msg content
func Unpack(packer []byte, msgs chan []byte) []byte {
	l := len(packer)
	// fmt.Println(l, " :", "packer : ", string(packer))
	var i = 0
	// because packer may get more data than msg
	for i = 0; i < l; i++ {
		headerSize := i + HEAD_SIZE + MSG_SIZE
		if l < headerSize {
			break
		}

		// whether head is right
		if string(packer[i:i+HEAD_SIZE]) != HEAD {
			continue
		}

		// get size
		msgSize := Btoi(packer[i+HEAD_SIZE : headerSize])
		if l < headerSize+msgSize {
			break
		}

		msgs <- packer[headerSize : headerSize+msgSize]
		// fmt.Println(i, " : upack msg in for loop: ", string(packer[headerSize:headerSize+msgSize]))
		// return packer[headerSize+msgSize:]
	}

	// if msg is absolute the same,return an empty slice
	if i == l {
		// fmt.Println("equal *************")
		return make([]byte, 0)
	}

	// if packer get more msg: return the left data(may be part of the msg)
	return packer[i:]
}

// Itob : convert int to int32(if int is by default 64, this convert is kind of danger),and return bigendian bytes
func Itob(l int) []byte {
	i32 := int32(l)

	bf := new(bytes.Buffer)

	err := binary.Write(bf, binary.BigEndian, i32)
	if err != nil {
		return nil
	}

	return bf.Bytes()
}

// Btoi : well opposite of Itob,change bigendian bytes to int
func Btoi(b []byte) int {

	bf := bytes.NewBuffer(b)

	var size int32

	err := binary.Read(bf, binary.BigEndian, &size)
	if err != nil {
		return 0
	}

	return int(size)

}
