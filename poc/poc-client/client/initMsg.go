package client

import (
	"math/big"
	"poc-client/msg/handshake"
	"poc-client/msg/request"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

func (c *Client) InitHandshakeMsg(contractAddress string, duration time.Duration, deposit *big.Int, secret *big.Int) *handshake.Msg {
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

func (c *Client) InitRequestMsg(cID int, t int) *request.Msg {

	return &request.Msg{
		ChannelID: 1,
		Type:      0,
		RequestBody: request.RequestBody{
			ChannelID:   1,
			RequestByte: []byte("curl https://api.coindesk.com/v1/bpi/currentprice.json"),
		},
	}
}
