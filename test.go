package main

import (
	"fmt"
	"math/big"
)
func sub(in []byte) []byte {
	a := new(big.Int).SetUint64(1)
	return append(in, a.Bytes()...)
}
func main() {
	b := [32]byte{}
	a := new(big.Int).SetUint64(1)

	//append(b[:0], a.Bytes()...)

	fmt.Println(append(b[:0], a.Bytes()...))
}
