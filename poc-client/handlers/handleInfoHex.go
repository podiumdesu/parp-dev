package handlers

import (
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/parp/resmsg"

	"github.com/ethereum/go-ethereum/common"
)

func HandleInfoHex(msg []byte) error {
	var serverMsg resmsg.ServerMsg
	err := json.Unmarshal(msg, &serverMsg)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("[INFOHEX]: ", common.BytesToHash(serverMsg.Info))
	return nil
}
