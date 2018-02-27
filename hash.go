package main

import (
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"fmt"
)

func main() {
	result := [32]byte{}
	data := []byte("01234567890123456789012345678901")

	var i uint64

	hw := sha3.NewKeccak256()
	for i=0; i<1<<20; i++ {
		hw.Write(data)
		hw.Sum(result[:0])
		data = result[:]
		hw.Reset()
		//fmt.Printf("%X\n", result)
	}

	//b := []byte("1234")
	fmt.Printf("%x\n", result)

}

func sum(b []byte) []byte {
	return append(b, 'a')
}
