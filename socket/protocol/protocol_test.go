package protocol

import (
	// "fmt"
	"reflect"
	"testing"
)

func TestItob(t *testing.T) {
	testData := []struct {
		i int
		b []byte
	}{
		{10, []byte{0, 0, 0, 10}},
		{256, []byte{0, 0, 1, 0}},
		{1025, []byte{0, 0, 4, 1}},
	}
	for _, data := range testData {
		b := Itob(data.i)
		// if b != data.b {
		// 	t.Errorf(" want %s but %s \n", data.b, b)
		// }
		// fmt.Println(b, " result : ", reflect.DeepEqual(b, data.b))
		if !reflect.DeepEqual(b, data.b) {
			t.Errorf(" want %s but %s \n", data.b, b)
		}
	}

}

func TestBtoi(t *testing.T) {

	testData := []struct {
		i int
		b []byte
	}{
		{10, []byte{0, 0, 0, 10}},
		{256, []byte{0, 0, 1, 0}},
		{1025, []byte{0, 0, 4, 1}},
	}

	for _, data := range testData {
		i := Btoi(data.b)
		if i != data.i {
			t.Errorf("Btoi want %d, but get %d \n", data.i, i)
		}
	}
}

func TestPack(t *testing.T) {
	testData := []struct {
		msg    []byte
		packer []byte
	}{
		{[]byte("hello"), append(append([]byte(HEAD), Itob(len("hello"))...), []byte("hello")...)},
		{[]byte(""), append(append([]byte(HEAD), Itob(len(""))...), []byte("")...)},
		{[]byte("hellofdas fljdsa ;fdlsajfldjsafj9oiepwfjoepowqjfepfoewjqfldsja"), append(append([]byte(HEAD), Itob(len("hellofdas fljdsa ;fdlsajfldjsafj9oiepwfjoepowqjfepfoewjqfldsja"))...), []byte("hellofdas fljdsa ;fdlsajfldjsafj9oiepwfjoepowqjfepfoewjqfldsja")...)},
		{[]byte("fwa"), append(append([]byte(HEAD), Itob(len("fwa"))...), []byte("ziw")...)},
	}

	for i, data := range testData {
		p := Pack(data.msg)

		if i != len(testData)-1 && !reflect.DeepEqual(p, data.packer) {
			t.Errorf("Pack want %v, but get %v \n", data.packer, p)
		}
		if i == len(testData)-1 && reflect.DeepEqual(p, data.packer) {
			t.Errorf("Pack want %v, but get %v \n", data.packer, p)
		}
	}
}

// func TestUnpack(t *testing.T) {
// 	testData := []struct {
// 		msg    []byte
// 		packer []byte
// 	}{
// 		{[]byte("hello"), append([]byte("what the hell "), append(append([]byte(HEAD), Itob(len("hello"))...), []byte("hello")...)...)},
// 		{[]byte(""), append(append([]byte(HEAD), Itob(len(""))...), []byte("")...)},
// 		{[]byte("ab"), append(append(append([]byte(HEAD), Itob(len("ab"))...), []byte("ab")...), []byte("more data in the end")...)},
// 		{nil, []byte(HEAD)},
// 		{[]byte("hellofdas fljdsa ;fdlsajfldjsafj9oiepwfjoepowqjfepfoewjqfldsja"), append(append([]byte(HEAD), Itob(len("hellofdas fljdsa ;fdlsajfldjsafj9oiepwfjoepowqjfepfoewjqfldsja"))...), []byte("hellofdas fljdsa ;fdlsajfldjsafj9oiepwfjoepowqjfepfoewjqfldsja")...)},
// 		{[]byte("fwa"), append(append([]byte(HEAD), Itob(len("fwa"))...), []byte("ziw")...)},
// 	}

// 	for i, data := range testData {
// 		unp := Unpack(data.packer)

// 		if i != len(testData)-1 && !reflect.DeepEqual(data.msg, unp) {
// 			t.Errorf("Pack want %v, but get %v \n", data.msg, unp)
// 		}
// 		if i == len(testData)-1 && reflect.DeepEqual(data.msg, unp) {
// 			t.Errorf("Pack want %v, but get %v \n", data.msg, unp)
// 		}
// 	}
// }
