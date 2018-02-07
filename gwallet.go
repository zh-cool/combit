package groupwallet

import (
	"go-unitcoin/libraries/common"
	"go-unitcoin/libraries/chain/util"
	"go-unitcoin/libraries/chain/space/protocol"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"fmt"
	//"math/big"
)

type conbin struct {
	pos uint
	len uint
}

type Req_verify struct {
	From	common.Address
	Number  uint
}

type Wallet struct {
	Wallets map[common.Address][]conbin
	//bins    map[common.Address] map[int]*big.Int
	binpath string

	CID     uint
	ROOT    common.Hash
	GHASH   []common.Hash
	GADDR 	[]common.Address
	ID      uint

	txpool     *util.TxPool
	txCH 		chan protocol.TxPreEvent
	ids		[]int
	Verify  chan Req_verify
}

func getList(pos []conbin, begin, end , number uint) []uint {
	var bpos, epos uint
	for i, v := range pos {
		if v.pos + v.len > begin && bpos==0 {
			bpos = uint(i)
		}

		if v.pos < end {
			epos = uint(i)
		}
	}

	fmt.Println("bpos:", bpos, "epos", epos)
	var list []uint
	for i:=epos; i >= bpos; i-- {
		rg := pos[i]
		for k:= rg.len; k>0; k-- {
			if rg.pos + k < end && rg.pos + k >= begin {
				list = append(list, rg.pos+k)
			}
		}
	}
	return list[0:number]
}



type myReader struct {
	r     io.Reader
	bits  uint
	count uint
	prev  byte
}

func newMyReader(file string) *myReader {

	r, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	return &myReader{r, 0, 0, 0}
}

var conv = map[byte]uint{
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
	'A': 10,
	'b': 11,
	'B': 11,
	'c': 12,
	'C': 12,
	'd': 13,
	'D': 13,
	'e': 14,
	'E': 14,
	'f': 15,
	'F': 15,
	'g': 0,
	'h': 1,
}

func (rd *myReader) Read(p []byte) (n int, err error) 	{
	data := make([]byte, 1<<20, 1<<20)
	var v byte

	for {
		n, err = rd.r.Read(data)
		if err == io.EOF {
			break
		}
		for i:=0; i<n; i++ {
			v = data[i]
			if v=='\n' || v==' ' {
				continue
			}
			//v := conv[b]
			if v == '0' || v == '1' {
				if rd.prev == '0' {
					if rd.count==0 {
						rd.count = 1
					}
					for rd.count >= 0 {
						p[rd.bits/8] &^= (1 << (7-(rd.bits % 8)))
						rd.bits++
						rd.count--
						if rd.count == 0 {
							break
						}
					}
				} else if rd.prev == '1' {
					if rd.count==0 {
						rd.count = 1
					}
					for rd.count >= 0 {
						p[rd.bits/8] |= (1 << (7 - (rd.bits % 8)))
						rd.bits++
						rd.count--
						if rd.count == 0 {
							break
						}
					}
				}
				rd.prev = v
				rd.count = 0
			} else {
				rd.count = rd.count*16 + conv[v]
			}
		}

		if rd.prev == '0' {
			if rd.count==0 {
				rd.count = 1
			}
			for rd.count >= 0 {
				p[rd.bits/8] &^= (1 << (7-(rd.bits % 8)))
				rd.bits++
				rd.count--
				if rd.count == 0 {
					break
				}
			}
		} else if rd.prev == '1' {
			if rd.count==0 {
				rd.count = 1
			}
			for rd.count >= 0 {
				p[rd.bits/8] |= (1 << (7 - (rd.bits % 8)))
				rd.bits++
				rd.count--
				if rd.count == 0 {
					break
				}
			}
		}
		rd.prev = v
		rd.count = 0

	}

	return n, nil
}

type cbit struct {
	b     uint8
	count uint
}

