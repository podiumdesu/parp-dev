package handlers

import (
	"encoding/json"
	"log"
	"poc-client/client"
	"poc-client/utils/cryptoutil"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/parp/resmsg"
)

func HandleOpenChan(msg []byte, client *client.Client) error {

	var resMsg resmsg.ResponseMsg

	err := json.Unmarshal(msg, &resMsg)
	if err != nil {
		log.Println("Unmarshal error: ", err)
		return err
	}

	log.Println("[CHANOPEN] Response Message: ", resMsg)

	resMsgBodyHash := resMsg.Keccak256Hash()
	res := cryptoutil.VerifyHash(crypto.FromECDSAPub(client.ServerPublicKey), resMsgBodyHash, resMsg.Signature)
	log.Println("[CHANOPEN] Verify Response signature:", res)

	client.SetChanId(resMsg.ChannelId)
	log.Println("[CHANOPEN] ChannelID: ", client.ChannelID)
	return nil
}
