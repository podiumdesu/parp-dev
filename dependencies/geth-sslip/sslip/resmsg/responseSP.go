package resmsg

import (
	"encoding/hex"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// response of storage proof

// for storage proof
type ResponseSPBody struct {
	SignedReqBody []byte
	Proof         [][]byte
	// Address       common.Address
	// BlockNr *big.Int
}
type ResponseSPMsg struct {
	Type               string
	ChannelId          common.Hash
	Amount             uint
	ReqBodyHash        common.Hash
	SignedReqBody      []byte
	CurrentBlockHeight *big.Int
	ReturnValue        []byte
	Proof              [][]byte
	Address            common.Address
	// BlockNr            *big.Int
	Signature []byte
}

func (r *ResponseSPMsg) Bytes() []byte {
	return marshalToJson(r)
}

func (rb *ResponseSPBody) Keccak256Hash() common.Hash {
	data := []byte{}
	data = append(data, rb.SignedReqBody...)

	for _, proofItem := range rb.Proof {
		data = append(data, []byte(proofItem)...) // Proof as bytes array
	}
	// data = append(data, rb.Address.Bytes()...)
	// data = append(data, rb.BlockNr.Bytes()...)

	// log.Println(hex.EncodeToString(rb.BlockNr.Bytes()))
	hash := crypto.Keccak256Hash(data)

	// log.Println("----------------------NOW I WANT TO CHECK THE HASH!!!!----------------------")
	// log.Println(rb.Proof)
	// log.Println(rb.TxHash)
	// log.Println(rb.TxIdx)
	// log.Println(data)

	// log.Println("----------------------END!----------------------")
	return hash
}

func (r *ResponseSPMsg) RlpBytes() string {
	// Encode the ResponseMsg to RLP bytes
	serializedBytes, err := rlp.EncodeToBytes(r)
	if err != nil {
		log.Fatalf("Failed to RLP encode ResponseMsg: %v", err)
	}

	// Convert the serialized RLP bytes to a hex string, prefixed with "0x"
	hexString := "0x" + hex.EncodeToString(serializedBytes)

	return hexString
}

func (r *ResponseSPMsg) BodyHashBytes() []byte {
	data := ResponseSPBody{
		SignedReqBody: r.SignedReqBody,
		Proof:         r.Proof,
		// Address:       r.Address,
		// BlockNr:       r.BlockNr,
	}
	return hashData(data)
}

func (r *ResponseSPBody) HashBytes() []byte {
	return hashData(r)
}
