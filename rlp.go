package main

import (
	"go-unitcoin/libraries/rlp"
	"bytes"
	"fmt"
	"io"
	"math/big"
)
type MyCoolType struct {
	Name []byte
	a, b uint
}

// EncodeRLP writes x as RLP list [a, b] that omits the Name field.
func (x *MyCoolType) EncodeRLP(w io.Writer) (err error) {
	// Note: the receiver can be a nil pointer. This allows you to
	// control the encoding of nil, but it also means that you have to
	// check for a nil receiver.
	if x == nil {
		err = rlp.Encode(w, []uint{0, 0})
	} else {
		err = rlp.Encode(w, []uint{x.a, x.b})
	}
	return err
}

func ExampleEncoder() {
	var t *MyCoolType // t is nil pointer to MyCoolType
	cbytes, _ := rlp.EncodeToBytes(t)
	fmt.Printf("%v → %X\n", t, cbytes)

	t = &MyCoolType{Name: []byte("foobar"), a: 5, b: 6}
	cbytes, _ = rlp.EncodeToBytes(t)
	fmt.Printf("%v → %X\n", t, cbytes)
}


type structWithTail struct {
	A, B uint
	C    []uint `rlp:"tail"`
}

func ExampleDecode_structTagTail() {
	// In this example, the "tail" struct tag is used to decode lists of
	// differing length into a struct.
	var val structWithTail

	err := rlp.Decode(bytes.NewReader([]byte{0xC4, 0x01, 0x02, 0x03, 0x04}), &val)
	fmt.Printf("with 4 elements: err=%v val=%v\n", err, val)

	err = rlp.Decode(bytes.NewReader([]byte{0xC6, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}), &val)
	fmt.Printf("with 6 elements: err=%v val=%v\n", err, val)

	// Note that at least two list elements must be present to
	// fill fields A and B:
	err = rlp.Decode(bytes.NewReader([]byte{0xC1, 0x01}), &val)
	fmt.Printf("with 1 element: err=%q\n", err)

	// Output:
	// with 4 elements: err=<nil> val={1 2 [3 4]}
	// with 6 elements: err=<nil> val={1 2 [3 4 5 6]}
	// with 1 element: err="rlp: too few elements for rlp.structWithTail"
}

type TestRlpStruct struct {
	A      uint
	B      string
	C      []byte
	BigInt *big.Int
}

//rlp用法
func TestRlp() {
	//1.将一个整数数组序列化
	arrdata, err := rlp.EncodeToBytes([]uint{32, 28})
	fmt.Printf("unuse err:%v\n", err)
	//fmt.Sprintf("data=%s,err=%v", hex.EncodeToString(arrdata), err)
	//2.将数组反序列化
	var intarray []uint
	err = rlp.DecodeBytes(arrdata, &intarray)
	//intarray 应为{32,28}
	fmt.Printf("intarray=%v\n", intarray)

	//3.将一个布尔变量序列化到一个writer中
	writer := new(bytes.Buffer)
	err = rlp.Encode(writer, true)
	//fmt.Sprintf("data=%s,err=%v",hex.EncodeToString(writer.Bytes()),err)
	//4.将一个布尔变量反序列化
	var b bool
	err = rlp.DecodeBytes(writer.Bytes(), &b)
	//b:true
	fmt.Printf("b=%v\n", b)

	//5.将任意一个struct序列化
	//将一个struct序列化到reader中
	_, r, err := rlp.EncodeToReader(TestRlpStruct{3, "44", []byte{0x12, 0x32}, big.NewInt(32)})
	data := [128]byte{}
	r.Read(data[:])
	fmt.Println(data)
	/*
	var teststruct TestRlpStruct
	err = rlp.Decode(r, &teststruct)
	//{A:0x3, B:"44", C:[]uint8{0x12, 0x32}, BigInt:32}
	fmt.Printf("teststruct=%#v\n", teststruct)
	*/
}

type testRlp struct {
	A	uint
	B	uint
	C	string
}

func main() {
	ts := testRlp {1, 2, "Austin"}

	b := new(bytes.Buffer)
	rlp.Encode(b, &ts)

	var tt testRlp
	rlp.Decode(b, &tt)
	fmt.Println(tt)
}
