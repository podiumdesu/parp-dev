package connection

import (
	"poc-client/client"

	"math/rand"
	"strconv"

	"github.com/gorilla/websocket"
)

// Connect PoC-Client-Backend to PoC-Server
func ConnectToServer(client *client.Client) (*websocket.Conn, error) {
	dialer := websocket.DefaultDialer
	// const serverInfo = "8081:20"
	// s := strings.Split(serverInfo, ":")
	// port, clientID := s[0], s[1]
	// url := "ws://localhost:" + port + "/ws/" + clientID
	id := rand.Intn(10000000000)

	// Connect to the end server with the generated clientID
	client.ServerEndpoint = client.ServerEndpoint + strconv.Itoa(id)
	conn, _, err := dialer.Dial(client.ServerEndpoint, nil)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
