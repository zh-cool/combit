package gwallet

import (
	"crypto/ecdsa"
	"go-unitcoin/libraries/chain/transaction/types"
	"go-unitcoin/libraries/crypto"
	"go-unitcoin/libraries/event"
	"go-unitcoin/libraries/chain/transaction"
)

type Verification struct {
	gw    *GWallet
	rsVryCH  event.Feed
	prv   *ecdsa.PrivateKey
	transactionPool *transaction.TransactionPool
	txMsg chan TransactionMsgPreEvent
}
type VerPreEvent struct{ ve *Res_verify }

type TransactionMsgPreEvent struct{ Msg *types.Transaction_Message }

func NewVerification(gw *GWallet, prv *ecdsa.PrivateKey) *Verification {
	if prv == nil {
		prv, _ = crypto.GenerateKey()
	}

	return &Verification{
		gw:    gw,
		prv:   prv,
	}
}

func (vr *Verification) Start(transactionPool *transaction.TransactionPool) {
	vr.txMsg = make(chan TransactionMsgPreEvent)
	transactionPool.SubscribeVerifyEvent(vr.txMsg)
	go vr.Worker()
}

func (vr *Verification) SubscribeRsVerifyEvent(rsCh chan VerPreEvent) {
	vr.rsVryCH.Subscribe(rsCh)
}

func (vr *Verification) Worker() {
	for {
		select {
		case v := <-vr.txMsg:

			to, err := vr.gw.IDtoAddr(v.Msg.To())
			if err != nil {

			}

			from, err := vr.gw.IDtoAddr(v.Msg.From())
			if err != nil {

			}

			pos := v.Msg.CoinIndex()
			vr.gw.Move(to, from, uint64(len(pos)))

			rs := &Res_verify{&vr.gw.ROOT, vr.gw.GID, pos, v.Msg.TxHash(), []byte{}}
			rs.Sign(vr.prv)
			vr.rsVryCH.Send(VerPreEvent{rs})
		}
	}
}
