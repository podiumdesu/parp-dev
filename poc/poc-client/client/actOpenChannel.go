package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"poc-client/protocol"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) sendOpenChanTx(wg *sync.WaitGroup, hubSend chan<- []byte) {
	log.Println("\n------------------Send OpenChan Tx request--------------------")
	OpenChanTx, err := c.createOpenChanTx()
	if err != nil {
		log.Println("Failed to create OpenChan request: ", err)
	}

	select {
	case hubSend <- OpenChanTx:
		log.Println("Request message sent successfully")
	default:
		log.Println("Failed to send request message: channel is full or closed")
	}
	fmt.Println("------------------------------------------------------")

}

func (c *Client) createOpenChanTx() ([]byte, error) {
	// in @openChan.go as well

	bcClient, err := c.ConnectToBlockchain()
	if err != nil {
		return nil, err
	}

	blockHeader, err := bcClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Block Hash: ", blockHeader.Hash().Hex())

	// log.Println("\n------------------Send OpenChan request--------------------")
	// TODO: Forgot why it is arbitrary here, need to revise it

	fnAddr := common.HexToAddress("0xA2131E7503F7Dd11ff5dAAC09fa7c301e7Fe0f30")
	deposit := big.NewInt(200000)
	openChanSignTx := protocol.OpenChanTx(bcClient, c.PrivateKey, fnAddr, deposit, common.HexToAddress(c.ContractAddress))

	msg := c.generateParpRequest(20, c.Amount+100, openChanSignTx, blockHeader.Hash())

	c.Amount += 100 // record the payment amount
	msgWType := append([]byte("TX:"), msg...)

	log.Println("Sending: ", string(msgWType))

	return msgWType, nil
}
