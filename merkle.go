package main

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/common"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
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

	tried.Update([]byte("1234"), []byte("Austin"))
	tried.Update([]byte("12345"), []byte("Austin key"))
	root, _ := tried.Commit(nil)

	tried, _ = trie.New(root,triedb)
	val, err := tried.TryGet([]byte("1234"));
	if err!=nil {
		fmt.Println(err)
	}

	fmt.Printf("%s\n", val)

	val = tried.Get([]byte("12345"))
	fmt.Printf("%s\n", val)

	b := crypto.Keccak256([]byte("Austin"))
	fmt.Printf("%x\n", b)
}