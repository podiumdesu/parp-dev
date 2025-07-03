package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"poc-client/protocol"
	"poc-client/utils/cryptoutil"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) sendOpenChanTx(wg *sync.WaitGroup, hubSend chan<- []byte) (common.Hash, error) {
	defer wg.Done()
	log.Println("\n------------------Send OpenChan Tx request--------------------")
	OpenChanTx, txHash, err := c.createOpenChanTx()
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

	return txHash, nil
}

func (c *Client) createOpenChanTx() ([]byte, common.Hash, error) {
	// in @openChan.go as well

	bcClient, err := c.ConnectToBlockchain()
	if err != nil {
		return nil, common.Hash{}, err
	}

	blockHeader, err := bcClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, common.Hash{}, err
	}

	onChainNonce, err := bcClient.PendingNonceAt(context.Background(), c.Address)
	if err != nil {
		log.Fatal(err)
	}
	nonce := c.GetNonce(onChainNonce)
	fnAddr := common.HexToAddress(cryptoutil.PubkeyToHexAddr(c.ServerPublicKey))
	log.Println("fnAddr: ", fnAddr)

	deposit := big.NewInt(200000)

	log.Println("contract address: ", c.ContractAddress)
	openChanSignTx, txHash := protocol.OpenChanTx(bcClient, c.PrivateKey, fnAddr, deposit, nonce, common.HexToAddress(c.ContractAddress))

	msg := c.generateParpRequest(c.ChannelID, c.Amount+100, openChanSignTx, blockHeader.Hash())
	c.Amount += 100 // record the payment amount
	msgWType := append([]byte("OpenChan:"), msg...)

	log.Println("Sending: ", string(msgWType))

	return msgWType, txHash, nil
}
