// manage the hub functionalities

package hub

import (
	"poc-client/hub/wsClient"
)

type WsClient = wsClient.Client
type Hub struct {
	// registered client
	Set_fe     chan *WsClient
	Set_fn     chan *WsClient
	Send_fe    chan []byte
	Send_fn    chan []byte
	Receive_fe chan []byte
	Receive_fn chan []byte
	fe         *WsClient
	fn         *WsClient
}

func NewHub() *Hub {
	return &Hub{
		Set_fe:     make(chan *WsClient),
		Set_fn:     make(chan *WsClient),
		Send_fe:    make(chan []byte, 10000),
		Send_fn:    make(chan []byte, 10000),
		Receive_fe: make(chan []byte, 10000),
		Receive_fn: make(chan []byte, 10000),
		fe:         &WsClient{},
		fn:         &WsClient{},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Set_fe:
			h.fe = client
		case client := <-h.Set_fn:
			h.fn = client
		case msg := <-h.Send_fe:
			h.fe.Send <- msg
		case msg := <-h.Send_fn:
			h.fn.Send <- msg
		}
	}
}
