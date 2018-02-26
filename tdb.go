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
	"errors"
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
	Data []byte
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
/*
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
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
		'g': 0,
		'h': 1,
	}

	var i, k int64

	data := w.Data[:]
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
						return sum, length
					}
					sum = sum*16 + conv[ch]
					length++
				}
				return sum, length
			}(data[i+1:])

			i += l
			if s == 0 {
				s = 1
			}
			for k = 0; k < s; k++ {
				big.SetBit(big, bigpos, uint(bit))
				bigpos++
			}
		}
	}
	big.SetBit(big, int(bitpos), 1)
	//w.Data = CreateOneFG(big.Bytes())
}
*/
func (w *Wallet) SetBit(bitpos int64, bit int) {
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
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
		'g': 0,
		'h': 1,
	}
	convint := func(num int64) string {
		if num == 0 || num == 1{
			return string("")
		}
		hexbyte := []byte(strconv.FormatInt(num, 16))
		for i, v := range hexbyte {
			if v == '1' {
				hexbyte[i] = 'h'
			}
			if v == '0' {
				hexbyte[i] = 'g'
			}
		}
		return string(hexbyte)
	}

	wlength := len(w.Data)
	org := w.Data[:]
	type POS struct{
		val int64
		pos int64
		len int64
		off int64
		offlen int64
	}

	per, cur, next, current := POS{}, POS{}, POS{}, POS{}

	var found bool
	for i:=0; i<wlength; {
		ch := w.Data[i]
		val := conv[ch]
		if ch == '1' || ch == '0' {
			if (found == false) {
				per = cur
			}
			cur.pos += cur.len
			cur.len = 0
			cur.val = val
			cur.off = int64(i)
			cur.offlen = 0

			s, l := func(data []byte) (sum int64, length int64) {
				sum = 0
				length = 1
				for i := 0; i < len(data); i++ {
					ch := data[i]
					if ch == '1' || ch == '0' {
						return sum, length
					}
					sum = sum*16 + conv[ch]
					length++
				}
				return sum, length
			}(w.Data[i+1:])

			if s==0 {
				s=1
			}
			cur.len = s
			cur.offlen += l

			if(found) {
				next = cur
				break
			}

			i += int(l)
			if bitpos >=cur.pos && bitpos<(cur.pos+cur.len) {
				found = true
				current = cur
			}
		}
	}

	if current.val == int64(bit) {
		return
	}

	//fmt.Println(per, current, next)

	cur = current
	var end int64
	if(next.off > 0){
		end = next.off+next.offlen
	}else{
		end = cur.off+cur.offlen
	}

	if bitpos == current.pos {
		per.len += 1
		per.val = cur.val^1
		cur.len -= 1
	}else if bitpos == current.pos+current.len-1 {
		next.len += 1
		next.val = cur.val^1
		cur.len -= 1
	}else {
		end = cur.off+cur.offlen
		per.len = bitpos - cur.pos
		per.val = cur.val
		per.off = cur.off

		next.len = current.pos+current.len - bitpos -1
		next.val = cur.val
		next.off = 0

		cur.val = cur.val^1
		cur.len = 1
	}

	if( cur.len == 0){
		per.len += next.len
		next.len = 0
	}

	zo := []string{"0", "1"}
	data := []byte{}

	if (per.len > 0) {
		b := bytes.Buffer{}
		data = append(data, org[0:per.off]...)
		b.WriteString(zo[per.val])
		b.WriteString(convint(per.len))
		data = append(data, b.Bytes()...)
	}

	if(cur.len > 0){
		b := bytes.Buffer{}
		b.WriteString(zo[cur.val])
		b.WriteString(convint(cur.len))
		data = append(data, b.Bytes()...)
	}

	if(next.len > 0){
		b := bytes.Buffer{}
		b.WriteString(zo[next.val])
		b.WriteString(convint(next.len))
		data = append(data, b.Bytes()...)
	}

	data = append(data, org[end:]...)

	w.Data = data
	//fmt.Printf("%s\n", data)
}

