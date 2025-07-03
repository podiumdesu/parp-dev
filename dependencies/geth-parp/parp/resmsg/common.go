package resmsg

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/crypto"
)

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

type DataMsg struct {
	Type string      `json:"type"`
	Data ResponseMsg `json:"data"`
}

func (d *DataMsg) Bytes() []byte {
	return marshalToJson(d)
}

func marshalToJson(v interface{}) []byte {
	jsonMsg, _ := json.Marshal(v)
	return jsonMsg
}

func hashData(v interface{}) []byte {
	jsonBody, _ := json.Marshal(v)
	hash := crypto.Keccak256Hash(jsonBody)
	return hash.Bytes()
}
