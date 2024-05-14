package connection

import (
	"log"

	"poc-client/client"

	"github.com/gorilla/websocket"
)

// Connect PoC-Client-Backend to PoC-Server
func ConnectToServer(client *client.Client) *websocket.Conn {
	dialer := websocket.DefaultDialer
	// const serverInfo = "8081:20"
	// s := strings.Split(serverInfo, ":")
	// port, clientID := s[0], s[1]
	// url := "ws://localhost:" + port + "/ws/" + clientID
	conn, _, err := dialer.Dial(client.ServerEndpoint, nil)

	if err != nil {
		log.Println("ConnectToServer Error upgrading connection: ", err)
		return nil
	}

	log.Println("Connected to server: ", client.ServerEndpoint)
	if err != nil {
		log.Println("Error upgrading connection: ", err)
	}

	return conn
}
