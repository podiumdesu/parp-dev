package web

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// homeTemplate.Execute(w, nil)
	http.ServeFile(w, r, "web/index.html")

	// conn, err := upgrader.Upgrade(w, r, nil)

	// if err != nil {
	// 	log.Printf("Error home upgrading connection: %v\n", err)
	// 	return nil, err
	// }

	// defer conn.Close()

	// return conn, nil
}
