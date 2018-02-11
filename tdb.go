package main

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
	"fmt"
	"go-unitcoin/libraries/common"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	CONBIN_BIT_SIZE  = (1 << 32)
	CONBIN_BYTE_SIZE = (CONBIN_BIT_SIZE >> 3)

	FG_LENGTH    = (1 << 10)
	FG_BIT_SIZE  = (CONBIN_BIT_SIZE / FG_LENGTH)
	FG_BYTE_SIZE = (FG_BIT_SIZE >> 3)
)

type FG_data struct {
	//	Len int
	Data []byte
}

type Wallet struct {
	ID   common.Address
	Data [FG_LENGTH]FG_data
}

func (w *Wallet) Bytes() []byte {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(w)
	if err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}

func (w *Wallet) toWallet(data []byte) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(w)
	if err != nil {
		fmt.Println(err)
	}
}

func (w *Wallet) Set(bitpos int64) {
	conv := map[byte]int64{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'a': 10,
		'b': 11,
		'c': 13,
		'd': 14,
		'e': 15,
		'f': 16,
		'g': 0,
		'h': 1,
	}

	var i, k int64
	fgpos := bitpos / FG_BIT_SIZE
	bitpos = bitpos % FG_BIT_SIZE

	data := w.Data[fgpos].Data
	length := len(data)

	big := &big.Int{}
	big.SetBytes(make([]byte, FG_BYTE_SIZE))

	bigpos := 0
	for i = 0; i < int64(length); {
		ch := data[i]
		bit := conv[ch]
		if ch == '1' || ch == '0' {
			s, l := func(data []byte) (sum int64, length int64) {
				sum = 0
				length = 1
				for i := 0; i < len(data); i++ {

					ch := data[i]
					if ch == '1' || ch == '0' {
						return sum+1, length
					}
					sum = sum*16 + conv[ch]
					length++
				}
				return sum+1, length
			}(data[i+1:])

			i += l
			for k=0; k<s; k++ {
				big.SetBit(big, bigpos, uint(bit))
				bigpos++
			}
		}
	}
	big.SetBit(big, int(bitpos), 1)
	w.Data[fgpos]=CreateOneFG(big.Bytes())
}

func (w *Wallet) FGCompress() {
	for _, v := range w.Data {
		b := bytes.NewBuffer([]byte{})
		wr := zlib.NewWriter(b)
		wr.Write(v.Data)
		wr.Close()
		v.Data = b.Bytes()
	}
}

func (w *Wallet) Compress() []byte {
	b := bytes.NewBuffer([]byte{})
	W := zlib.NewWriter(b)

	for _, v := range w.Data {
		W.Write(v.Data)
	}
	W.Close()
	return b.Bytes()
}

func (w *Wallet) Statistics() {

	var sum int64
	b := bytes.NewBuffer([]byte{})
	W := zlib.NewWriter(b)

	for _, v := range w.Data {
		sum += int64(len(v.Data))
		W.Write(v.Data)
	}
	W.Close()

	l := len(b.Bytes())
	fmt.Printf("Org Data:%dByte %dKB\n", sum, sum/1024)
	fmt.Printf("Compress %dByte %dKB\n", l, l/1024)
}

type GWallet struct {
	db *leveldb.DB
}

func (gw *GWallet) Get(address common.Address) (w *Wallet, err error) {
	data, err := gw.db.Get(address[:], nil)
	w = &Wallet{}
	w.toWallet(data)
	return w, err
}

func (gw *GWallet) Put(w *Wallet) error {
	return gw.db.Put(w.ID[:], w.Bytes(), nil)
}

func NewGWallet(path string) (*GWallet, error) {
	db, err := leveldb.OpenFile(path, nil)
	return &GWallet{db}, err
}

func (gw *GWallet) ReleaseGWallet() {
	gw.db.Close()
}

