package handlers

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/poc-client/msg/request"
	"github.com/ethereum/go-ethereum/poc-client/utils/cryptoutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/sslip/manager"
	"github.com/ethereum/go-ethereum/sslip/mpt"
	"github.com/ethereum/go-ethereum/sslip/resmsg"
	"github.com/gorilla/websocket"
)

func handler_tx(clientID string, body string, m *manager.Manager, conn *websocket.Conn, mt int) error {

	client := m.GetClient(clientID)

	log.Println("----------------SignTx from clientID ", clientID, "---------------------")

	// 1. Unmarshal request body
	var req request.RequestMsg
	err := json.Unmarshal([]byte(body), &req)

	// err := rlp.DecodeBytes([]byte(body), &req)
	if err != nil {
		log.Fatal("Unmarshal error: ", err)
		return err
	}
	jsonReq, _ := json.Marshal(req)
	log.Println("Request: ", string(jsonReq))

	requestByte := req.ReqByte
	// log.Println(requestByte)

	// 2. Verify the signature
	log.Println("SignedReqBody after verification:", hex.EncodeToString(req.SignedReqBody))

	sigFlag, reqHash := m.VerifyRequestWithSig(clientID, req)
	log.Println("SignedReqBody after verification:", hex.EncodeToString(req.SignedReqBody))
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

	log.Println("*************************Pre check passed!*************************")

	// Send the result of the signature check to the client
	conn.WriteMessage(mt, msg.Bytes())

	// 3. Send the tx to the network
	wsEndpoint := "ws://localhost:8100"
	bcClient, err := ethclient.Dial(wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	txSend := new(types.Transaction)
	rlp.DecodeBytes(requestByte, &txSend)

	err = bcClient.SendTransaction(context.Background(), txSend)
	if err != nil {
		log.Fatal(err)
	}

	log.Println()
	log.Printf("----------------Submited Tx: %s--------------------", txSend.Hash().Hex())

	msg = resmsg.ServerMsg{
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

	// Retrieve result
	const contractABI = `[{"inputs":[{"indexed":true,"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"ChannelOpened","type":"event"},{"inputs":[{"internalType":"address","name":"addr","type":"address"}],"name":"balance","outputs":[{"internalType":"uint256","name":"bal","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"closeChan","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"confirmClosure","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"}],"name":"greeting","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"senderDeposit","type":"uint256"}],"name":"openChan","outputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"paychanCheck","outputs":[{"components":[{"internalType":"bytes32","name":"id","type":"bytes32"},{"internalType":"address payable","name":"sender","type":"address"},{"internalType":"address payable","name":"recipient","type":"address"},{"internalType":"uint256","name":"senderDeposit","type":"uint256"},{"internalType":"uint256","name":"startTime","type":"uint256"},{"internalType":"uint256","name":"status","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"},{"internalType":"uint256","name":"disputeStartTime","type":"uint256"},{"internalType":"uint256","name":"disputeDuration","type":"uint256"},{"internalType":"bool","name":"senderConfirm","type":"bool"},{"internalType":"bool","name":"recipientConfirm","type":"bool"}],"internalType":"struct paychan.PayChan","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"paychanSelectedArguments","outputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"rec","type":"address"},{"internalType":"uint256","name":"status","type":"uint256"},{"internalType":"uint256","name":"senderB","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	contractAbi, err := abi.JSON(strings.NewReader(string(contractABI)))
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
				m.SetClientChannelID(clientID, string(channelId.Hex()))
				fmt.Printf("Channel ID: %s\n", channelId.Hex())
				fmt.Printf("Channel ID: %s\n", m.GetClientChannelID(clientID))
				msg = resmsg.ServerMsg{
					Type: "info-hex",
					Info: channelId.Bytes(),
				}
				client.Send(msg.Bytes())
			}
		}
	}

	log.Println()
	log.Printf("----------------Generating Proof--------------------")

	// To delete later
	proof, idx, txHash, tx30 := FuncTestTransactionRootAndProof()
	// blockNr := big.NewInt(int64(10467135))
	txRootHash := common.HexToHash("bb345e208bda953c908027a45aa443d6cab6b8d2fd64e83ec52f1008ddeafa58")
	printProofToSubmit(proof, idx, txRootHash)

	// ----end

	// To add later: Generate the proof of the transaction
	blockHash := txReceipt.BlockHash
	block, _ := bcClient.BlockByHash(context.Background(), blockHash)
	blockNr := block.Number()
	// txHash := txReceipt.TxHash
	// txRootHash := block.TxHash()
	// proof, idx := generateProof(block, txHash)

	log.Println("*************************Generation end!*************************")

	log.Println()
	log.Printf("----------------Verifying Proof--------------------")
	if proof == nil {
		log.Println("Error: unable to generate the proof")
	} else {
		res := verifyProofTmpt(txHash, proof, idx, tx30)

		log.Println("Proof Verification: ", res)

	}
	log.Println("*************************Proof verification end!*************************")

	// Send the response to the client

	log.Println()
	log.Printf("----------------Generating response!!--------------------")

	responseBody := resmsg.ResponseBody{
		SignedReqBody: req.SignedReqBody,
		Proof:         proof.CustomRLPSerialize(),
		TxHash:        txHash,
		TxIdx:         idx,
	}

	sig := cryptoutil.SignHash(m.PrivateKey, responseBody.Keccak256Hash())

	responseMsg := resmsg.ResponseMsg{
		Type:               "response",
		ChannelId:          channelId,
		Amount:             req.Amount,
		ReqBodyHash:        reqHash,
		SignedReqBody:      req.SignedReqBody,
		CurrentBlockHeight: blockNr,
		ReturnValue:        txReceipt.Bloom.Bytes(),
		Proof:              proof.CustomRLPSerialize(),
		TxHash:             txHash,
		TxIdx:              idx,
		Signature:          sig,
	}

	log.Println("SignedReqBody:", hex.EncodeToString(req.SignedReqBody))
	log.Println("Signature:", hex.EncodeToString(sig))
	log.Println("resHash: ", responseBody.Keccak256Hash())
	log.Println("reqHash: ", reqHash)

	fmt.Println("-=-=-=-=-= Now print response message bytes -=-=-=-=-=-=")
	log.Println(responseMsg.RlpBytes())

	fmt.Println("-=-=-=-=-= Now print request body bytes -=-=-=-=-=-=")
	reqBody := request.ReqBody{
		// ChannelID:      req.ChannelID,
		Amount:         req.Amount,
		LocalBlockHash: req.LocalBlockHash,
		ReqByte:        req.ReqByte,
	}
	reqBodyBytesString := reqBody.RlpBytes()
	log.Println(reqBodyBytesString)

	fmt.Println("*********************************************************************")

	_ = conn.WriteMessage(mt, responseMsg.Bytes())

	return nil
}

func verifyProofTmpt(txHash common.Hash, proof mpt.Proof, key []byte, tx *types.Transaction) bool {
	// wsEndpoint := "ws://localhost:8100"
	// bcClient, err := ethclient.Dial(wsEndpoint)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// the old one
	// block, _ := bcClient.HeaderByNumber(context.Background(), blockNr)
	// txRootHash := block.TxHash
	// tx, _, _ := bcClient.TransactionByHash(context.Background(), txHash)
	// txRLP, _ := rlp.EncodeToBytes(tx)
	// txProofRLP, _ := mpt.VerifyProof(txRootHash[:], key, proof)
	// log.Println("txProofRLP: ", txProofRLP)
	// log.Println("txRLP: ", txRLP)
	// log.Println("proof: ", proof.Serialize())

	// return bytes.Equal(txRLP, txProofRLP)

	// To match the new one

	txRootHash, _ := hex.DecodeString("bb345e208bda953c908027a45aa443d6cab6b8d2fd64e83ec52f1008ddeafa58")
	txRLP, _ := rlp.EncodeToBytes(tx)
	txProofRLP, _ := mpt.VerifyProof(txRootHash[:], key, proof)
	// log.Println("txProofRLP: ", txProofRLP)
	// log.Println("txRLP: ", txRLP)
	// log.Println("proof: ", proof.Serialize())

	return bytes.Equal(txRLP, txProofRLP)
}
