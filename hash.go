package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

func main() {

	result := [32]byte{}
	data := [32]byte{}

	var i uint64

	for i = 0; i < 1<<2; i++ {
		hw := sha3.NewKeccak256()
		hw.Write(data[:])
		hw.Sum(result[:0])
	}

	hw := sha3.NewKeccak256()
	hw.Write(data[:1])
	hw.Write(data[1:2])
	hw.Sum(result[:0])
	fmt.Printf("%x\n", result)

	hw = sha3.NewKeccak256()
	hw.Write(data[:2])
	hw.Sum(result[:0])
	fmt.Printf("%x\n", result)

	hw.Sum(result[:0])
	fmt.Printf("%x\n", result)

	hw.Sum(result[:0])
	fmt.Printf("%x\n", result)

}

func Hash() {
	var data [1 << 23]common.Hash

	i := len(data)/2 - 1

	for ; i >= 1; i-- {
		left := i << 1
		right := i<<1 + 1
		hw := sha3.NewKeccak256()
		hw.Write(data[left][:])
		hw.Write(data[right][:])
		hw.Sum(data[i][:0])
	}
	fmt.Printf("%x\n", data[1])
}

func sum(b []byte) []byte {
	return append(b, 'a')
}