func CreateOneFG(data []byte) FG_data {
	big := &big.Int{}
	big.SetBytes(data)
	bitlen := FG_BIT_SIZE//len(data) * 8
	big.SetBit(big, bitlen, big.Bit(bitlen-1)^1)

	perbit := big.Bit(0)
	percount := int64(1)

	var buf bytes.Buffer
	ch := []string{"0", "1"}
	for i := 1; i < bitlen+1; i++ {
		if perbit == big.Bit(i) {
			percount++
		} else {
			buf.WriteString(ch[perbit])
			if percount > 1 {
				hexbyte := []byte(strconv.FormatInt(percount, 16))
				for i, v := range hexbyte {
					if v == '1' {
						hexbyte[i] = 'h'
					}
					if v == '0' {
						hexbyte[i] = 'g'
					}
				}
				buf.WriteString(string(hexbyte))
			}
			perbit = big.Bit(i)
			percount = 1
		}
	}

	return FG_data{buf.Bytes()}
}

func CreateOneWallet(addr common.Address, data []byte) *Wallet {
	w := &Wallet{ID: addr}
	/*
		for i:=FG_LENGTH-1; i<FG_LENGTH; i++ {
			begin := FG_BYTE_SIZE*i
			w.Data[i] = CreateOneFG(data[begin:begin+FG_BYTE_SIZE])
		}
	*/
	FT_SIZE := 4
	PT_LENGTH := FG_LENGTH / FT_SIZE
	fgCH := make(chan int, FT_SIZE)
	for i := 0; i < FT_SIZE; i++ {
		go CreatePartWallet(i*PT_LENGTH, i*PT_LENGTH+PT_LENGTH, data, w, fgCH)
	}

	for i := 0; i < FT_SIZE; i++ {
		<-fgCH
	}

	return w
}

func CreatePartWallet(start int, end int, data []byte, w *Wallet, fgCH chan int) {
	for i := start; i < end; i++ {
		begin := FG_BYTE_SIZE * i
		w.Data[i] = CreateOneFG(data[begin : begin+FG_BYTE_SIZE])
	}
	fgCH <- 1
}

func RandData(length int64, data []byte) {

	src := rand.NewSource(time.Now().Unix() + int64(0))
	rd := rand.New(src)

	for j := 0; j < int(length>>22); j++ {
		pos := rd.Int63n(length)
		bytepos := pos >> 3
		bitpos := pos % 8
		data[bytepos] |= (1 << uint(bitpos))
	}
}

/*
func merge(left, right []byte) []byte {
	conv := map[byte]int {
		'0' : 0,
		'1' : 1,
		'2' : 2,
		'3' : 3,
		'4' : 4,
		'5' : 5,
		'6' : 6,
		'7' : 7,
		'8' : 8,
		'9' : 9,
		'a' : 10,
		'b' : 11,
		'c' : 13,
		'd' : 14,
		'e' : 15,
		'f' : 16,
		'g' : 0,
		'h' : 1,
	}

	l := len(left)
	r := len(right)

	pos := l-1
	for pos > 0 {
		if(left[pos] == '1') || (left[pos] == '0') {
			break
		}
		pos--
	}

	var sum int
	sum = 0
	for i:=pos; i<l; i++ {
		sum = sum*16 + conv[left[i]]
	}
	llen := sum

	sum = 0
	for i := 1; i<r && right[i]!='1' && right[i]!='0'; i++ {
		sum += sum*16 + conv[right[i]]
	}



	if left[l] == right[0] {

	}
}

func divid_conquer(){

}
*/
func main() {
	/*
		addr := common.StringToAddress("123")

		data := make([]byte, CONBIN_BYTE_SIZE)
		//RandData(CONBIN_BIT_SIZE, data)

		w := CreateOneWallet(addr, data)

		w.Statistics()

	*/




	gw, err := NewGWallet("path/to/db")
	if err != nil {
		fmt.Println(err)
	}
	defer gw.ReleaseGWallet()

	addr := common.StringToAddress("1")
	w, err := gw.Get(addr)

	for i:=0; i<1<<20; i++ {
		addr = common.StringToAddress(strconv.FormatInt(int64(i), 16))
		w, err = gw.Get(addr)
		w.Set(int64(i))
		gw.Put(w)

		fmt.Println(i)

	}

	fmt.Printf("%v %s", w.ID, w.Data)

}
