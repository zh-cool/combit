package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)


func main() {


	db, _ := ethdb.NewLDBDatabase("trie", 0, 0)
	defer db.Close()


	tree, err := trie.New(common.Hash{}, trie.NewDatabase(db))
	if err!=nil {
		fmt.Println(err)
	}
	tree.Update()

	result := [32]byte{}
	data := []byte{}

	hw := sha3.NewKeccak256()
	hw.Write(data[:])
	hw.Sum(result[:0])
	fmt.Printf("%x\n", result)
	fmt.Printf("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
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
