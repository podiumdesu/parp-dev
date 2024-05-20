package protocol

import (
	"encoding/json"
	"log"
	"poc-client/client"
	"poc-client/msg/request"

	"github.com/ethereum/go-ethereum/common"
)

// PoC client generates the request

func GenerateRequest(c *client.Client, ch int, amount uint, reqByte []byte, blockHash common.Hash) []byte {
	// log.Println(reqByte)
	requestBody := request.ReqBody{
		ChannelID:      ch,
		Amount:         amount,
		ReqByte:        reqByte,
		LocalBlockHash: blockHash,
	}

	signedReqBody := c.Sign(requestBody.PreHashByte())

	paymentBody := request.PaymentBody{
		ChannelID: ch,
		Amount:    amount,
	}

	signedPayBody := c.Sign(paymentBody.PreHashByte())

	// log.Println("_--------")
	// log.Println(hex.EncodeToString(paymentBody.PreHashByte()))
	// log.Println(hex.EncodeToString(signedPayBody))
	// log.Println("_--------")
	request := request.RequestMsg{
		ChannelID:         ch,
		Type:              0,
		Amount:            amount,
		ReqByte:           reqByte,
		LocalBlockHash:    blockHash,
		SignedReqBody:     signedReqBody,
		SignedPaymentBody: signedPayBody,
	}

	c.Amount = amount

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
	}

	// log.Println("Request is ready to be sent: ", string(jsonRequest))

	// log.Println("Now I want to check the verification")
	// log.Println("Verification: ", c.Verify(requestBody.HashByte(), signedReqBody))
	// log.Println("Verification: ", cryptoutil.Verify(c.PubKeyBytes(), requestBody.HashByte(), signedReqBody))
	return jsonRequest
}
