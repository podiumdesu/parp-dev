package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"poc-client/client"
	utils "poc-client/utils/common"
	"poc-client/utils/cryptoutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/sslip/mpt"
	"github.com/ethereum/go-ethereum/sslip/resmsg"
)

func HandleResponseSP(msg []byte, client *client.Client) error {
	var resMsg resmsg.ResponseSPMsg
	err := json.Unmarshal(msg, &resMsg)
	if err != nil {
		log.Println("Unmarshal error: ", err)
		return err
	}
	resBodyHash := resMsg.Keccak256Hash()
	res := cryptoutil.VerifyHash(crypto.FromECDSAPub(client.ServerPublicKey), resBodyHash, resMsg.Signature)
	log.Println("[RES-SP] Verify Response signature:", res)

	proof, _ := mpt.DeserializeProof(resMsg.Proof)
	if err != nil {
		log.Println("Error deserializing proof: ", err)
	}

	bcClient, err := client.ConnectToBlockchain()
	if err != nil {
		return err
	}
	blockHeader, _ := bcClient.HeaderByNumber(context.Background(), resMsg.CurrentBlockNr)
	log.Println()
	log.Println("[RES-SP] Print block header bytes: ")
	log.Println(utils.BhRlpBytes(blockHeader))

	result, _ := verifySPProof(client, proof, resMsg.CurrentBlockNr, resMsg.Address)
	log.Println("[RES-SP] Proof Verification: ", result)

	fmt.Println("-=-=-=-=-= [RES-SP] Now print response message bytes -=-=-=-=-=-=")
	log.Println(resMsg.RlpBytes())
	fmt.Println("*********************************************************************")

	return nil
}

// func verifySPProof(c *client.Client, proofTrie mpt.Proof, blockNr *big.Int, account common.Address) (bool, []byte) {
func verifySPProof(c *client.Client, proofTrie mpt.Proof, blockNr *big.Int, account []byte) (bool, []byte) {

	wsEndpoint := c.BcWsEndpoint
	client, err := ethclient.Dial(wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	block, _ := client.HeaderByNumber(context.Background(), blockNr)
	stateRoot := block.Root
	// log.Println("stateRootHash: ", stateRoot.Hex)
	// account = common.HexToAddress("0x50D69B935A828113Dd0E4c7Fc721105632014a1d")
	// stateRootHash := common.HexToHash("0x1a204339d2c548efc90e51c67b24d1cd9fcd5ae5b221b490b31cec272e863f7c")
	// log.Println(stateRootHash.Hex)
	// verify the proof against the stateRootHash
	validAccountState, err := mpt.VerifyProof(
		// stateRoot.Bytes(), crypto.Keccak256(account.Bytes()), proofTrie)
		stateRoot.Bytes(), account, proofTrie)

	if err != nil {
		log.Fatal(err)
		return false, nil
	}
	PrintProofToSubmit(proofTrie, account, stateRoot)

	return true, validAccountState
}

func PrintProofToSubmit(proof mpt.Proof, key []byte, rootHash common.Hash) {

	log.Println("root Hash: ", rootHash)
	// 2. Print the serialized proof as a bytes array (proof is []string).
	serializedProofStr := proof.CustomRLPSerializeString()
	fmt.Print("Proof (bytes[] proof): [")
	for i, node := range serializedProofStr {
		if i > 0 {
			fmt.Print(", ")
		}
		// Add "0x" prefix only once for each proof element.
		if len(node) >= 2 && node[:2] == "0x" {
			// If the proof element already has "0x" prefix, use it directly.
			fmt.Printf("\"%s\"", node)
		} else {
			// Otherwise, add "0x" prefix.
			fmt.Printf("\"0x%s\"", node)
		}
	}
	fmt.Println("]")

	keyHex := "0x" + hex.EncodeToString(key)
	// Print the key as bytes[] for Remix.
	fmt.Printf("Keys (bytes[] keys): [\"%s\"]\n", keyHex)
}
