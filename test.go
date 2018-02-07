package main

import (
	"fmt"
	"go-unitcoin/libraries/chain/account/groupwallet"
)

func main() {
	var v byte
	v = 1

	if v == 1 || v == 0 {
		fmt.Println("V is 0 1")
	}
	w := groupwallet.NewGroupWallet("/home/austin/data", nil)
	/*
		var address common.Address
		for address, _ = range w.Wallets {
			break
		}

		go w.Worker()
		w.Verify <- groupwallet.Req_verify{address, 8}
	*/
	for v, s := range w.Wallets {
		fmt.Printf("%x:%v\n", v, s)
	}
}
