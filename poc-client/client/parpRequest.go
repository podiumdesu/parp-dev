package client

import (
	"encoding/json"
	"log"
	"poc-client/msg/request"

	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) generateParpRequest(chanID common.Hash, amount uint, reqByte []byte, blockHash common.Hash) []byte {

	request := request.RequestMsg{
		ChannelID:         chanID,
		Type:              0,
		Amount:            amount,
		ReqByte:           reqByte,
		LocalBlockHash:    blockHash,
		SignedReqBody:     nil,
		SignedPaymentBody: nil,
	}

	// Sign Request Body and Payment
	requestBodyHash := request.BodyKeccak256Hash()
	signedReqBody := c.Sign(requestBodyHash.Bytes())
	paymentHash := request.PaymentKeccak256Hash()
	signedPayBody := c.Sign(paymentHash.Bytes())

	request.SignedReqBody = signedReqBody
	request.SignedPaymentBody = signedPayBody

	log.Println()
	log.Println("-=-=-=-=-= Now print request body bytes -=-=-=-=-=-=")
	log.Println(request.RequestBodyRlpBytes())
	log.Println("*********************************************************************")

	// log.Println("Body message Hash: ", requestBodyHash.Hex())
	// log.Println("Payment message Hash: ", paymentHash.Hex())
	// log.Println("Payment message sig: ", hex.EncodeToString(signedPayBody))
	// log.Println("light client signature over request body: ", hex.EncodeToString(signedReqBody))
	// log.Println("amount: ", amount)
	// log.Println("local block hash: ", blockHash)
	// log.Println("Req byte: ", hex.EncodeToString(reqByte))
	// log.Println("Payment body byte: ", request.PaymentBodyRlpBytes())

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
	}

	return jsonRequest
}
