package client

import (
	"fmt"
	"log"
	"math/big"
	"poc-client/msg/handshake"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

func (c *Client) generateHandshakeMsg(contractAddress string, duration time.Duration, deposit *big.Int, secret *big.Int) *handshake.Msg {
	b := &handshake.MsgBody{
		PubKB:           crypto.FromECDSAPub(c.PublicKey),
		ContractAddress: contractAddress,
		Duration:        duration * time.Hour,
		Deposit:         deposit,
		SecretN:         secret,
	}
	h := b.DigestHash()
	sig := c.Sign(h)

	return &handshake.Msg{
		Body: b,
		Sig:  sig,
	}
}

func (c *Client) SendHandshakeMsg(wg *sync.WaitGroup, hubSend chan<- []byte) {
	defer wg.Done()
	msg := c.generateHandshakeMsg(c.ContractAddress, 10, big.NewInt(100000), big.NewInt(200))
	b := append([]byte("HANDSHAKE:"), msg.Bytes()...)

	log.Println("Sending: ", b)
	hubSend <- b
	fmt.Println("----------------------------------------------")
	hubSend <- []byte("FE: Connected to server")
}
