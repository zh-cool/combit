package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

const BIT_LENGTH = 1 << 32

type Bigbit struct {
	bit   [][]byte
	cmbit [][]byte
}

func NewBitbit(num int, length int64) *Bigbit {

	var bit Bigbit
	buf := new(bytes.Buffer)
	for i := 0; i < num; i++ {
		src := rand.NewSource(time.Now().Unix() + int64(i))
		rd := rand.New(src)
		data := rd.Perm(int(length >> 5))

		for _, v := range data {
			binary.Write(buf, binary.BigEndian, int32(v))
		}

		s := buf.Bytes()
		bit.bit = append(bit.bit, s)
		bit.bit = append(bit.bit, s)

		fmt.Println(len(s), s)
		fmt.Println(len(data), data)
		fmt.Println(bit)

	}

	return &bit
}

func main() {
	NewBitbit(1, 1<<10)
}
