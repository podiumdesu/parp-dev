package client

import (
	"crypto/ecdsa"
	"math/big"
	"poc-client/msg/handshake"
	"poc-client/msg/request"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Client struct {
	BcWsEndpoint    string
	BcRpcEndpoint   string
	ServerEndpoint  string
	PrivateKey      *ecdsa.PrivateKey
	PublicKey       *ecdsa.PublicKey
	ServerPublicKey *ecdsa.PublicKey
	ContractAddress string
	Address         common.Address
	ConnectFN       common.Address
	FeeStandards    int
	Secret          *big.Int
	ChannelID       *big.Int
	StartTime       time.Time
	Duration        time.Duration
	Amount          uint
	LastSignedProof []byte
	BlockHeader     *types.Header
	Checkpoint      *types.Header
	Step            int // 0-IDLE, 1-HANDSHAKING, 3-UNSTAKED, 4-COMMITTED, 5-UNBONDING
}

type IClient interface {
	Init(p string) *Client
	PrivKeyBytes() []byte
	PubKeyBytes() []byte
	AddrHex() string
	Sign(hash []byte) []byte
	Verify(hash []byte, sig []byte) bool

	InitHandshakeMsg(duration time.Duration, deposit *big.Int, secret *big.Int) *handshake.Msg
	InitRequestMsg(cID int, t int) *request.Msg
}
