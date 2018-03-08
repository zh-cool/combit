package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"time"
)

func main() {

	fmt.Println(time.Now().Unix())
	Hash()
	fmt.Println(time.Now().Unix())

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

