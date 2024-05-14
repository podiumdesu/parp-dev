package protocol

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Body)
	// defer conn.Close()

	// fmt.Println("This is one of the handlers... Currently I dont know how to write the code")

	// buffer := make([]byte, 1024)

	// for {
	// 	n, err := conn.Read(buffer)

	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 		return
	// 	}

	// 	fmt.Printf("Received from %s: %s \n", conn.RemoteAddr(), buffer[:n])

	// 	conn.Write([]byte("Messsage received. Server"))
	// }
}
