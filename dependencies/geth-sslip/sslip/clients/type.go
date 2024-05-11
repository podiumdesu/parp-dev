package mClient

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	ID      string
	Conn    *websocket.Conn
	PubKeyB []byte
}

func New(id string, conn *websocket.Conn) *Client {
	return &Client{
		ID:      id,
		Conn:    conn,
		PubKeyB: []byte{},
	}
}
