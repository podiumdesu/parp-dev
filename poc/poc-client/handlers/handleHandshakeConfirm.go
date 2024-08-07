package handlers

import (
	"encoding/json"
	"poc-client/client"
	"poc-server/resmsg"

	"github.com/ethereum/go-ethereum/crypto"
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

	return nil
}
