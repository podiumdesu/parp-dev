package handlers

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"poc-client/client"
	"poc-client/utils/cryptoutil"

	"github.com/ethereum/go-ethereum/sslip/mpt"

	"github.com/ethereum/go-ethereum/sslip/resmsg"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func HandleResponse(msg []byte, client *client.Client) error {

	var resMsg resmsg.ResponseMsg

	err := json.Unmarshal(msg, &resMsg)
	if err != nil {
		log.Println("Unmarshal error: ", err)
		return err
	}

	resMsgBodyHash := resMsg.Keccak256Hash()
	res := cryptoutil.VerifyHash(crypto.FromECDSAPub(client.ServerPublicKey), resMsgBodyHash, resMsg.Signature)
	log.Println("[RESPONSE] Verify Response signature:", res)

	proof := resMsg.Proof // get the proof string
	if err != nil {
		log.Println("[RESPONSE] Error deserializing proof: ", err)
	}
	res = verifySerializedProof(client, resMsg.TxHash, proof, resMsg.CurrentBlockHeight, resMsg.TxIdx)

	log.Println("[RESPONSE] Proof Verification: ", res)

	fmt.Println("-=-=-=-=-= Now print response message bytes -=-=-=-=-=-=")
	log.Println(resMsg.RlpBytes())
	fmt.Println("*********************************************************************")

	// fmt.Println("*********************************************************************")

	return nil
}

// For serialized proof string

func verifySerializedProof(client *client.Client, txHash common.Hash, proofBytes [][]byte, blockNr *big.Int, key []byte) bool {
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

	proof, _ := mpt.DeserializeProof(proofBytes)

	txProofRLP, _ := mpt.VerifyProof(txRootHash[:], key, proof)

	return bytes.Equal(txRLP, txProofRLP)

}

func PrintProof(txRootHash common.Hash, proofBytes [][]byte, key []byte) {
	log.Printf("Root Hash: %s\n", txRootHash.Hex())
	// 2. Print the serialized proof as a bytes array (proof is []string).
	fmt.Print("Proof (bytes[] proof): [")
	for i, node := range proofBytes {
		if i > 0 {
			fmt.Print(", ")
		}
		// Add "0x" prefix for each proof element
		fmt.Printf("\"0x%s\"", hex.EncodeToString(node))
	}
	fmt.Println("]")
	// 3. Print the keys as a bytes array.
	// RLP encode the index as the key (assuming resMsg.TxIdx is the index)
	// Convert the RLP-encoded key to hex format.
	keyHex := "0x" + hex.EncodeToString(key)

	// Print the key as bytes[] for Remix.
	fmt.Printf("Keys (bytes[] keys): [\"%s\"]\n", keyHex)
}
