// handle the client's state, configuration, and interactions with the blockchain and server

package client

import (
	"log"
	"math/big"
	"poc-client/config"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (c *Client) SetChanId(id common.Hash) {
	c.ChannelID = id
}

func (c *Client) GetNonce(n uint64) uint64 {
	if n > c.LocalNonce {
		c.LocalNonce = n
	}
	return c.LocalNonce
}
func NewClient() (*Client, error) {
	config, err := config.LoadConfig("localConfig.json")
	if err != nil {
		return nil, err
	}
	client := &Client{}
	return client.Init(config.PrivateKeyFilePath, config.BcWsEndpoint, config.BcRpcEndpoint, config.ServerEndpoint, config.ContractAddress), nil
}

func (c *Client) Start(hubSend chan<- []byte) {
	var wg sync.WaitGroup
	wg.Add(1)

	log.Println("\n------------------Handshake-------------------")
	// Step 1: Handshake
	c.SendHandshakeMsg(&wg, hubSend)

	wg.Wait()

	wg.Add(1)
	// Step 2: Open a channel
	c.sendOpenChanTx(&wg, hubSend)
	log.Println()

	// Step 3: Wait for ChannelID to be set
	for c.ChannelID == (common.Hash{}) { // Wait until ChannelID is set
		log.Println("Waiting for ChannelID.....")
		time.Sleep(2000 * time.Millisecond)
	}

	// Now the preparation has been done.
	// The light client can start sending requests

	wg.Add(1)
	// Step 4: Send tokens to other address
	c.sendTokenTx(&wg, hubSend, "0xf3D1dBbC7Db2CC7eAE8b44e7c4422DC041993178", big.NewInt(10000000))
	wg.Wait()

	// wg.Add(1)
	// // Step 5: Send a request to check balance
	// c.SendBalanceCheckRequest(&wg, hubSend)
	// wg.Wait()
}

func (c *Client) ConnectToBlockchain() (*ethclient.Client, error) {
	return ethclient.Dial(c.BcWsEndpoint)
}
