package main

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/common"
	"fmt"
	"go-unitcoin/libraries/crypto/sha3"
)

func newEmpty() *trie.Trie {
	diskdb, _ := ethdb.NewMemDatabase()
	trie, _ := trie.New(common.Hash{}, trie.NewDatabase(diskdb))
	return trie
}

func main() {
	diskdb, _ := ethdb.NewMemDatabase()
	triedb := trie.NewDatabase(diskdb)

	tried, _ := trie.New(common.Hash{}, triedb)

	hs := [2][32]byte{}
	hw := sha3.NewKeccak256()
	hw.Write([]byte("Austin"))
	hw.Sum(hs[0][:0])

	hw = sha3.NewKeccak256()
	hw.Write([]byte("Austin key"))
	hw.Sum(hs[1][:0])

	tried.Update(hs[0][:], []byte("Austin"))
	tried.Update(hs[1][:], []byte("Austin key"))
	h := tried.Hash()
	fmt.Printf("%x\n", h.Bytes())

	result := [32]byte{}
	hw = sha3.NewKeccak256()
	hw.Write(hs[0][:])
	hw.Write(hs[1][:])
	hw.Sum(result[:0])
	fmt.Printf("%x\n", result)
}