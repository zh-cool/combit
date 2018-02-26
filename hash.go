package main

import (
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"fmt"
)

func main() {
	result := [32]byte{}
	data := []byte("0")

	var i uint64

	hw := sha3.NewKeccak256()
	for i=0; i<1<<2; i++ {
		hw.Write(data)
		hw.Sum(result[:0])
		hw.Reset()
		fmt.Printf("%X\n", result)
	}

	b := []byte("1234")
	fmt.Printf("%x\n", b)

}

func sum(b []byte) []byte {
	return append(b, 'a')
}
