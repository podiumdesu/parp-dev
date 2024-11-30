package client

import (
	"poc-client/msg/request"

	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) GenerateRequestMsg(cID common.Hash, t int) *request.Msg {

	return &request.Msg{
		ChannelID: cID,
		Type:      0,
		RequestBody: request.RequestBody{
			ChannelID:   cID,
			RequestByte: []byte("curl https://api.coindesk.com/v1/bpi/currentprice.json"),
		},
	}
}
