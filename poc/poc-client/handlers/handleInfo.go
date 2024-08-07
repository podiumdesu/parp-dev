package handlers

import (
	"encoding/json"
	"log"
	"poc-server/resmsg"
)

func HandleInfo(msg []byte) error {
	var infoMsg resmsg.ServerMsg
	err := json.Unmarshal(msg, &infoMsg)
	if err != nil {
		log.Println("Unmarshal error: ", err)
		return err
	}
	log.Println(string(infoMsg.Info))
	return nil
}
