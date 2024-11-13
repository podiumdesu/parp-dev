package handlers

import (
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
)

func HandleResponseSP(msg []byte, client *client.Client) error {
	log.Println(string(msg))
	var resMsg resmsg.ResponseSPMsg
	err := json.Unmarshal(msg, &resMsg)
	if err != nil {
		log.Println("Unmarshal error: ", err)
		return err
	}
	res := cryptoutil.Verify(crypto.FromECDSAPub(client.ServerPublicKey), resMsg.BodyHashBytes(), resMsg.Signature)
	log.Println("Verify Response signature:", res)
	proof, _ := mpt.DeserializeProof(resMsg.Proof)
	// txProofRLP, _ := mpt.VerifyProof(txRootHash[:], key, proof)

	// proof, err := mpt.DeserializeProofNodes(resMsg.Proof)
	if err != nil {
		log.Println("Error deserializing proof: ", err)
	}
	result, validState := verifySPProof(client, proof, resMsg.BlockNr, resMsg.Address)
	log.Println(result, validState)
	return nil
}

func verifySPProof(c *client.Client, proofTrie mpt.Proof, blockNr *big.Int, account common.Address) (bool, []byte) {
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
		stateRoot.Bytes(), crypto.Keccak256(account.Bytes()), proofTrie)

	if err != nil {
		log.Fatal(err)
		return false, nil
	}

	return true, validAccountState
}
