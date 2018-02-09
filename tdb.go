package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"

	//"github.com/syndtr/goleveldb/leveldb/util"
	"go-unitcoin/libraries/common"
	"bytes"
	"encoding/gob"
	"strconv"
	"math/big"
	//"compress/zlib"
	//"os"
	//"io"
)

const (
	FG_LENGTH = 1
)

type FG_data struct {
	Len int
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

func CreteOneFG(data []byte) FG_data {
	big := &big.Int{}
	big.SetBytes(data)
	bitlen := len(data)*8
	big.SetBit(big, bitlen, big.Bit(bitlen-1)^1)

	perbit := big.Bit(0)
	percount := int64(1)
	fmt.Print(big.Bit(0))

	var buf bytes.Buffer
	ch := []string{"0", "1"}
	for i:=1; i<bitlen+1; i++ {
		fmt.Print(big.Bit(i))
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

	hexdata := buf.Bytes()
	hexlen := len(hexdata)
	fg := FG_data{hexlen, hexdata}
	fmt.Printf("\nlen:%d data:%s\n", hexlen, hexdata)
/*
	buf.Reset()
	w := zlib.NewWriter(&buf)
	w.Write(hexdata)
	w.Close()
	r,_ := zlib.NewReader(&buf)
	io.Copy(os.Stdout, r)
*/
	return fg
}

func main() {
	data := []byte{1,2,3,4,5,6,7,8,9,5,6,4,3,2,3,2,2,5,7,3,2,1}
	fmt.Printf("\n%08b\n", data)
	fg := CreteOneFG(data)
	fmt.Printf("\nlen:%d data:%s", fg.Len, fg.Data)
	/*
	gw, err := NewGWallet("path/to/db")
	if(err != nil) {
		fmt.Println(err)
	}
	defer gw.ReleaseGWallet()

	addr := common.StringToAddress("123456")
	wt := &Wallet{ID:addr}

	for i, v := range wt.Data {
		v.Len = i
		v.Data = []byte(strconv.FormatInt(int64(i), 10))
		wt.Data[i] = v
	}

	fmt.Println(wt)
	//fmt.Println("conten:",wt.Bytes())
	fmt.Println("conten len:", len(wt.Bytes()))

	gw.Put(wt)
	gwt, err := gw.Get(addr)
	fmt.Println(gwt)
	*/
}