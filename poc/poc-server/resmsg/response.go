package resmsg

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	TxIdx         uint32
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
	TxIdx              uint32
	Signature          []byte
}

func (rb *ResponseBody) HashBytes() []byte {
	jsonBody, _ := json.Marshal(rb)
	hash := crypto.Keccak256Hash(jsonBody)
	return hash.Bytes()
}

func (r *ResponseMsg) Bytes() []byte {
	jsonMsg, _ := json.Marshal(r)
	return jsonMsg
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
