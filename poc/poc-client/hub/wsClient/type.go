package wsClient

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	ChannelId []byte
	Send      chan []byte
	Receive   chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.Receive <- message
	}

}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		// log.Println("wsClient Send: ", string(message))

		w, err := c.Conn.NextWriter(websocket.TextMessage)

		if err != nil {
			log.Println("wsClient: NextWriter error: ", err)
		}

		if _, err := w.Write(message); err != nil {
			log.Println("wsClient: Write error: ", err)
		}

		if err := w.Close(); err != nil {
			log.Println("Error closing writer: ", err)
		}
		// }
	}

	log.Println("WritePump stopped")

}
