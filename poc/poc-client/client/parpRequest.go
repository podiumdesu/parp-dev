package client

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"poc-client/msg/request"

	"github.com/ethereum/go-ethereum/common"
)

func (c *Client) generateParpRequest(ch uint32, amount uint, reqByte []byte, blockHash common.Hash) []byte {
	log.Println(reqByte)

	// blockHash = common.HexToHash("0xeb96a0494c2e86c4597beb442028b1f490e69aa8ec2a80e8c837204dc642dcf0")
	// reqByte = []byte("This is a fucking test")
	requestBody := request.ReqBody{
		// ChannelID:      ch,
		Amount:         amount,
		LocalBlockHash: blockHash,
		ReqByte:        reqByte,
	}
	log.Println("request body: ", requestBody.RlpBytes())
	// log.Println(requestBody.RlpBytes())
	// signedReqBody := c.Sign(requestBody.PreHashByte())
	signedReqBody := c.Sign(requestBody.Keccak256Hash().Bytes())

	log.Println("message Hash: ", requestBody.Keccak256Hash().Hex())
	log.Println("light client signature over request body: ", hex.EncodeToString(signedReqBody))

	log.Println("amount: ", amount)
	log.Println("local block hash: ", blockHash)
	log.Println("Req byte: ", hex.EncodeToString(reqByte))

	paymentBody := request.PaymentBody{
		// ChannelID: ch,
		Amount: amount,
	}

	signedPayBody := c.Sign(paymentBody.PreHashByte())

	// log.Println("_--------")
	// log.Println(hex.EncodeToString(paymentBody.PreHashByte()))
	// log.Println(hex.EncodeToString(signedPayBody))
	// log.Println("_--------")
	request := request.RequestMsg{
		// ChannelID:         ch,
		Type:              0,
		Amount:            amount,
		ReqByte:           reqByte,
		LocalBlockHash:    blockHash,
		SignedReqBody:     signedReqBody,
		SignedPaymentBody: signedPayBody,
	}

	c.Amount = amount

	jsonRequest, err := json.Marshal(request)
	// jsonRequest, err := rlp.EncodeToBytes(request)
	if err != nil {
		log.Println(err)
	}

	// log.Println("Request is ready to be sent: ", string(jsonRequest))

	// log.Println("Now I want to check the verification")
	// log.Println("Verification: ", c.Verify(requestBody.HashByte(), signedReqBody))
	// log.Println("Verification: ", cryptoutil.Verify(c.PubKeyBytes(), requestBody.HashByte(), signedReqBody))
	return jsonRequest
}
