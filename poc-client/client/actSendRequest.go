package client

import (
	"context"
	"encoding/json"
	"log"
	"poc-client/msg/request"
	"sync"
)

func (c *Client) SendBalanceCheckRequest(wg *sync.WaitGroup, hubSend chan<- []byte) {
	defer wg.Done()
	log.Println("\n------------------Send Balance Check Request--------------------")

	request, err := c.createBalanceCheckRequest(wg, hubSend, 100)
	if err != nil {
		log.Fatal("Failed to create a balance check requeest: ", err)
	}

	select {
	case hubSend <- request:
		log.Println("Request message sent successfully")
	default:
		log.Println("Failed to send request message: channel is full or closed")
	}
	log.Println("------------------------------------------------------")
}

// TODO: Can be refactored
// actOpenChannel has a similar code base

func (c *Client) createBalanceCheckRequest(wg *sync.WaitGroup, hubSend chan<- []byte, amount uint) ([]byte, error) {

	bcClient, err := c.ConnectToBlockchain()
	if err != nil {
		log.Fatal(err)
	}

	blockHeader, _ := bcClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	request := request.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBalance",
		Params:  []interface{}{"0xA2131E7503F7Dd11ff5dAAC09fa7c301e7Fe0f30", "latest"},
		ID:      1,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
	}

	msg := c.generateParpRequest(c.ChannelID, c.Amount+amount, jsonRequest, blockHeader.Hash())
	c.Amount = amount

	msgWType := append([]byte("REQ:"), msg...)
	log.Println("Sending Request: ", string(msgWType))

	return msgWType, nil
}
