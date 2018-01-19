package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const BIT_LENGTH = 1 << 32

type Cbit struct {
	b     uint8
	count int
}

type Bigbit struct {
	bit       [][]byte
	cmbit     []Cbit
	hexdata   []byte
	hexstring string
}

func NewBitbit(num int, length int64) *Bigbit {
	var bit Bigbit

	//buf := new(bytes.Buffer)
	for i := 0; i < num; i++ {
		bit.bit = append(bit.bit, make([]byte, length>>3))
		src := rand.NewSource(time.Now().Unix() + int64(i))
		rd := rand.New(src)
		/*
			data := rd.Perm(int(length >> 5))
			for _, v := range data {
				binary.Write(buf, binary.BigEndian, int32(v)|int32(v<<16))
			}

			s := buf.Bytes()
			bit.bit = append(bit.bit, s)
		*/
		for j := 0; j < int(length>>8); j++ {
			pos := rd.Int63n(length)
			bytepos := pos >> 8
			bitpos := pos % 8
			bit.bit[i][bytepos] |= (1 << uint(bitpos))
		}
	}

	return &bit
}

func (bit *Bigbit) Count() {
	var cb Cbit
	var TBIT uint8
	var i uint
	TBIT = 128
	var shex string
	var hexstring []string

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
		bit.cmbit = append(bit.cmbit, cb)
	}

	for _, vv := range bit.cmbit[1:] {
		hsdata := []string{}
		if vv.count > 1 {
			//shex = fmt.Sprintf("%b%x", uint8(vv.b), uint32(vv.count))
			shex = strconv.FormatInt(int64(vv.b), 16)
			hsdata = append(hsdata, shex)
			shex = strconv.FormatUint(uint64(vv.count), 16)
			hsdata = append(hsdata, shex)
		} else {
			//shex = fmt.Sprintf("%b", vv.b)
			shex = strconv.FormatInt(int64(vv.b), 16)
			hsdata = append(hsdata, shex)
		}
		/*
			hbyte := []byte(strings.Join(hsdata, ""))
			for i, v := range hbyte[1:] {
				if v == 0 {
					hbyte[i] = 'g'
				} else if v == 1 {
					hbyte[i] = 'h'
				}
			}
		*/
		hexstring = append(hexstring, strings.Join(hsdata, ""))
	}
	bit.hexstring = strings.Join(hexstring, "")
	bit.hexdata = []byte(bit.hexstring)
}

func (bit *Bigbit) Statistics() {
	var bitsize int64
	var bytesize int64
	var zosize [2]int64

	fmt.Println("Org data statistics:")
	for _, v := range bit.bit {
		bytesize += int64(len(v))
	}
	bitsize = bytesize << 3

	for _, v := range bit.cmbit[1:] {
		zosize[v.b] += int64(v.count)
	}

	fmt.Printf("size:=(%dbits %dbyte %dKB, %dMB) 1:%d 0:%d\n",
		bitsize, bytesize, bytesize>>10, bytesize>>20,
		zosize[1], zosize[0])

	fmt.Println("Compress data statistics:")
	bytesize = int64(len(bit.hexdata))
	bitsize = bytesize << 3
	fmt.Printf("size:=(%dbits %dbyte %dKB, %dMB)\n",
		bitsize, bytesize, bytesize>>10, bytesize>>20,
	)

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
	fmt.Println("Compress data")
	fmt.Println(bit.hexstring)
}

func main() {
	bit := NewBitbit(1, 1<<12)
	bit.Count()
	bit.Print()
	bit.Statistics()
}