func (w *Wallet) Bit(bitpos int64) int {
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
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
		'g': 0,
		'h': 1,
	}

	wlength := len(w.Data)

	type POS struct{
		val int64
		pos int64
		len int64
	}

	cur := POS{}

	for i:=0; i<wlength; {
		ch := w.Data[i]
		val := conv[ch]
		if ch == '1' || ch == '0' {
			cur.pos += cur.len
			cur.len = 0
			cur.val = val

			s, l := func(data []byte) (sum int64, length int64) {
				sum = 0
				length = 1
				for i := 0; i < len(data); i++ {
					ch := data[i]
					if ch == '1' || ch == '0' {
						return sum, length
					}
					sum = sum*16 + conv[ch]
					length++
				}
				return sum, length
			}(w.Data[i+1:])

			if s==0 {
				s=1
			}
			cur.len = s

			i += int(l)
			if bitpos >=cur.pos && bitpos<(cur.pos+cur.len) {
				return int(cur.val)
			}
		}
	}
	return int(cur.val)
}

func (w *Wallet) FGCompress() {

		b := bytes.NewBuffer([]byte{})
		wr := zlib.NewWriter(b)
		wr.Write(w.Data)
		wr.Close()
		w.Data = b.Bytes()
}

func (w *Wallet) Compress() []byte {
	b := bytes.NewBuffer([]byte{})
	W := zlib.NewWriter(b)

	W.Write(w.Data)

	W.Close()
	return b.Bytes()
}

func (w *Wallet) Statistics() {

	var sum int64
	b := bytes.NewBuffer([]byte{})
	W := zlib.NewWriter(b)


		sum += int64(len(w.Data))
		W.Write(w.Data)

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
	w = &Wallet{address, data}
	return w, err
}

func (gw *GWallet) Put(w *Wallet) error {
	return gw.db.Put(w.ID[:], w.Data, nil)
}

func (gw *GWallet) Move(dst, src common.Address, bin int) error {
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
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
		'g': 0,
		'h': 1,
	}
	var wd, ws *Wallet
	var err error
	if wd, err = gw.Get(dst); err != nil {
		return err
	}

	if ws, err = gw.Get(src); err != nil {
		return err
	}

	dlen := len(ws.Data)
	var pos int64
	var spos []int64
	for i:=0; i<dlen; {
		ch := ws.Data[i]
		if ch == '0' || ch == '1'{
			s, l := func(data []byte) (sum int64, length int64) {
				sum = 0
				length = 1
				for i := 0; i < len(data); i++ {
					ch := data[i]
					if ch == '1' || ch == '0' {
						return sum, length
					}
					sum = sum*16 + conv[ch]
					length++
				}
				return sum, length
			}(ws.Data[i+1:])
			if s == 0 {
				s = 1
			}
			i += int(l)

			if ch == '0' {
				pos += s
			}else{
				for s > 0 && bin > 0 {
					spos = append(spos, pos)
					pos++
					s--
					bin--
				}
				if bin <= 0 {
					break
				}
			}
		}
	}

	if bin > 0 {
		return errors.New("Not Enough ubin")
	}

    fmt.Printf("%v %s\n", ws.ID, ws.Data)
    fmt.Printf("%v %s\n", wd.ID, wd.Data)
	for _, v := range spos {
		ws.SetBit(v, 0)
		wd.SetBit(v, 1)
	}


    fmt.Printf("%v %s\n", ws.ID, ws.Data)
    fmt.Printf("%v %s\n", wd.ID, wd.Data)
	return nil
}

func NewGWallet(path string) (*GWallet, error) {
	db, err := leveldb.OpenFile(path, nil)
	return &GWallet{db}, err
}

func (gw *GWallet) ReleaseGWallet() {
	gw.db.Close()
}

func CreateOneFG(data []byte) []byte {
	big := &big.Int{}
	big.SetBytes(data)

	bitlen := len(data) << 3
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

	return buf.Bytes()
}

