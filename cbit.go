package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

const BIT_LENGTH = 1 << 32

type Cbit struct {
	b     uint8
	count int
}

type Bigbit struct {
	bit   [][]byte
	cmbit []Cbit
}

func NewBitbit(num int, length int64) *Bigbit {

	var bit Bigbit
	buf := new(bytes.Buffer)
	for i := 0; i < num; i++ {
		src := rand.NewSource(time.Now().Unix() + int64(i))
		rd := rand.New(src)
		data := rd.Perm(int(length >> 5))

		for _, v := range data {
			binary.Write(buf, binary.BigEndian, int32(v)|int32(v<<16))
		}

		s := buf.Bytes()
		bit.bit = append(bit.bit, s)
	}
	fmt.Println(bit.bit)

	return &bit
}

func (bit *Bigbit) Count() {
	var cb Cbit
	var TBIT uint8
	var i uint
	TBIT = 128
	for _, v := range bit.bit {
		cb = Cbit{128, -1}
		for _, bt := range v {
			for i = 0; i < 8; i++ {
				bi := uint8(bt) & (TBIT >> i)
				bi = bi >> (7 - i)

				if bi == cb.b {
					cb.count++
				} else {
					bit.cmbit = append(bit.cmbit, cb)
					cb.b = bi
					cb.count = 1
				}
			}
		}
	}
}

func (bit *Bigbit) Print() {
	fmt.Println("Org data")
	for _, v := range bit.bit {
		fmt.Printf("[")
		for _, b := range v {
			fmt.Printf(" %08b", b)
		}
		fmt.Printf("]\n")
	}
	fmt.Println(bit.cmbit)
}

func main() {
	bit := NewBitbit(1, 1<<12)
	bit.Count()
	bit.Print()
}
