package request

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// the request information sent to PoC Server

// Which will be used as the proof of this request
type RequestBody struct {
	ChannelID   int
	RequestByte []byte // either request or tx, for example curl xxxx.com
}

type PaymentBody struct {
	ChannelID int
	Amount    uint // the total money the client is willing to pay
}
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type Msg struct {
	ChannelID     int
	Type          int // 0 for regular request, 1 for on-chain tx
	RequestBody   RequestBody
	PaymentBody   PaymentBody
	SignedRequest []byte
	SignedPayment []byte
}

type ReqBody struct {
	ChannelID      int
	Amount         uint
	ReqByte        []byte
	LocalBlockHash common.Hash
}

type RequestMsg struct {
	ChannelID         int
	Type              int
	Amount            uint
	ReqByte           []byte
	LocalBlockHash    common.Hash
	SignedReqBody     []byte
	SignedPaymentBody []byte
}

func (rb *ReqBody) HashByte() []byte {
	jsonReq, _ := json.Marshal(rb)
	hash := crypto.Keccak256Hash(jsonReq)
	return hash.Bytes()
}

func (rb *ReqBody) PreHashByte() []byte {
	return GeneratePrefixedHash(rb.HashByte())
}

func (pb *PaymentBody) HashByte() []byte {

	bytes := encodeData(pb.ChannelID, pb.Amount)
	return bytes
	// jsonPay, _ := json.Marshal(pb)
	// hash := crypto.Keccak256Hash(jsonPay)
	// return hash.Bytes()
}

func (pb *PaymentBody) PreHashByte() []byte {
	log.Println("----------PB----")

	// prefixed hash for ethereum signautre message
	return GeneratePrefixedHash(pb.HashByte())
}
func encodeData(channelID int, amount uint) []byte {
	buf := new(bytes.Buffer)

	// Create a bytes slice fully padded for uint256
	paddedChannelID := make([]byte, 32)
	paddedAmount := make([]byte, 32)

	// Convert integers to big-endian, 64-bit integers
	binary.BigEndian.PutUint64(paddedChannelID[24:], uint64(channelID)) // last 8 bytes for uint64
	binary.BigEndian.PutUint64(paddedAmount[24:], uint64(amount))       // last 8 bytes for uint64

	// Write padded data to buffer
	buf.Write(paddedChannelID)
	buf.Write(paddedAmount)

	return buf.Bytes()
}
func GeneratePrefixedHash(hash []byte) []byte {
	// prefixed hash for ethereum signautre message
	prefixedHash := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(hash)) + string(hash))).Bytes()
	return prefixedHash
}
