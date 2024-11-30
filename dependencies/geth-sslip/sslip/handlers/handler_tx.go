package handlers

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/poc-client/utils/cryptoutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/sslip/manager"
	"github.com/ethereum/go-ethereum/sslip/mpt"
	"github.com/ethereum/go-ethereum/sslip/resmsg"
	"github.com/gorilla/websocket"
)

func handler_tx(clientID string, body string, m *manager.Manager, conn *websocket.Conn, mt int) error {

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
	reqChannelId := m.GetClientChannelID(clientID)
	if reqChannelId == req.ChannelID {
		log.Println("Channel ID matched", req.ChannelID)
	}

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

	log.Println()
	log.Printf("----------------Generating Proof--------------------")

	// // To delete later
	// proof, idx, txHash, tx30 := FuncTestTransactionRootAndProof()
	// // blockNr := big.NewInt(int64(10467135))
	// txRootHash := common.HexToHash("bb345e208bda953c908027a45aa443d6cab6b8d2fd64e83ec52f1008ddeafa58")

	// ----end

	// To add later: Generate the proof of the transaction
	blockHash := txReceipt.BlockHash
	block, _ := bcClient.BlockByHash(context.Background(), blockHash)
	blockNr := block.Number()
	txHash := txReceipt.TxHash
	txRootHash := block.TxHash()
	proof, idx := generateProof(block, txHash)

	// -----Log Helper: for block information
	// currentBlockHeader, _ := bcClient.HeaderByNumber(context.Background(), blockNr)
	// ComputeBlockHash(currentBlockHeader)
	// PrintBlockInfo(currentBlockHeader)
	// log.Println(BhRlpBytes(currentBlockHeader))
	// -----end

	log.Println("*************************Generation end!*************************")

	log.Println()
	log.Printf("----------------Verifying Proof--------------------")
	if proof == nil {
		log.Println("Error: unable to generate the proof")
	} else {
		// res := verifyProofTmpt(txHash, proof, idx, tx30)
		res := verifyProof(txHash, proof, blockNr, idx)

		// -----Log Helper: print proof information to submit on-chain
		// PrintProofToSubmit(proof, idx, txRootHash)
		// -----end
		log.Println("Proof Verification: ", res)

	}
	log.Println("*************************Proof verification end!*************************")

	log.Println()
	log.Printf("----------------Generating response!!--------------------")

	channelId := m.GetClientChannelID(clientID)
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
		Signature:          nil,
		RootHash:           txRootHash,
	}

	responseBodyHash := responseMsg.Keccak256Hash()
	sig := cryptoutil.SignHash(m.PrivateKey, responseBodyHash)
	responseMsg.Signature = sig

	// -----Log Helper: print response message to verify hash values
	// log.Println("channelId: ", channelId)
	// log.Println("RootHash: ", txRootHash)
	// log.Println("SignedReqBody:", hex.EncodeToString(req.SignedReqBody))
	// log.Println("Signature:", hex.EncodeToString(sig))
	// log.Println("resHash: ", responseBodyHash)
	// log.Println("reqHash: ", reqHash)
	// -----end

	// -----Log Helper: print response message bytes
	fmt.Println("-=-=-=-=-= Now print response message bytes -=-=-=-=-=-=")
	log.Println(responseMsg.RlpBytes())

	fmt.Println("-=-=-=-=-= Now print request body bytes -=-=-=-=-=-=")
	reqBodyBytesString := req.RequestBodyRlpBytes()
	log.Println(reqBodyBytesString)

	fmt.Println("*********************************************************************")
	// -----end

	log.Println("-=-=-=-=-=-= Now print request payment bytes -=-=-=-=-=-=")
	log.Println("Payment body byte: ", req.PaymentBodyRlpBytes())
	log.Println("*********************************************************************")

	_ = conn.WriteMessage(mt, responseMsg.Bytes())

	return nil
}

func generateProof(block *types.Block, txHash common.Hash) (mpt.Proof, []byte) {
	txs := block.Transactions()
	idx := -1
	for index, tx := range txs {
		if tx.Hash() == txHash {
			idx = index
		}
	}
	if idx < 0 {
		return nil, nil
	}

	mptTrie := mpt.NewTrie()
	for i, tx := range txs {
		key, _ := rlp.EncodeToBytes(uint(i))

		transaction := mpt.FromEthTransaction(tx)

		rlp, _ := transaction.GetRLP()

		mptTrie.Put(key, rlp)
	}

	// generate the proof and verify it
	key, _ := rlp.EncodeToBytes(uint(idx))
	proof, _ := mptTrie.Prove(key)
	// proofSize := len(proof.Serialize()[0])
	// log.Println("proofSize: ", proofSize)
	// fmt.Printf("proof: %x, found: %v\n", proof, found)

	return proof, key
}

func verifyProofTmpt(txHash common.Hash, proof mpt.Proof, key []byte, tx *types.Transaction) bool {

	txRootHash, _ := hex.DecodeString("bb345e208bda953c908027a45aa443d6cab6b8d2fd64e83ec52f1008ddeafa58")
	txRLP, _ := rlp.EncodeToBytes(tx)
	txProofRLP, _ := mpt.VerifyProof(txRootHash[:], key, proof)
	// log.Println("txProofRLP: ", txProofRLP)
	// log.Println("txRLP: ", txRLP)
	// log.Println("proof: ", proof.Serialize())

	return bytes.Equal(txRLP, txProofRLP)
}
