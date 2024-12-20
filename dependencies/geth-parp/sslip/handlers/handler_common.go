package handlers

import (
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/poc-client/msg/request"
	"github.com/ethereum/go-ethereum/sslip/manager"
	"github.com/ethereum/go-ethereum/sslip/resmsg"
)

func verifyReqSignature(m *manager.Manager, clientID string, req request.RequestMsg) (bool, common.Hash, resmsg.ServerMsg, error) {

	sigFlag, reqHash := m.VerifyRequestWithSig(clientID, req)
	// log.Println("SignedReqBody after verification:", hex.EncodeToString(req.SignedReqBody))

	var msg resmsg.ServerMsg
	if sigFlag {
		log.Println("PASS: Signature verified")
		msg = resmsg.ServerMsg{
			Type: "info",
			Info: []byte("SigCheck: Passed"),
		}
	} else {
		msg = resmsg.ServerMsg{
			Type: "info",
			Info: []byte("SigCheck: WRONG signature"),
		}
	}

	return sigFlag, reqHash, msg, nil
}

func unmarshalRequest(body string) (request.RequestMsg, error) {
	// 1. Unmarshal request body
	var req request.RequestMsg
	err := json.Unmarshal([]byte(body), &req)
	if err != nil {
		log.Fatal("Unmarshal error: ", err)
		return req, err
	}
	// jsonReq, _ := json.Marshal(req)
	// log.Println("Request: ", string(jsonReq))
	return req, nil
}
