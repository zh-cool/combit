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

	"errors"
	"go-unitcoin/libraries/crypto/sha3"
	"go-unitcoin/libraries/db/lvldb"
	"go-unitcoin/libraries/event"
	"encoding/binary"
)

var (
	conv = map[byte]uint64{
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
)

const (
	CONBIN_BIT_SIZE  = (1 << 32)
	CONBIN_BYTE_SIZE = (CONBIN_BIT_SIZE >> 3)

	FG_LENGTH    = (1 << 10)
	FG_BIT_SIZE  = (CONBIN_BIT_SIZE / FG_LENGTH)
	FG_BYTE_SIZE = (FG_BIT_SIZE >> 3)
)

type Req_verify struct {
	From   common.Address
	To     common.Address
	Number uint64
	TxHash common.Hash
}

type FG_data struct {
	//	Len int
	Data []byte
}

type Wallet struct {
	Addr common.Address
	Sid  uint64
	id   []byte
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
func (w *Wallet) SetBit(bitpos uint64, bit uint64) {
	convint := func(num uint64) string {
		if num == 0 || num == 1 {
			return string("")
		}
		hexbyte := []byte(strconv.FormatUint(num, 16))
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
	type POS struct {
		val    uint64
		pos    uint64
		len    uint64
		off    uint64
		offlen uint64
	}

	per, cur, next, current := POS{}, POS{}, POS{}, POS{}

	var found bool
	for i := 0; i < wlength; {
		ch := w.Data[i]
		val := conv[ch]
		if ch == '1' || ch == '0' {
			if found == false {
				per = cur
			}
			cur.pos += cur.len
			cur.len = 0
			cur.val = val
			cur.off = uint64(i)
			cur.offlen = 0

			s, l := func(data []byte) (sum uint64, length uint64) {
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

			if s == 0 {
				s = 1
			}
			cur.len = s
			cur.offlen += l

			if found {
				next = cur
				break
			}

			i += int(l)
			if bitpos >= cur.pos && bitpos < (cur.pos+cur.len) {
				found = true
				current = cur
			}
		}
	}

	if current.val == uint64(bit) {
		return
	}

	//fmt.Println(per, current, next)

	cur = current
	var end uint64
	if next.off > 0 {
		end = next.off + next.offlen
	} else {
		end = cur.off + cur.offlen
	}

	if bitpos == current.pos {
		per.len += 1
		per.val = cur.val ^ 1
		cur.len -= 1
	} else if bitpos == current.pos+current.len-1 {
		next.len += 1
		next.val = cur.val ^ 1
		cur.len -= 1
	} else {
		end = cur.off + cur.offlen
		per.len = bitpos - cur.pos
		per.val = cur.val
		per.off = cur.off

		next.len = current.pos + current.len - bitpos - 1
		next.val = cur.val
		next.off = 0

		cur.val = cur.val ^ 1
		cur.len = 1
	}

	if cur.len == 0 {
		per.len += next.len
		next.len = 0
	}

	zo := []string{"0", "1"}
	data := []byte{}

	if per.len > 0 {
		b := bytes.Buffer{}
		data = append(data, org[0:per.off]...)
		b.WriteString(zo[per.val])
		b.WriteString(convint(per.len))
		data = append(data, b.Bytes()...)
	}

	if cur.len > 0 {
		b := bytes.Buffer{}
		b.WriteString(zo[cur.val])
		b.WriteString(convint(cur.len))
		data = append(data, b.Bytes()...)
	}

	if next.len > 0 {
		b := bytes.Buffer{}
		b.WriteString(zo[next.val])
		b.WriteString(convint(next.len))
		data = append(data, b.Bytes()...)
	}

	data = append(data, org[end:]...)

	w.Data = data
	//fmt.Printf("%s\n", data)
}

func (w *Wallet) Bit(bitpos uint64) int {
	wlength := len(w.Data)

	type POS struct {
		val uint64
		pos uint64
		len uint64
	}

	cur := POS{}

	for i := 0; i < wlength; {
		ch := w.Data[i]
		val := conv[ch]
		if ch == '1' || ch == '0' {
			cur.pos += cur.len
			cur.len = 0
			cur.val = val

			s, l := func(data []byte) (sum uint64, length uint64) {
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

			if s == 0 {
				s = 1
			}
			cur.len = s

			i += int(l)
			if bitpos >= cur.pos && bitpos < (cur.pos+cur.len) {
				return int(cur.val)
			}
		}
	}
	return int(cur.val)
}

func (w *Wallet) GetPos(bin uint64) []uint64 {
	dlen := len(w.Data)
	var pos uint64
	var spos []uint64

	for i := 0; i < dlen; {
		ch := w.Data[i]
		if ch == '0' || ch == '1' {
			s, l := func(data []byte) (sum uint64, length uint64) {
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
			if s == 0 {
				s = 1
			}
			i += int(l)

			if ch == '0' {
				pos += s
			} else {
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
	return spos
}

func (w *Wallet) GetGPos(bin uint64, gid uint64) []uint64 {
	dlen := len(w.Data)
	var pos uint64
	var spos []uint64

	begin := gid * FG_BIT_SIZE
	end := begin + FG_BIT_SIZE

	for i := 0; i < dlen; {
		ch := w.Data[i]
		if ch == '0' || ch == '1' {
			s, l := func(data []byte) (sum uint64, length uint64) {
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
			if s == 0 {
				s = 1
			}
			i += int(l)

			if ch == '0' {
				pos += s
			} else {
				if pos+s < begin {
					pos += s
					continue
				}
				for s > 0 && bin > 0 {
					if pos >= begin && pos < end {
						spos = append(spos, pos)
					}
					pos++
					s--
					bin--
				}
				if bin <= 0 {
					break
				}
			}
			if pos >= end {
				break
			}
		}
	}
	return spos
}

func (w *Wallet) GetBalance() uint64 {
	dlen := len(w.Data)
	var balance uint64

	for i := 0; i < dlen; {
		ch := w.Data[i]
		if ch == '0' || ch == '1' {
			s, l := func(data []byte) (sum uint64, length uint64) {
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
			if s == 0 {
				s = 1
			}
			i += int(l)

			if ch == '1' {
				balance += s
			}
		}
	}
	return balance
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
	db      *lvldb.LDBDatabase
	posCh   event.Feed
	PosINfo []common.Address
	hashTree 	[]common.Hash

	GID   uint64
	ROOT  common.Hash
	GHASH []common.Hash
	GADDR []common.Address
}

type ReqPosINfo struct {
	Pos []byte
}
type ReqPosINfoEvent struct{ Pos *ReqPosINfo }

type ResPosInfo struct {
	Pos  []byte
	Addr common.Address
}
type ResPosINfoList []*ResPosInfo

func (gw *GWallet) Hash() common.Hash {
	return gw.ROOT
}

func (gw *GWallet) CHash() common.Hash {
	var i uint64
	i = uint64(len(gw.hashTree)/2 - 1)
	for ; i >= 1; i-- {
		left := i << 1
		right := i<<1 + 1
		hw := sha3.NewKeccak256()
		hw.Write(gw.hashTree[left][:])
		hw.Write(gw.hashTree[right][:])
		hw.Sum(gw.hashTree[i][:0])
	}
	gw.ROOT = gw.hashTree[1]
	return gw.hashTree[1]
}

func (gw *GWallet) uHash(pos uint64) {
	begin := gw.GID*FG_BIT_SIZE
	end := begin+FG_BIT_SIZE
	if pos<begin || pos>=end {
		return
	}

	var parent uint64
	parent = pos <<1
	for ; parent >= 1;  {
		left := parent >> 1
		right := parent >> 1 + 1
		hw := sha3.NewKeccak256()
		hw.Write(gw.hashTree[left][:])
		hw.Write(gw.hashTree[right][:])
		hw.Sum(gw.hashTree[parent][:0])
		parent = parent << 1
	}
}

func (gw *GWallet) IDtoAddr(id uint64) (common.Address, error) {
	b := new(big.Int).SetUint64(id)
	addr, err := gw.db.Get(b.SetBit(b, 64, 1).Bytes()[1:])
	if err != nil {
		return common.Address{}, err
	}

	return common.BytesToAddress(addr), nil
}

func (gw *GWallet) Get(address common.Address) (w *Wallet, err error) {
	data, err := gw.db.Get(address[:])
	if err != nil {
		return nil, err
	}
	sid := new(big.Int).SetBytes(data[:8])

	w = &Wallet{address, sid.Uint64(), data[:8], data[8:]}
	return w, err
}

func (gw *GWallet) Put(w *Wallet) error {
	if _, err := gw.db.Get(w.id); err != nil {
		gw.db.Put(w.id, w.Addr[:])
	}
	return gw.db.Put(w.Addr[:], append(w.id, w.Data...))
}

func (gw *GWallet) Move(dst, src common.Address, bin uint64) error {

	var wd, ws *Wallet
	var err error
	if wd, err = gw.Get(dst); err != nil {
		return err
	}

	if ws, err = gw.Get(src); err != nil {
		return err
	}

	if dst==src {
		return nil
	}

	dlen := len(ws.Data)
	var pos uint64
	var spos []uint64
	for i := 0; i < dlen; {
		ch := ws.Data[i]
		if ch == '0' || ch == '1' {
			s, l := func(data []byte) (sum uint64, length uint64) {
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
			} else {
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

	for _, v := range spos {
		ws.SetBit(v, 0)
		wd.SetBit(v, 1)
		gw.uHash(v)
	}

	return nil
}

func (gw *GWallet) RequestLostPos() {
	var idx, rp, offset uint64
	var pos []byte

	addr := common.Address{}
	offset = gw.GID*FG_BIT_SIZE

	for idx = 0; idx <= FG_BIT_SIZE; idx++ {
		if idx != FG_BIT_SIZE && gw.PosINfo[idx] == addr {
			rp++
		} else {
			if rp == 1 {
				b := new(big.Int).SetUint64(offset + idx - 1)
				b = b.SetBit(b, 64, 1)
				pos = append(pos, b.Bytes()[1:]...)
			} else if rp > 1 {
				begin := idx - rp
				end := idx - 1

				b := new(big.Int).SetUint64(offset + begin)
				b = b.SetBit(b, 64, 1)
				pos = append(pos, b.Bytes()[1:]...)

				pos = append(pos, '-')

				b = new(big.Int).SetUint64(offset + end)
				b = b.SetBit(b, 64, 1)
				pos = append(pos, b.Bytes()[1:]...)
			}
			rp = 0
		}
	}
	if len(pos) == 0 {
		return
	}

	req := ReqPosINfo{
		Pos: pos,
	}
	gw.posCh.Send(ReqPosINfoEvent{&req})
}

func (gw *GWallet) AddPosInfo(pos ResPosINfoList) {
	for _, v := range pos {
		offset := 0
		var pre uint64 = 0

		for offset < len(v.Pos) {
			cur, flag := func(coin []byte) (idx uint64, flag bool) {
				if coin[0] == '-' {
					offset += 9
					return new(big.Int).SetBytes(coin[1:9]).Uint64(), true
				}
				offset += 8
				return new(big.Int).SetBytes(coin[0:8]).Uint64(), false

			}(v.Pos[offset:])

			var i uint64
			if flag == true {
				length := cur - pre
				for i = 0; i < length; i++ {
					pre++
					gw.PosINfo[pre] = v.Addr
				}
			} else {
				pre = cur
				gw.PosINfo[cur] = v.Addr
			}
		}
	}
	gw.RequestLostPos()
}

func (gw *GWallet) GetPosInfo(req *ReqPosINfo) ResPosINfoList {
	pos := req.Pos
	posinfo := make(map[common.Address][]uint64)
	empty := common.Address{}

	offset := 0
	var pre uint64 = 0

	for offset < len(pos) {
		cur, flag := func(coin []byte) (idx uint64, flag bool) {
			if coin[0] == '-' {
				offset += 9
				return new(big.Int).SetBytes(coin[1:9]).Uint64(), true
			}
			offset += 8
			return new(big.Int).SetBytes(coin[0:8]).Uint64(), false
		}(pos[offset:])

		var i uint64
		if flag == true {
			length := cur - pre
			for i = 0; i < length; i++ {
				pre++
				addr := gw.PosINfo[pre]
				if addr == empty{
					continue
				}
				_, ok := posinfo[addr]
				if ok == true {
					posinfo[addr] = append(posinfo[addr], pre)
				} else {
					posinfo[addr] = []uint64{pre}
				}
			}
		} else {
			addr := gw.PosINfo[cur]
			if addr == empty{
				continue
			}
			_, ok := posinfo[addr]
			if ok == true {
				posinfo[addr] = append(posinfo[addr], cur)
			} else {
				posinfo[addr] = []uint64{cur}
			}
			pre = cur
		}
	}

	var poslist ResPosINfoList
	for addr, idx := range posinfo {
		spos := []byte{}

		for _, v := range idx {
			b := new(big.Int)
			b.SetUint64(v).SetBit(b, 64, 1)
			spos = append(spos, b.Bytes()[1:]...)
		}

		v := &ResPosInfo{
			Pos:  spos,
			Addr: addr,
		}
		poslist = append(poslist, v)
	}

	return poslist
}

func (gw *GWallet) SubscribeReqPosINfoEvent(posCh chan ReqPosINfoEvent) {
	gw.posCh.Subscribe(posCh)
}

func (gw *GWallet) AddWallets(wallets []*Wallet) {
	leafe := len(gw.hashTree)
	for _, v := range wallets {
		gw.Put(v)
		var idx uint64 = gw.GID*FG_BIT_SIZE

		for i:=0; i<FG_BIT_SIZE; i++ {
			if v.Bit(idx)==1 {
				gw.PosINfo[i] = v.Addr
				hw := sha3.NewKeccak256()
				hw.Write(v.Data)
				hw.Sum(gw.hashTree[leafe+i][:0])
				gw.uHash(idx)
				idx++
			}
		}
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, gw.PosINfo)
	gw.db.Put([]byte("POSINFO"), buf.Bytes())
}

/*
func (gw *GWallet) Update(block *protocol.Block) {
	txs := block.Transactions()

	for _, tx := range txs {
		signer := protocol.MakeSigner()
		from, _:= protocol.Sender(signer, tx)
		value := tx.Value()
		to := tx.To()
		gw.Move(*to, from, value.Uint64())
	}
}

func (gw *GWallet) CheckTxBin() {
	for {
		select {
		case event := <-gw.txCH:
			fmt.Println("austin:", event)
			from, to, num := func(tx *protocol.Tx) (common.Address, common.Address, uint64) {
				signer := protocol.MakeSigner()
				from, _ := protocol.Sender(signer, tx)
				value := tx.Value()
				fmt.Printf("TX Hash:%x\n",tx.Hash())
				return from, *tx.To(), value.Uint64()

			}(event.Transaction)
			gw.Verify <- Req_verify{from, to, num, event.Transaction.Hash()}
		}
	}
}

func (gw *GWallet) Worker() {
	for {
		select {
		case v := <-gw.Verify:
			w, _ := gw.Get(v.From)
			pos := w.GetPos(v.Number)
			gw.Move(v.To, v.From, v.Number)
			rs := &Res_verify{&gw.ROOT, gw.GID, pos, &v.TxHash, []byte{}}
			rs.Sign(gw.prv)
			fmt.Println(pos)
		}
	}
}

func (gw *GWallet) Sync() {
	forceSync := time.NewTicker(10 * time.Second )
	defer forceSync.Stop()

	for {
		select {
		case <-forceSync.C:
			cur := gw.blockchain.CurrentBlock().GetNumberU64()
			for ; gw.blockNum < cur+1; gw.blockNum++ {
				block := gw.blockchain.GetBlockByNumber(gw.blockNum)
				gw.Update(block)
			}
		}
	}
}

func (gw *GWallet) Start() {
    gw.txCH = make(chan protocol.TxPreEvent)
	gw.txpool.SubscribeTxPreEvent(gw.txCH)
	gw.Verify = make(chan Req_verify)
	fmt.Println("Start")
	go gw.CheckTxBin()
	go gw.Worker()
	go gw.Sync()
}
*/

func NewGWallet(db lvldb.Database) *GWallet {
	gw := &GWallet{db: db.(*lvldb.LDBDatabase)}
	posInfo, err := gw.db.Get([]byte("POSINFO"))
	if err!= nil {
		return gw
	}

	buf := bytes.NewBuffer(posInfo)
	gw.PosINfo = make([]common.Address, FG_BIT_SIZE)
	binary.Read(buf, binary.BigEndian, gw.PosINfo[:FG_BIT_SIZE])

	gw.hashTree = make([]common.Hash, FG_BIT_SIZE<<1)
	leaf := len(gw.hashTree)/2
	addr := common.Address{}
	for i:=0; i<FG_BIT_SIZE; i++ {
		if addr != gw.PosINfo[i] {
			hw := sha3.NewKeccak256()
			hw.Write(gw.PosINfo[i][:])
			hw.Sum(gw.hashTree[leaf+i][:0])
		}
	}

	gw.CHash()
	return gw
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
	w := &Wallet{Addr: addr}
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
		if num == 0 || num == 1 {
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

	pos = l - 1
	for pos > 0 {
		if (left[pos] == '1') || (left[pos] == '0') {
			break
		}
		pos--
	}

	sum = 0
	for i = pos + 1; i < l; i++ {
		sum = sum*16 + conv[left[i]]
	}
	if sum == 0 {
		sum = 1
	}
	llen := sum

	sum = 0
	for i = 1; i < r && right[i] != '1' && right[i] != '0'; i++ {
		sum = sum*16 + conv[right[i]]
	}
	if sum == 0 {
		sum = 1
	}

	data := []byte{}
	if left[pos] == right[0] {
		length = llen + sum
		data = append(data, left[0:pos]...)
		b := bytes.Buffer{}
		if right[0] == '0' {
			b.WriteString("0")
		} else {
			b.WriteString("1")
		}
		b.WriteString(convint(length))
		data = append(data, b.Bytes()...)
		data = append(data, right[i:]...)
	} else {
		data = append(data, left[:]...)
		data = append(data, right[:]...)
	}
	return data
}

func divid_conquer(data []byte) []byte {
	if len(data) > 1 {
		mid := len(data) / 2
		left := divid_conquer(data[0:mid])
		right := divid_conquer(data[mid:len(data)])
		return merge(left, right)
	}
	dt := int(data[0])
	return []byte(strconv.FormatInt(int64(dt), 16))
}

func main() {

	addr := []string{"0xD339ac7AeBBCC86d601aa6Ff7e9c05B3224bBF22",
		"0x4d7405BF97a5aA98F916cCcBA0711451D8b59fb2",
		"0x473F76d3bca7ea32D26D59470B5BE30a3395abfe",
		"0xD05F8b9B1F90E73554FD5Ac724Feeb1a3C654D06",
	}
	_ = addr
	db, _ := lvldb.NewLDBDatabase("/home/austin/go/src/go-unitcoin/build/20050/unitcoin/transactionchaindata", 0, 0)
	gw := NewGWallet(db)
	defer gw.ReleaseGWallet()
/*
	gw.RequestLostPos()

	pos := []byte{}
	for _, v := range []uint64{1,5,7,9,0,15} {
		if v==0 {
			pos = append(pos, '-')
			continue
		}
		b := new(big.Int).SetUint64(v)
		b.SetBit(b, 64, 1)
		pos = append(pos, b.Bytes()[1:]...)
	}
	//fmt.Println(pos)
	list := gw.GetPosInfo(&ReqPosINfo{pos})
	//fmt.Println(list[0])
	gw.AddPosInfo(list)

	gw.ReleaseGWallet()
*/
	id := [8]byte{}
	w := &Wallet{common.Address{}, 0, id[:], []byte("0hgggggggg")}
	defer gw.ReleaseGWallet()

	var glist [1<<22]common.Address

	var i, j uint64
	for i = 0; i < 4; i++ {
		b := new(big.Int).SetUint64(i)
		b.SetBit(b, 64, 1)
		w.Addr = common.HexToAddress(addr[i])
		w.Sid = i

		w.id = b.Bytes()[1:]
		gw.Put(w)
		fmt.Printf("%x, %v, %v, %s\n", w.Addr, w.Sid, w.id, w.Data)
	}
	//fmt.Printf("%x, %v, %v, %s\n", w.Addr, w.Sid, w.id, w.Data)
	fmt.Println()
	for i = 0; i < 4; i++ {
		addr := common.HexToAddress(addr[i])
		w, _ := gw.Get(addr)
		for j=0; j<100; j++ {
			w.SetBit(i*100+j, 1)
			glist[i*100+j] = addr
		}
		gw.Put(w)
		fmt.Printf("%x, %v, %v, %s\n", w.Addr, w.Sid, w.id, w.Data)
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, glist)
	gw.db.Put([]byte("POSINFO"), buf.Bytes())

/*
	var mTree [1 << 23]common.Hash
	var leafe uint64

	leafe = uint64(len(mTree) / 2)
	for i = 0; i < 4; i++ {
		w, _ := gw.Get(common.HexToAddress(addr[i]))
		hw := sha3.NewKeccak256()
		hw.Write(w.Addr[:])

		for j=0; j<100; j++ {
			w.SetBit(i*100+j, 1)
			hw.Sum(mTree[leafe+i*100+j][:0])
		}
		gw.Put(w)
		fmt.Printf("%x, %v, %v, %s\n", w.Addr, w.Sid, w.id, w.Data)
	}

	i = uint64(len(mTree)/2 - 1)
	for ; i >= 1; i-- {
		left := i << 1
		right := i<<1 + 1
		hw := sha3.NewKeccak256()
		hw.Write(mTree[left][:])
		hw.Write(mTree[right][:])
		hw.Sum(mTree[i][:0])
		hw.Reset()
	}
	fmt.Printf("%x %d\n", mTree[1], len(mTree[leafe:]))
	gw.ROOT = mTree[1]

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, mTree[leafe:])
	gw.db.Put([]byte("GHASH"), buf.Bytes())


	w.SetBit(0, 1)
	w.SetBit(1, 1)
	w.SetBit(7, 1)
	w.SetBit(FG_BIT_SIZE, 1)
	pos := w.GetGPos(16, 0)
	fmt.Println(pos)
	pos = w.GetPos(16)
	fmt.Println(pos)
*/
}

