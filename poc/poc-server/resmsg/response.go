package resmsg

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// type ServerMsg interface {
// 	MessageType() string
// }

type ServerMsg struct {
	Type string `json:"type"`
	Info []byte `json:"info"`
}

func (s *ServerMsg) MessageType() string {
	return s.Type
}

func (s *ServerMsg) Bytes() []byte {
	jsonMsg, _ := json.Marshal(s)
	return jsonMsg
}

type HandshakeMsg struct {
	Type             string `json:"type"`
	ServerPublicKeyB []byte `json:"serverPublicKeyB"`
	Signature        []byte `json:"signature"`
}

type HandshakeMsgBody struct {
	Type             string `json:"type"`
	ServerPublicKeyB []byte `json:"serverPublicKeyB"`
}

func (h *HandshakeMsg) HashBytes() []byte {
	data := HandshakeMsgBody{
		Type:             h.Type,
		ServerPublicKeyB: h.ServerPublicKeyB,
	}
	jsonBody, _ := json.Marshal(data)
	hash := crypto.Keccak256Hash(jsonBody)
	return hash.Bytes()
}
func (h *HandshakeMsg) Bytes() []byte {
	jsonMsg, _ := json.Marshal(h)
	return jsonMsg
}

type DataMsg struct {
	Type string      `json:"type"`
	Data ResponseMsg `json:"data"`
}

func (d *DataMsg) Bytes() []byte {
	jsonMsg, _ := json.Marshal(d)
	return jsonMsg
}

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
	SignedReqBody      []byte
	CurrentBlockHeight *big.Int
	ReturnValue        []byte
	Proof              [][]byte
	TxHash             common.Hash
	TxIdx              []byte
	Signature          []byte
}

func (rb *ResponseBody) HashBytes() []byte {
	jsonBody, _ := json.Marshal(rb)
	hash := crypto.Keccak256Hash(jsonBody)
	return hash.Bytes()
}

func (rb *ResponseBody) Keccak256Hash() common.Hash {
	data := []byte{}

	data = append(data, rb.SignedReqBody...)
	data = append(data, rb.TxHash.Bytes()...)
	data = append(data, rb.TxIdx...)

	for _, proofItem := range rb.Proof {
		data = append(data, []byte(proofItem)...) // Proof as bytes array
	}
	hash := crypto.Keccak256Hash(data)

	return hash
}

func (r *ResponseMsg) Bytes() []byte {
	jsonMsg, _ := json.Marshal(r)
	return jsonMsg
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

func (r *ResponseMsg) BodyHashBytes() []byte {
	data := ResponseBody{
		SignedReqBody: r.SignedReqBody,
		Proof:         r.Proof,
		TxHash:        r.TxHash,
		TxIdx:         r.TxIdx,
	}
	jsonBody, _ := json.Marshal(data)
	hash := crypto.Keccak256Hash(jsonBody)
	return hash.Bytes()
}

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
	jsonMsg, _ := json.Marshal(r)
	return jsonMsg
}

func (r *ResponseSPMsg) BodyHashBytes() []byte {
	data := ResponseSPBody{
		SignedReqBody: r.SignedReqBody,
		Proof:         r.Proof,
		Address:       r.Address,
		BlockNr:       r.BlockNr,
	}
	jsonBody, _ := json.Marshal(data)
	hash := crypto.Keccak256Hash(jsonBody)
	return hash.Bytes()
}

func (r *ResponseSPBody) HashBytes() []byte {
	jsonBody, _ := json.Marshal(r)
	hash := crypto.Keccak256Hash(jsonBody)
	return hash.Bytes()
}
