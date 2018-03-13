package main

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
	"fmt"
	_"github.com/ethereum/go-ethereum/crypto/sha3"
	_"go-unitcoin/libraries/rlp"
	_"bytes"
	"bytes"
	"go-unitcoin/libraries/rlp"
	"go-unitcoin/libraries/crypto/sha3"
)

func newEmpty() *trie.Trie {
	diskdb, _ := ethdb.NewMemDatabase()
	trie, _ := trie.New(common.Hash{}, trie.NewDatabase(diskdb))
	return trie
}

func TestInsert() {
	trie := newEmpty()

	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")

	exp := common.HexToHash("8aad789dff2f538bca5d8ea56e8abe10f4c7ba3a5dea95fea4cd6e7c3a1168d3")
	root := trie.Hash()
	if root != exp {
		fmt.Printf("exp %x got %x", exp, root)
	}

	trie = newEmpty()
	updateString(trie, "A", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	exp = common.HexToHash("d23786fb4a010da3ce639d66d5e904a11dbc02746d1ce25029e53290cabf28ab")
	root, err := trie.Commit(nil)
	if err != nil {
		fmt.Printf("commit error: %v", err)
	}
	if root != exp {
		fmt.Printf("exp %x got %x\n", exp, root)
	}
}

type NOD struct {
	Key	[]byte
	Value []byte
}

var nod NOD = NOD{
	Key: []byte{32, 65},
	Value: []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
}

func main() {
	TestInsert()
	trie := newEmpty()
	updateString(trie, "A", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	hash := trie.Hash();
	fmt.Printf("%x\n", hash.Bytes())

	buf := new(bytes.Buffer)
	rlp.Encode(buf, nod)
	fmt.Println(buf.Bytes())

	hw := sha3.NewKeccak256()
	hw.Write(buf.Bytes())
	hw.Sum(hash[:0])
	fmt.Printf("%x\n", hash.Bytes())
}

func updateString(trie *trie.Trie, k, v string) {
	trie.Update([]byte(k), []byte(v))
}