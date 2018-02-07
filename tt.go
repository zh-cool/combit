package main

import (
	//	"flag"
	"fmt"
	//	"log"
	"math/rand"
	//	"os"
	//"runtime/pprof"
	"strconv"
	//	"strings"
	"bytes"
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
			bytepos := pos >> 3
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
	//var hexstring []string

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

	ozstring := []string{"0", "1"}
	var buffer bytes.Buffer
	for _, vv := range bit.cmbit[1:] {
		hsdata := []string{}
		if vv.count > 1 {
			hsdata = append(hsdata, ozstring[vv.b])
			shex = strconv.FormatUint(uint64(vv.count), 16)
			hsdata = append(hsdata, shex)
			bthex := []byte(hsdata[1])
			for i, v := range bthex {
				if v == '0' {
					bthex[i] = 'g'
				} else if v == '1' {
					bthex[i] = 'h'
				}
			}
			hsdata[1] = string(bthex)
			//hexstring = append(hexstring, strings.Join(hsdata, ""))
			buffer.WriteString(hsdata[0])
			buffer.WriteString(hsdata[1])
		} else {
			hsdata = append(hsdata, ozstring[vv.b])
			//hexstring = append(hexstring, strings.Join(hsdata, ""))
			buffer.WriteString(hsdata[0])
		}

	}
	//bit.hexstring = strings.Join(hexstring, "")
	bit.hexstring = buffer.String()
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
