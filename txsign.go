package main

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
	"fmt"
)

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func main() {
	data := []byte("Austin")

	h := rlpHash(data)
	prv, err := crypto.GenerateKey()
	fmt.Printf("prv: %x %v\n", *prv, err)


	sig, err := crypto.Sign(h[:], prv)
	fmt.Printf("sig:%x %v\n", sig, err)
	//prv.

	pub, err := crypto.Ecrecover(h[:], sig)
	fmt.Printf("pub:%x %d %v\n", pub, len(pub), err)

	ok := crypto.VerifySignature(pub, h[:], sig[:64])
	fmt.Println(ok)
}
