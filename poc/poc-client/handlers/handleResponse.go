package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math/big"
	"poc-client/client"
	"poc-client/utils/cryptoutil"
	"poc-server/mpt"
	"poc-server/resmsg"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func HandleResponse(msg []byte, client *client.Client) error {

	var resMsg resmsg.ResponseMsg

	log.Printf("Size of the Tx response: %d bytes", len(msg))
	err := json.Unmarshal(msg, &resMsg)
	if err != nil {
		log.Println("Unmarshal error: ", err)
		return err
	}

	resMsgBodyHash := resMsg.BodyHashBytes()
	res := cryptoutil.Verify(crypto.FromECDSAPub(client.ServerPublicKey), resMsgBodyHash, resMsg.Signature)
	log.Println("Verify Response signature:", res)

	proof, err := mpt.DeserializeProof(resMsg.Proof)
	if err != nil {
		log.Println("Error deserializing proof: ", err)
	}
	res = verifyProof(client, resMsg.TxHash, proof, resMsg.CurrentBlockHeight, uint32(resMsg.TxIdx))

	log.Println("Proof Verification: ", res)
	return nil
}

func verifyProof(client *client.Client, txHash common.Hash, proof mpt.Proof, blockNr *big.Int, idx uint32) bool {
	wsEndpoint := client.BcWsEndpoint
	bcClient, err := ethclient.Dial(wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	// query the block information
	block, _ := bcClient.HeaderByNumber(context.Background(), blockNr)
	txRootHash := block.TxHash
	tx, _, _ := bcClient.TransactionByHash(context.Background(), txHash)
	txRLP, _ := rlp.EncodeToBytes(tx)
	key, _ := rlp.EncodeToBytes(uint32(idx))
	txProofRLP, _ := mpt.VerifyProof(txRootHash[:], key, proof)
	log.Println("txProofRLP: ", txProofRLP)
	log.Println("txRLP: ", txRLP)
	log.Println("proof: ", proof.Serialize())

	return bytes.Equal(txRLP, txProofRLP)
}
