package client

import (
	"poc-client/msg/request"
)

func (c *Client) GenerateRequestMsg(cID int, t int) *request.Msg {

	return &request.Msg{
		ChannelID: 1,
		Type:      0,
		RequestBody: request.RequestBody{
			ChannelID:   1,
			RequestByte: []byte("curl https://api.coindesk.com/v1/bpi/currentprice.json"),
		},
	}
}
