package handlers

import (
	"encoding/json"
	"log"
	"poc-client/client"
	"poc-client/utils/cryptoutil"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/sslip/resmsg"
)

func HandleHandshakeConfirm(msg []byte, client *client.Client) error {
	var hsMsg resmsg.HandshakeMsg
	err := json.Unmarshal(msg, &hsMsg)
	if err != nil {
		return err
	}

	serverPubKeyByte := hsMsg.ServerPublicKeyB
	serverPubKeyECDSA, err := crypto.UnmarshalPubkey(serverPubKeyByte)
	if err != nil {
		return err
	}

	client.ServerPublicKey = serverPubKeyECDSA
	log.Println("Getting public key: ", cryptoutil.PubkeyToHexAddr(client.ServerPublicKey))

	return nil
}
