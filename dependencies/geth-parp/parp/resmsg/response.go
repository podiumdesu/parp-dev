package resmsg

import (
	"encoding/hex"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type ResponseBody struct {
	SignedReqBody []byte
	Proof         [][]byte
	TxHash        common.Hash
	TxIdx         []byte
}

type ResponseMsg struct {
	Type               string
	ChannelId          common.Hash
	Amount             uint
	ReqBodyHash        common.Hash
	SignedReqBody      []byte
	CurrentBlockHeight *big.Int
	ReturnValue        []byte
	Proof              [][]byte
	TxHash             common.Hash
	TxIdx              []byte
	Signature          []byte
	RootHash           common.Hash
}

// func (rb *ResponseBody) HashBytes() []byte {
// 	return hashData(rb)
// }

func (r *ResponseMsg) Bytes() []byte {
	return marshalToJson(r)
}

func (rb *ResponseMsg) Keccak256Hash() common.Hash {
	data := []byte{}
	data = append(data, rb.ChannelId.Bytes()...)
	data = append(data, rb.SignedReqBody...)

	for _, proofItem := range rb.Proof {
		data = append(data, []byte(proofItem)...) // Proof as bytes array
	}
	data = append(data, rb.TxHash.Bytes()...)
	data = append(data, rb.TxIdx...)

	hash := crypto.Keccak256Hash(data)

	// log.Println("----------------------NOW I WANT TO CHECK THE HASH!!!!----------------------")
	// log.Println(rb.Proof)
	// log.Println(rb.TxHash)
	// log.Println(rb.TxIdx)
	// log.Println(data)

	// log.Println("----------------------END!----------------------")
	return hash
}

func (r *ResponseMsg) RlpBytes() string {
	// Encode the ResponseMsg to RLP bytes
	serializedBytes, err := rlp.EncodeToBytes(r)
	if err != nil {
		log.Fatalf("Failed to RLP encode ResponseMsg: %v", err)
	}

	// Convert the serialized RLP bytes to a hex string, prefixed with "0x"
	hexString := "0x" + hex.EncodeToString(serializedBytes)

	return hexString
}

// func (rb *ResponseBody) Keccak256Hash() common.Hash {
// 	data := []byte{}
// 	data = append(data, rb.SignedReqBody...)

// 	for _, proofItem := range rb.Proof {
// 		data = append(data, []byte(proofItem)...) // Proof as bytes array
// 	}
// 	data = append(data, rb.TxHash.Bytes()...)
// 	data = append(data, rb.TxIdx...)

// 	hash := crypto.Keccak256Hash(data)

// 	// log.Println("----------------------NOW I WANT TO CHECK THE HASH!!!!----------------------")
// 	// log.Println(rb.Proof)
// 	// log.Println(rb.TxHash)
// 	// log.Println(rb.TxIdx)
// 	// log.Println(data)

// 	// log.Println("----------------------END!----------------------")
// 	return hash
// }

// func (r *ResponseMsg) BodyHashBytes() []byte {
// 	data := ResponseBody{
// 		// SignedReqBody: r.SignedReqBody,
// 		// Proof:  r.Proof,
// 		TxHash: r.TxHash,
// 		TxIdx:  r.TxIdx,
// 	}
// 	return hashData(data)
// }
