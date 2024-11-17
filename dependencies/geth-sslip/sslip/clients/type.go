package mClient

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID        string
	Conn      *websocket.Conn
	PubKeyB   []byte
	ChannelID string
	WriteChan chan []byte
}

func New(id string, conn *websocket.Conn) *Client {
	return &Client{
		ID:        id,
		Conn:      conn,
		PubKeyB:   []byte{},
		ChannelID: "",
		WriteChan: make(chan []byte, 256), // Buffered channel to avoid blocking
	}
}

func (c *Client) ConnectToBlockchain() (*ethclient.Client, error) {
	return ethclient.Dial("ws://127.0.0.1:8101")
}

func (c *Client) WritePump() {
	for message := range c.WriteChan {
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error writing message to client %s: %v", c.ID, err)
			break
		}
	}
	c.Conn.Close()
}

func (client *Client) Send(message []byte) {
	select {
	case client.WriteChan <- message:
	default:
		log.Printf("Client %s write channel is full. Dropping message.", client.ID)
	}
}
