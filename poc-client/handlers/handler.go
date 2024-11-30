package handlers

import (
	"encoding/json"
	"log"
	"poc-client/client"
	"poc-client/utils/cryptoutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/sslip/resmsg"
)

func HandleMesssage(msg []byte, client *client.Client) error {
	var serverMsg resmsg.ServerMsg

	err := json.Unmarshal(msg, &serverMsg)
	if err != nil {
		log.Fatal(err)
		return err
	}

	switch serverMsg.Type {
	case "info":
		HandleInfo(msg)
	case "info-hex":
		HandleInfoHex(msg)
	case "channelId":
		channelId := common.BytesToHash(serverMsg.Info)
		client.SetChanId(channelId)
	case "open-chan":
		HandleOpenChan(msg, client)
	case "HANDSHAKE-CONFIRMED":
		err := HandleHandshakeConfirm(msg, client)
		if err != nil {
			log.Fatal(err)
			return err
		}
		log.Println("Server public key address has been set: ", cryptoutil.PubkeyToHexAddr(client.ServerPublicKey))

	case "response":
		HandleResponse(msg, client)

	case "responseSP":
		HandleResponseSP(msg, client)
	default:
		log.Println("Unrecognized message type: ", serverMsg.Type)
	}
	return nil
}
