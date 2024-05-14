package handshake

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type Msg struct {
	Body *MsgBody
	Sig  []byte
}

type IMsg interface {
	Init()
}

func (m *Msg) Bytes() []byte {
	jsonMsg, _ := json.Marshal(m)
	return jsonMsg
}

type MsgBody struct {
	PubKB           []byte
	ContractAddress string
	Duration        time.Duration
	Deposit         *big.Int
	SecretN         *big.Int
}

func (b *MsgBody) DigestHash() (sig []byte) {
	jsonBody, _ := json.Marshal(b)
	hash := crypto.Keccak256Hash(jsonBody)
	return hash.Bytes()
}
