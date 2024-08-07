// handle the client's state, configuration, and interactions with the blockchain and server

package client

import (
	"log"
	"poc-client/config"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

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
	go c.SendHandshakeMsg(&wg, hubSend)

	wg.Wait()

	wg.Add(1)
	go c.sendOpenChanTx(&wg, hubSend)

}

func (c *Client) ConnectToBlockchain() (*ethclient.Client, error) {
	return ethclient.Dial(c.BcWsEndpoint)
}