func (rd *myReader) conbinpos() []conbin {
	var TBIT uint8 = 128
	var bin []conbin
	var cmbit []cbit
	var i, k uint
	data := make([]byte, 1<<29)

	rd.Read(data)
	count := (rd.bits+7)/8

	cb := cbit{128, 0}
	for k = 0; k<count; k++ {
		bt := data[k]
		for i = 0; (i < 8) && (rd.bits > 0); i++ {
			bi := uint8(bt) & (TBIT >> i)
			bi = bi >> (7 - i)

			if bi == cb.b {
				cb.count++
			} else {
				cmbit = append(cmbit, cb)
				cb.b = bi
				cb.count = 1
			}
			rd.bits--
		}
	}
	cmbit = append(cmbit, cb)
	fmt.Println(cmbit)

	var cn uint
	for _, cb = range cmbit[1:] {
		if cb.b == 1 {
			bin = append(bin, conbin{cn, cb.count})
		}
		cn += cb.count
	}

	return bin
}

func (rd *myReader) Address() common.Address {
	var addr common.Address
	buf := make([]byte, 40)
	rd.r.Read(buf[0:40])
	addr = common.HexToAddress(string(buf))
	return addr
}

func (rd *myReader) Close() {
	file := rd.r.(*os.File)
	file.Close()
}

func NewGroupWallet(conbindir string, txpool *util.TxPool) *Wallet {

	w := Wallet{}
	w.txpool = txpool
	w.Wallets = make(map[common.Address][]conbin)
	w.Verify = make(chan Req_verify)
	if filepath.IsAbs(conbindir) {
		w.binpath = conbindir
	}else{
		w.binpath, _ = filepath.Abs(conbindir)
	}
	files := w.binfiels(conbindir)

	for _, v := range files {
		var address common.Address
		rd := newMyReader(conbindir + "/" + v.Name())
		address = rd.Address()
		w.Wallets[address] = rd.conbinpos()
		rd.Close()
	}

	return &w
}

func (w *Wallet) Start() {
	w.txCH = make(chan protocol.TxPreEvent)
	w.txpool.SubscribeTxPreEvent(w.txCH)
	fmt.Println(("Austin debug"))
	go w.CheckTxBin()
	go w.Worker()
}

func (w *Wallet) CheckTxBin() {
	fmt.Println(("Austin debug CheckTxBin "))
	for {
		select {
		case event := <-w.txCH:
			fmt.Println(event)
			from, num := func(tx *protocol.Tx) (common.Address, uint) {
					signer := protocol.MakeSigner()
					from, _:= protocol.Sender(signer, tx)
					value := tx.Value()
					value.Int64()
					return from, uint(value.Uint64())

			}(event.Transaction)
			w.Verify <- Req_verify{from, num}
		}
	}
	fmt.Println(("Austin debug out"))
}

func (w *Wallet) Worker() {
	w.ID = 1
	for {
		select {
		case v := <-w.Verify:
			begin, end := func(id uint) (begin, end uint){
				return 1024*(id-1), 1024*id
			}(w.ID)

			binpos := w.Wallets[v.From]
			list := getList(binpos, begin, end, v.Number)
			fmt.Println(list)
		}
	}
}

func (w *Wallet) Bound(id int) {
	//v := big.Int{}
	//v.SetBytes()
}

func (w *Wallet) binfiels(conbindir string) []os.FileInfo {
	files, err := ioutil.ReadDir(conbindir)
	if err != nil {
		log.Fatal(err)
	}

	var rfiles []os.FileInfo
	for _, v := range files {
		if v.Mode().IsRegular() {
			rfiles = append(rfiles, v)
		}
	}

	return rfiles
}

func (w *Wallet) Blance(address common.Address) (int64, bool) {
	var blance int64 = 0
	if wallet, ok := w.Wallets[address]; ok {
		for _, v := range wallet {
			blance += int64(v.len)
		}
	}else{
		return -1, false
	}

	return blance, true
}
/*
func (w *Wallet) Verfiy(address common.Address, conbin int64) bool {
	if blance, ok := w.Blance(address); ok {
		return blance >= conbin
	}
	fmt.Println("OK")
	return false
}
*/
/*
func (w *Wallet) Sync() bool {
	files := w.binfiels(w.binpath)

	for address, wallet := range w.Wallets {
		fmt.Println(address, ":", wallet)
		files := w.binfiels(w.binpath)
	}

	return true
}
*/
func (w *Wallet) Update() {

}
