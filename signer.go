package main

import (
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/common"
	 "crypto/ecdsa"
	"fmt"
	"crypto/elliptic"
	"crypto/rand"
)

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func main() {
	data := []byte("Austin")

	prv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%x %x %x\n", prv, prv.X.Bytes(), prv.Y.Bytes())

	h := rlpHash(data)
	r, s, err := ecdsa.Sign(rand.Reader, prv, h[:])
	fmt.Printf("%x %d %x %d\n", r.Bytes(), len(r.Bytes()), s.Bytes(), len(s.Bytes()))

	ok := ecdsa.Verify(&prv.PublicKey, h[:], r, s)
	fmt.Println(ok)
}
