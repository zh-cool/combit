package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	BIT_SIZE    = (1 << 32)
	MATRIX_SIZE = 16
)

type Bigbit struct {
	bit [][]byte
}

func bprint(data []byte) {
	for i := 0; i < 16; i++ {
		fmt.Printf("%08b :%d %08b :%d\n", data[i], data[i], data[i+1], data[i+1])
	}
}

func (bits Bigbit) display(index int) {
	for i := 0; i < 16; i++ {
		fmt.Printf("%08b :%d %08b :%d\n", bits.bit[0][index*32+i], bits.bit[0][index*32+i], bits.bit[0][index*32+i+1], bits.bit[0][index*32+i+1])
	}
	fmt.Println()
}

func NewBitbit(num int, length int64) *Bigbit {
	var bit Bigbit

	for i := 0; i < num; i++ {
		bit.bit = append(bit.bit, make([]byte, length>>3))
		src := rand.NewSource(time.Now().Unix() + int64(i))
		rd := rand.New(src)

		for j := 0; j < int(length>>8); j++ {
			pos := rd.Int63n(length)
			bytepos := pos >> 3
			bitpos := pos % 8
			bit.bit[i][bytepos] |= (1 << uint(bitpos))
		}
	}

	return &bit
}

type Matric struct {
	matric [MATRIX_SIZE * MATRIX_SIZE >> 3]byte
	flag   uint8
	sum    uint
}

func (m Matric) display() {
	for i := 0; i < 16; i++ {
		fmt.Printf("%08b :%d %08b %d\n", m.matric[i], m.matric[i], m.matric[i+1], m.matric[i+1])
	}
	fmt.Println("Matric flag:", m.flag)
}

func CreateMatric(data []byte) []Matric {
	//bits := len(data) << 3
	bprint(data[0:32])
	CMatric := len(data) >> 5

	matric := make([]Matric, CMatric)

	for i := 0; i < CMatric; i++ {
		for j := 0; j < 32; j++ {
			matric[i].matric[j] |= data[j+i<<5]
			matric[i].flag |= uint8(data[j+i<<5])
			matric[i].sum += uint(data[j+i<<5])
			if i == 0 {
				fmt.Printf("T%d %d %d ", j, matric[i].sum, uint(data[j+i<<5]))
			}
		}
		/*
			if matric[i].flag > 0 {
				matric[i].flag = 1
			}
		*/
	}
	fmt.Println("CreateMatric", CMatric, matric[0].sum, matric[0].flag)
	bprint(data[0:32])
	fmt.Println()
	matric[0].display()
	fmt.Println("CreateMatric", CMatric)
	return matric
}

func prepareData(matric []Matric) []byte {
	size := len(matric) >> 3
	data := make([]byte, size)
	for i, v := range matric {
		if v.flag > 0 {
			data[i/8] |= 1 << uint(i%8)
			v.flag = 1
		}
	}
	return data
}

func main() {
	fmt.Println("Austin Test")
	bit := NewBitbit(1, 1<<16)
	bit.display(0)

	matric := CreateMatric(bit.bit[0])
	fmt.Println(len(matric))
	matric[0].display()

	data := prepareData(matric)
	for i := 0; i < 256; i++ {
		fmt.Printf("%d ", matric[i].flag)
	}
	fmt.Println()
	bprint(data[0:32])

}
