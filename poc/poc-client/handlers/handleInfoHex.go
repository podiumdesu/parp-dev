package handlers

import (
	"encoding/json"
	"log"
	"poc-server/resmsg"

	"github.com/ethereum/go-ethereum/common"
)

func HandleInfoHex(msg []byte) error {
	var serverMsg resmsg.ServerMsg
	err := json.Unmarshal(msg, &serverMsg)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println(common.BytesToHash(serverMsg.Info))
	return nil
}
