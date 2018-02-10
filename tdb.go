package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"go-unitcoin/libraries/common"
	"bytes"
	"encoding/gob"
	"strconv"
	"math/big"
	//"compress/zlib"
	"time"
	"math/rand"
)

const (
	CONBIN_BIT_SIZE		= (1<<32)
	CONBIN_BYTE_SIZE	= (CONBIN_BIT_SIZE>>3)

	FG_LENGTH 			= (1<<10)
	FG_BIT_SIZE 		= (CONBIN_BIT_SIZE/FG_LENGTH)
	FG_BYTE_SIZE		= (FG_BIT_SIZE>>3)
)

type FG_data struct {
//	Len int
	Data []byte
}

type Wallet struct {
	ID	common.Address
	Data [FG_LENGTH]FG_data
}

func (w *Wallet) Bytes() ([]byte) {
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
	if err!=nil {
		fmt.Println(err)
	}
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

func (gw *GWallet) Put (w *Wallet) error {
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
	bitlen := len(data)*8
	big.SetBit(big, bitlen, big.Bit(bitlen-1)^1)

	perbit := big.Bit(0)
	percount := int64(1)

	var buf bytes.Buffer
	ch := []string{"0", "1"}
	for i:=1; i<bitlen+1; i++ {
		if perbit == big.Bit(i) {
			percount++
		}else{
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
/*
	b := bytes.NewBuffer([]byte{})
	w := zlib.NewWriter(b)
	w.Write(buf.Bytes())
	w.Close()
*/
	hexdata := buf.Bytes()
	return FG_data{hexdata}
}

func CreateOneWallet(addr common.Address, data []byte) *Wallet {
	w := &Wallet{ID:addr}
	/*
	for i:=FG_LENGTH-1; i<FG_LENGTH; i++ {
		begin := FG_BYTE_SIZE*i
		w.Data[i] = CreateOneFG(data[begin:begin+FG_BYTE_SIZE])
	}
	*/
	FT_SIZE := 8
	PT_LENGTH := FG_LENGTH/FT_SIZE
	fgCH := make(chan int, FT_SIZE)
	for i:=0; i<FT_SIZE; i++ {
		go CreatePartWallet(i*PT_LENGTH, i*PT_LENGTH+PT_LENGTH, data, w, fgCH)
	}

	for i:=0; i<FT_SIZE; i++ {
		<-fgCH
	}

	return w
}

func CreatePartWallet(start int, end int, data []byte, w *Wallet, fgCH chan int) {
	for i:=start; i<end; i++ {
		begin := FG_BYTE_SIZE*i
		w.Data[i] = CreateOneFG(data[begin:begin+FG_BYTE_SIZE])
	}
	fgCH<-1
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

func main() {
	addr := common.StringToAddress("123")
	data := make([]byte, CONBIN_BYTE_SIZE)
	RandData(CONBIN_BIT_SIZE, data)

	/*
	hd := data[CONBIN_BYTE_SIZE-FG_BYTE_SIZE:]
	for _, v:=range hd {
		fmt.Printf("%08b", v)
	}
	*/

	w := CreateOneWallet(addr, data)

	fg := w.Data[FG_LENGTH-1]
	/*
	var b bytes.Buffer
	b.Write(fg.Data)
	r,_ := zlib.NewReader(&b)
	hexdata := make([]byte, 1024)
	len, _ := r.Read(hexdata)
	*/
	hexdata := fg.Data
	fmt.Printf("\n\nlen:%d data:%s", len(hexdata), hexdata)
/*
	gw, err := NewGWallet("path/to/db")
	if(err != nil) {
		fmt.Println(err)
	}
	defer gw.ReleaseGWallet()

	gw.Put(w)
	gwt, err := gw.Get(addr)
	fmt.Println(gwt)
*/
}