package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"

	"poc-client/protocol"

	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) sendTokenTx(wg *sync.WaitGroup, hubSend chan<- []byte, recipient string, tokenAmount *big.Int) common.Hash {
	defer wg.Done()
	log.Println("\n------------------Send Token Tx request--------------------")
	TokenTx, txHash, err := c.createTokenTx(recipient, tokenAmount)
	log.Println("token transaction hash: ", txHash.Hex())
	if err != nil {
		log.Println("Failed to create Token Tx request: ", err)
	}

	select {
	case hubSend <- TokenTx:
		log.Println("Token transaction message sent successfully")
	default:
		log.Println("Failed to send token transaction message: channel is full or closed")
	}
	fmt.Println("------------------------------------------------------")

	return txHash
}

func (c *Client) createTokenTx(recipient string, tokenAmount *big.Int) ([]byte, common.Hash, error) {
	// Connect to the blockchain client
	bcClient, err := c.ConnectToBlockchain()
	if err != nil {
		return nil, common.Hash{}, err
	}

	// Get the current block header
	blockHeader, err := bcClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, common.Hash{}, err
	}

	// Convert recipient address
	recipientAddr := common.HexToAddress(recipient)

	// Create the transaction
	tokenTx, txHash := protocol.TransferTokenTx(bcClient, c.PrivateKey, recipientAddr, tokenAmount, common.HexToAddress(c.ContractAddress))

	// Generate a message payload
	msg := c.generateParpRequest(c.ChannelID, c.Amount+100, tokenTx, blockHeader.Hash())
	c.Amount += 100

	// Add a message type prefix
	msgWType := append([]byte("TX:"), msg...)

	log.Println("Sending Token Transaction: ", string(msgWType))

	return msgWType, txHash, nil
}
