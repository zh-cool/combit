package main

import (
	"go-unitcoin/libraries/rlp"
	"bytes"
	"fmt"
	_ "github.com/ethereum/go-ethereum/trie"
)

type txdata  struct{
	A	uint
	B	uint
	Name	string
}

type TestRlp struct{
	Data txdata

	Value int32
	S     []uint
}

var tr TestRlp = TestRlp{
	Data:	txdata{
		A:	3,
		B:  4,
		Name: "Austin",
	},

	Value:	5,
	S:	[]uint{ 2, 3, 4, 55},
}

func main() {
	dat := []*TestRlp{&tr, &tr, &tr}

	//trie.New()

	buf := new(bytes.Buffer)

	err := rlp.Encode(buf, dat);
	fmt.Printf("%v %x\n", err, buf.Bytes())

	buf.Reset()
	err = rlp.Encode(buf, &dat);
	fmt.Printf("%v %x\n", err, buf.Bytes())

	var dtr []*TestRlp
	err = rlp.Decode(buf, &dtr)
	fmt.Printf("%v %v\n", err, dtr)
	for _, v := range dtr {
		fmt.Printf("%v %v\n", err, v)
	}

	fmt.Println("****************")
	buf.Reset()
	val := []string{"Hello", "Workd", "Austin"}
	rlp.Encode(buf, val)
	fmt.Printf("%x\n", buf.Bytes())

	var dval []string
	rlp.Decode(buf, &dval)
	fmt.Println(dval)


}
