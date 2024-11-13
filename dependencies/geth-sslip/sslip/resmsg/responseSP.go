package resmsg

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// response of storage proof

// for storage proof
type ResponseSPBody struct {
	SignedReqBody []byte
	Proof         [][]byte
	Address       common.Address
	BlockNr       *big.Int
}
type ResponseSPMsg struct {
	Type               string
	ChannelId          common.Hash
	Amount             uint
	SignedReqBody      []byte
	CurrentBlockHeight *big.Int
	ReturnValue        []byte
	Proof              [][]byte
	Address            common.Address
	BlockNr            *big.Int
	Signature          []byte
}

func (r *ResponseSPMsg) Bytes() []byte {
	return marshalToJson(r)
}

func (r *ResponseSPMsg) BodyHashBytes() []byte {
	data := ResponseSPBody{
		SignedReqBody: r.SignedReqBody,
		Proof:         r.Proof,
		Address:       r.Address,
		BlockNr:       r.BlockNr,
	}
	return hashData(data)
}

func (r *ResponseSPBody) HashBytes() []byte {
	return hashData(r)
}