func CreateOneWallet(addr common.Address, data []byte) *Wallet {
	w := &Wallet{ID: addr}
	w.Data = CreateOneFG(data[:])
	/*
	FT_SIZE := 4
	PT_LENGTH := FG_LENGTH / FT_SIZE
	fgCH := make(chan int, FT_SIZE)
	for i := 0; i < FT_SIZE; i++ {
		go CreatePartWallet(i*PT_LENGTH, i*PT_LENGTH+PT_LENGTH, data, w, fgCH)
	}

	for i := 0; i < FT_SIZE; i++ {
		<-fgCH
	}
	*/
	return w
}
/*
func CreatePartWallet(start int, end int, data []byte, w *Wallet, fgCH chan int) {
	for i := start; i < end; i++ {
		begin := FG_BYTE_SIZE * i
		CreateOneFG(data[begin : begin+FG_BYTE_SIZE])
	}
	fgCH <- 1
}
*/
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


func merge(left, right []byte) []byte {
	conv := map[byte]int64 {
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
		'c' : 12,
		'd' : 13,
		'e' : 14,
		'f' : 15,
		'g' : 0,
		'h' : 1,
	}
	convint := func(num int64) string {
		if num == 0 || num == 1{
			return string("")
		}
		hexbyte := []byte(strconv.FormatInt(num, 16))
		for i, v := range hexbyte {
			if v == '1' {
				hexbyte[i] = 'h'
			}
			if v == '0' {
				hexbyte[i] = 'g'
			}
		}
		return string(hexbyte)
	}

	var sum, i, length, r, l, pos int64
	l = int64(len(left))
	r = int64(len(right))

	pos = l-1
	for pos > 0 {
		if(left[pos] == '1') || (left[pos] == '0') {
			break
		}
		pos--
	}

	sum = 0
	for i=pos+1; i<l; i++ {
		sum = sum*16 + conv[left[i]]
	}
	if(sum == 0){
		sum = 1
	}
	llen := sum

	sum = 0
	for i = 1; i<r && right[i]!='1' && right[i]!='0'; i++ {
		sum = sum*16 + conv[right[i]]
	}
	if(sum ==0 ){
		sum = 1
	}

	data := []byte{}
	if left[pos] == right[0] {
		length = llen + sum
		data = append(data, left[0:pos]...)
		b := bytes.Buffer{}
		if right[0] == '0'{
			b.WriteString("0")
		}else{
			b.WriteString("1")
		}
		b.WriteString(convint(length))
		data = append(data, b.Bytes()...)
		data = append(data, right[i:]...)
	}else{
		data = append(data, left[:]...)
		data = append(data, right[:]...)
	}
	return data
}

func divid_conquer(data []byte) []byte{
	if(len(data) > 1){
		mid  := len(data)/2
		left := divid_conquer(data[0:mid])
		right:= divid_conquer(data[mid:len(data)])
		return merge(left, right)
	}
	dt := int(data[0])
	return []byte(strconv.FormatInt(int64(dt), 16))
}

func main() {

	/*
	data := make([]byte, 1<<10)
	result := divid_conquer(data)
	fmt.Printf("%s\n", result)
	*/

	/*
	addr := common.StringToAddress("0")
	data := make([]byte, CONBIN_BYTE_SIZE)
	//RandData(CONBIN_BIT_SIZE, data)
	w := CreateOneWallet(addr, data)
	w.Statistics()
	*/
	addr := common.StringToAddress("0")
	//w := &Wallet{addr, []byte("0hgggggggg")}


	gw, err := NewGWallet("path/to/db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer gw.ReleaseGWallet()


	//w, _ = gw.Get(addr)
	/*
	fmt.Printf("%v %s\n", w.ID, w.Data)
	w.SetBit(0,1)
	w.SetBit(1,1)
	w.SetBit(2,1)
	w.SetBit(3,1)
	w.SetBit(4,1)
	w.SetBit(5,1)
	fmt.Printf("%v %s\n", w.ID, w.Data)
	gw.Put(w)
	*/
    addr1 := common.StringToAddress("1")
	gw.Move(addr1, addr, 5)
	fmt.Println("---------")
	gw.Move(addr1, addr, 9)

/*
	var i  int64
	for i=0; i<1<<20; i++ {
		addr := common.StringToAddress(strconv.FormatInt(i, 16))
		w.ID = addr
		gw.Put(w)

		if(i%8196 == 0) {
			fmt.Printf("%d %v %s\n", i, w.ID, w.Data)
		}
	}
*/
}
