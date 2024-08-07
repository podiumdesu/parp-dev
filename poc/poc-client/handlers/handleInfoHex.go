package handlers

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"poc-server/resmsg"
)

func HandleInfoHex(msg []byte) error {
	var serverMsg resmsg.ServerMsg
	err := json.Unmarshal(msg, &serverMsg)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println(hex.EncodeToString(serverMsg.Info))
	return nil
}
