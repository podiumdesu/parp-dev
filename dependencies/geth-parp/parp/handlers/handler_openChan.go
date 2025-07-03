package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/parp/manager"
	"github.com/ethereum/go-ethereum/parp/resmsg"
	"github.com/ethereum/go-ethereum/poc-client/utils/cryptoutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/websocket"
)

func handler_openChan(clientID string, body string, m *manager.Manager, conn *websocket.Conn, mt int) error {

	client := m.GetClient(clientID)

	req, _ := unmarshalRequest(body)
	requestByte := req.ReqByte

	log.Println("----------------SignTx from clientID ", clientID, "---------------------")
	_, reqHash, sigCheckMsg, err := verifyReqSignature(m, clientID, req)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Send the result of the signature check to the client
	conn.WriteMessage(mt, sigCheckMsg.Bytes())

	log.Println("*************************Pre check passed!*************************")

	// 3. Send the tx to the network

	bcClient, _ := client.ConnectToBlockchain()

	txSend := new(types.Transaction)
	rlp.DecodeBytes(requestByte, &txSend)

	err = bcClient.SendTransaction(context.Background(), txSend)
	if err != nil {
		log.Fatal(err)
	}

	log.Println()
	log.Printf("----------------Submited Tx: %s--------------------", txSend.Hash().Hex())

	msg := resmsg.ServerMsg{
		Type: "info-hex",
		Info: txSend.Hash().Bytes(),
	}
	client.Send(msg.Bytes())

	var txReceipt *types.Receipt
	for txReceipt == nil {
		// Query the transaction receipt
		txReceipt, err = bcClient.TransactionReceipt(context.Background(), txSend.Hash())
		if err != nil {
			log.Println("Waiting for transaction to be mined...")
			time.Sleep(5 * time.Second) // Adjust the sleep duration based on expected block time
		}
	}
	log.Printf("Transaction mined in block %d", txReceipt.BlockNumber.Uint64())

	log.Println("*************************Block mining end!*************************")
	// By directly monitor the mempool

	if txReceipt.Status == types.ReceiptStatusFailed {
		log.Printf("❌ Transaction %s failed", txSend.Hash().Hex())
		msg := resmsg.ServerMsg{
			Type: "info-hex",
			Info: []byte("Transaction failed"),
		}
		client.Send(msg.Bytes())
		return err
	}

	// Retrieve result
	contractAbi, err := abi.JSON(strings.NewReader(string(contractABI)))
	if err != nil {
		log.Println("Error parsing contract ABI:", err)
		return err
	}
	var channelId common.Hash
	for _, vLog := range txReceipt.Logs {
		fmt.Printf("Log Address: %s\n", vLog.Address.Hex())
		if len(vLog.Topics) > 0 {
			eventName, err := contractAbi.EventByID(vLog.Topics[0])
			if err != nil {
				log.Println("Error finding event name:", err)
				continue
			}
			fmt.Printf("Event Name: %s\n", eventName.Name)

			// var results []interface{}
			if len(vLog.Topics) > 1 {

				channelId = common.BytesToHash(vLog.Topics[1].Bytes())
				m.SetClientChannelID(clientID, channelId)
				fmt.Printf("Channel ID: %x\n", channelId)
				// fmt.Printf("Channel ID: %s\n", m.GetClientChannelID(clientID))
				// channelIdString, _ := hex.DecodeString(channelId.Hex()[2:])
				msg = resmsg.ServerMsg{
					Type: "channelId",
					Info: channelId.Bytes(),
				}
				client.Send(msg.Bytes())
			}
		}
	}

	log.Println()
	log.Printf("----------------Generating Proof--------------------")

	// To delete later
	// proof, idx, txHash, tx30 := FuncTestTransactionRootAndProof()
	// txRootHash := common.HexToHash("bb345e208bda953c908027a45aa443d6cab6b8d2fd64e83ec52f1008ddeafa58")
	// printProofToSubmit(proof, idx, txRootHash)

	// log.Println()
	// log.Printf("----------------Verifying Proof--------------------")
	// if proof == nil {
	// 	log.Println("Error: unable to generate the proof")
	// } else {
	// 	res := verifyProofTmpt(txHash, proof, idx, tx30)

	// 	log.Println("Proof Verification: ", res)

	// }
	// log.Println("*************************Proof verification end!*************************")

	// ----end

	// To add later: Generate the proof of the transaction
	blockHash := txReceipt.BlockHash
	block, _ := bcClient.BlockByHash(context.Background(), blockHash)
	blockNr := block.Number()
	txHash := txReceipt.TxHash
	// txRootHash := block.TxHash()
	proof, idx := generateProof(block, txHash)

	log.Println("*************************Generation end!*************************")

	log.Printf("----------------Verifying Proof--------------------")
	if proof == nil {
		log.Println("Error: unable to generate the proof")
	} else {
		res := verifyProof(txHash, proof, blockNr, idx)

		log.Println("Proof Verification: ", res)

	}
	log.Println("*************************Proof verification end!*************************")

	// Send the response to the client

	log.Println()
	log.Printf("----------------Generating response!!--------------------")

	// responseBody := resmsg.ResponseBody{
	// 	SignedReqBody: req.SignedReqBody,
	// 	Proof:         proof.CustomRLPSerialize(),
	// 	TxHash:        txHash,
	// 	TxIdx:         idx,
	// }

	responseMsg := resmsg.ResponseMsg{
		Type:               "open-chan",
		ChannelId:          channelId,
		Amount:             req.Amount,
		ReqBodyHash:        reqHash,
		SignedReqBody:      req.SignedReqBody,
		CurrentBlockHeight: blockNr,
		ReturnValue:        txReceipt.Bloom.Bytes(),
		Proof:              proof.CustomRLPSerialize(),
		TxHash:             txHash,
		TxIdx:              idx,
		Signature:          []byte(""),
	}
	sig := cryptoutil.SignHash(m.PrivateKey, responseMsg.Keccak256Hash())
	responseMsg.Signature = sig

	// log.Println("SignedReqBody:", hex.EncodeToString(req.SignedReqBody))
	// log.Println("Signature:", hex.EncodeToString(sig))
	// log.Println("resHash: ", responseMsg.Keccak256Hash())
	// log.Println("reqHash: ", reqHash)

	// fmt.Println("-=-=-=-=-= Now print response message bytes -=-=-=-=-=-=")
	// log.Println(responseMsg.RlpBytes())

	// fmt.Println("-=-=-=-=-= Now print request body bytes -=-=-=-=-=-=")

	// reqBodyBytesString := req.RequestBodyRlpBytes()
	// log.Println(reqBodyBytesString)

	// fmt.Println("*********************************************************************")

	_ = conn.WriteMessage(mt, responseMsg.Bytes())

	return nil
}
