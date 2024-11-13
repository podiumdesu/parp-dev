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

	proof := resMsg.Proof // get the proof string
	if err != nil {
		log.Println("Error deserializing proof: ", err)
	}
	res = verifySerializedProof(client, resMsg.TxHash, proof, resMsg.CurrentBlockHeight, resMsg.TxIdx)

	log.Println("Proof Verification: ", res)

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

	log.Println("The bytes received: ", proofBytes)

	proof, _ := mpt.DeserializeProof(proofBytes)

	log.Println("Deserialized Proof:", proof)
	// log.Println("GGSSGSDGS")
	// log.Println(proof)
	// log.Println("GGSSGSDGS")
	txProofRLP, _ := mpt.VerifyProof(txRootHash[:], key, proof)
	log.Println("txProofRLP: ", txProofRLP)
	log.Println("txRLP: ", txRLP)
	log.Println("proof: ", proofBytes)

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

	return bytes.Equal(txRLP, txProofRLP)

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
