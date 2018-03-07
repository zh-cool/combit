package main

import (
	"fmt"
	_ "reflect"
	"go-unitcoin/libraries/common"
	"go-unitcoin/libraries/crypto/sha3"
	"bytes"
	"encoding/binary"
)

func main() {
	var mTree [1 << 2]common.Hash
	b := []byte{0}

	hw := sha3.NewKeccak256()
	hw.Write(b)
	hw.Sum(mTree[2][:0])
	hw.Sum(mTree[3][:0])
	hw.Reset()

	hw.Write(mTree[2][:])
	hw.Write(mTree[3][:])
	hw.Sum(mTree[1][:0])
	fmt.Println(mTree)

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, mTree)
	//fmt.Println(buf.Bytes())
	if err != nil {
		panic(err)
	}

	var gTree [1 << 2]common.Hash
	binary.Read(buf, binary.BigEndian, &gTree)
	fmt.Println(gTree)

	var hash []byte
	hash = []byte("string")
	fmt.Println(hash)
}
