package connection

import (
	"log"
	"net/http"
	"poc-client/hub"
	"poc-client/hub/wsClient"

	"github.com/gorilla/websocket"
)

func ConnectToFE(hub *hub.Hub) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{}

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println("Setup error upgrading connection: ", err)
		}

		log.Println("Connection with front end established :)")

		fe := &wsClient.Client{Conn: conn, Send: make(chan []byte, 256), Receive: make(chan []byte, 256)}
		hub.Set_fe <- fe
		go fe.WritePump()
		go fe.ReadPump()

		go func() {
			hub.Send_fe <- []byte("FE: Connected to server")
		}()

		// send handshake information
		go func() {
			hub.Send_fn <- []byte("HANDSHAKE")
		}()
		// defer fe_conn.Close()
	}
}
