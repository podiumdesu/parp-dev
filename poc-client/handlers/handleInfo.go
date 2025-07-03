package handlers

import (
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/parp/resmsg"
)

func HandleInfo(msg []byte) error {
	var infoMsg resmsg.ServerMsg
	err := json.Unmarshal(msg, &infoMsg)
	if err != nil {
		log.Println("Unmarshal error: ", err)
		return err
	}
	log.Println("[INFO]: ", string(infoMsg.Info))
	return nil
}
