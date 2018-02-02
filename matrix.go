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
		idx := i << 1
		fmt.Printf("%08b :%3d %08b :%3d\n", data[idx], data[idx], data[idx+1], data[idx+1])
	}
}

func (bits Bigbit) display(index int) {
	bprint(bits.bit[0][index*32 : index*32+32])
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
	bprint(m.matric[:])
	fmt.Printf("Matric flag:%d\n", m.flag)
}

func CreateMatric(data []byte) []Matric {

	CMatric := len(data) >> 5

	matric := make([]Matric, CMatric)

	for i := 0; i < CMatric; i++ {
		for j := 0; j < 32; j++ {
			matric[i].matric[j] |= data[j+i<<5]
			matric[i].flag |= uint8(data[j+i<<5])
			matric[i].sum += uint(data[j+i<<5])
		}
		if matric[i].flag > 0 {
			matric[i].flag = 1
		}
	}

	return matric
}

func Statistics(m []Matric) int {
	count := 0
	for _, v := range m {
		count += int(v.flag)
	}
	return count
}

func prepareData(matric []Matric) []byte {
	size := len(matric) >> 3
	data := make([]byte, size)
	for i, v := range matric {
		if v.flag > 0 {
			data[i/8] |= 1 << (7 - uint(i%8))
			v.flag = 1
		}
	}
	return data
}

func main() {
	bit := NewBitbit(1, 1<<32)
	fmt.Println("Org data")
	bit.display(0)

	matric_l1 := CreateMatric(bit.bit[0])
	fmt.Printf("Matric L1 len:%d contain data:%d %v\n", len(matric_l1), Statistics(matric_l1), float64(Statistics(matric_l1))/float64(len(matric_l1)))
	matric_l1[0].display()

	data := prepareData(matric_l1)
	fmt.Println()
	bprint(data[0:32])

	fmt.Println("---------next level-------")
	matric_l2 := CreateMatric(data)
	data = prepareData(matric_l2)
	fmt.Printf("Matric L2 len:%d contain data:%d %v\n", len(matric_l2), Statistics(matric_l2), float64(Statistics(matric_l2))/float64(len(matric_l2)))
	matric_l2[0].display()
	fmt.Println()
	bprint(data[0:32])

	fmt.Println("---------next level-------")
	matric_l3 := CreateMatric(data)
	data = prepareData(matric_l3)
	fmt.Printf("Matric L3 len:%d contain data:%d %v\n", len(matric_l3), Statistics(matric_l3), float64(Statistics(matric_l3))/float64(len(matric_l3)))
	matric_l3[0].display()
	fmt.Println()
	bprint(data[0:32])

	fmt.Println("---------next level-------")
	matric_l4 := CreateMatric(data)
	fmt.Println("Matric L4 len:", len(matric_l4))
	matric_l4[0].display()
}
