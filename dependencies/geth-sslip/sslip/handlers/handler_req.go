package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/poc-client/msg/request"
	"github.com/ethereum/go-ethereum/poc-client/utils/cryptoutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/sslip/manager"
	"github.com/ethereum/go-ethereum/sslip/mpt"
	"github.com/ethereum/go-ethereum/sslip/resmsg"
	"github.com/gorilla/websocket"
)

func handler_req(clientID string, body string, m *manager.Manager, conn *websocket.Conn, mt int) error {
	client := m.GetClient(clientID)

	req, _ := unmarshalRequest(body)
	// requestByte := req.ReqByte

	log.Println("----------------Request from clientID ", clientID, "---------------------")
	_, reqHash, sigCheckMsg, err := verifyReqSignature(m, clientID, req)

	if err != nil {
		log.Fatal("Unmarshal error: ", err)
		return err
	}

	// Send the result of the signature check to the client
	conn.WriteMessage(mt, sigCheckMsg.Bytes())

	log.Println("*************************Pre check passed!*************************")

	bcClient, _ := client.ConnectToBlockchain()

	var balanceRequest request.JSONRPCRequest
	err = json.Unmarshal(req.ReqByte, &balanceRequest)
	if err != nil {
		log.Fatal("Unmarshal error: ", err)
		return err
	}

	log.Println()
	log.Printf("----------------Process the request--------------------")

	params := balanceRequest.Params
	log.Println(params[0])
	address := common.HexToAddress(params[0].(string))

	balance, err := bcClient.BalanceAt(context.Background(), address, nil)
	balanceInEther := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(math.Pow10(18)))

	// Print the balance in Ether
	fmt.Printf("Balance of %s: %s ETH\n", address.Hex(), balanceInEther.Text('f', 18)) // formatted to 18 decimal places

	// blockHeader, _ := bcClient.HeaderByNumber(context.Background(), nil)
	currentBlockHeader, _ := bcClient.HeaderByNumber(context.Background(), nil)
	log.Println("Block number: ", currentBlockHeader.Number.Text(16))

	channelId := m.GetClientChannelID(clientID)
	// blockHeader, _ = bcClient.HeaderByNumber(context.Background(), big.NewInt(168))

	channelIdBytes := common.HexToHash(channelId).Bytes()
	var position [32]byte
	data := append(channelIdBytes, position[:]...)
	slot := crypto.Keccak256Hash(data)

	log.Println()
	log.Printf("----------------Generate Storage Proof--------------------")

	storageProof := generateStorageProof(m.ContractAddress, slot.Hex(), currentBlockHeader.Number)

	log.Println("*************************Generation end!*************************")

	log.Println()
	log.Printf("----------------Verifying Proof--------------------")

	res, validState := verifyStorageProof(storageProof, m.ContractAddress, currentBlockHeader.Root)
	log.Println(res)

	log.Println("*************************Generation & Verification DONE!*************************")

	log.Println()
	log.Printf("----------------Generate Response--------------------")

	responseSPBody := resmsg.ResponseSPBody{
		SignedReqBody: req.SignedReqBody,
		Proof:         storageProof.CustomRLPSerialize(),
		// Address:       common.HexToAddress(m.ContractAddress),
		// BlockNr: currentBlockHeader.Number,
	}

	sig := cryptoutil.SignHash(m.PrivateKey, responseSPBody.Keccak256Hash())
	responseSPMsg := resmsg.ResponseSPMsg{
		Type:               "responseSP",
		ChannelId:          common.HexToHash(channelId),
		Amount:             req.Amount,
		ReqBodyHash:        reqHash,
		SignedReqBody:      req.SignedReqBody,
		CurrentBlockHeight: currentBlockHeader.Number,
		ReturnValue:        []byte(validState),
		Proof:              storageProof.CustomRLPSerialize(),
		Address:            common.HexToAddress(m.ContractAddress),
		// BlockNr:            currentBlockHeader.Number,
		Signature: sig,
	}

	// log.Println(responseSPMsg)
	// client.Send(responseSPMsg.Bytes())

	log.Println("SignedReqBody:", hex.EncodeToString(req.SignedReqBody))
	log.Println("Signature:", hex.EncodeToString(sig))
	log.Println("resHash: ", responseSPBody.Keccak256Hash())
	log.Println("reqHash: ", reqHash)

	fmt.Println("-=-=-=-=-= Now print response message bytes -=-=-=-=-=-=")
	log.Println(responseSPMsg.RlpBytes())

	fmt.Println("*********************************************************************")

	reqBody := request.ReqBody{
		// ChannelID:      req.ChannelID,
		Amount:         req.Amount,
		LocalBlockHash: req.LocalBlockHash,
		ReqByte:        req.ReqByte,
	}
	reqBodyBytesString := reqBody.RlpBytes()
	log.Println(reqBodyBytesString)

	fmt.Println("*********************************************************************")

	_ = conn.WriteMessage(mt, responseSPMsg.Bytes())

	// storageKeys := []string{"0x0"} // Example storage key
	// blockNumber := big.NewInt(blockHeader.Number.Int64())

	// proofResponse := bcClient.Client().CallContext(context.Background(), nil, "eth_getProof", address, []string{"0x0"}, "0x2EE")
	// log.Println(proofResponse)
	// GetProof(context.Background(), address, blockNumber, storageKeys)

	// 		curl http://localhost:8100 \
	//   -X POST \
	//   -H "Content-Type: application/json" \
	//   -d '{"jsonrpc":"2.0","method":"eth_getProof","params":["0x50D69B935A828113Dd0E4c7Fc721105632014a1d",["0x0"], "0x2EE"],"id":1}'
	return nil
}

func generateStorageProof(address string, pos string, blockN *big.Int) *mpt.ProofDB {
	// blockHex := blockN.Text(16)
	rpcClient, _ := rpc.Dial("http://localhost:8000")
	var resultProof json.RawMessage
	// var paramsGenerated = []interface{}{
	// 	address, []string{pos}, blockHex,
	// }
	log.Println(address)
	log.Println(pos)
	log.Println(blockN.Text(16))
	var paramsGenerated = []interface{}{
		address, []string{pos}, "0x" + blockN.Text(16),
	}

	err := rpcClient.Call(&resultProof, "eth_getProof", paramsGenerated...)
	if err != nil {
		log.Fatalf("Failed to execute request: %v", err)
	}
	log.Println(resultProof)

	extendedResultProof := append([]byte(`{"jsonrpc":"2.0","id":1,"result":`), resultProof...)
	extendedResultProof = append(extendedResultProof, []byte(`}`)...)

	log.Println(string(extendedResultProof))
	// load into the struct
	var response mpt.EthGetProofResponse
	err = json.Unmarshal(extendedResultProof, &response)
	if err != nil {
		log.Fatal(err)
	}
	result := response.Result

	// account := common.HexToAddress("0xcca577ee56d30a444c73f8fc8d5ce34ed1c7da8b")
	account := common.HexToAddress(address)
	fmt.Println(fmt.Sprintf("decoded account state data from untrusted source for address %x: balance is %x, nonce is %x, codeHash: %x, storageHash: %x",
		account, result.Balance, result.Nonce, result.CodeHash, result.StorageHash))

	// get the state root hash from etherscan: https://etherscan.io/block/11045195
	// stateRootHash := common.HexToHash("0x8c571da4c95e212e508c98a50c2640214d23f66e9a591523df6140fd8d113f29")

	// create a proof trie, and add each node from the account proof
	proofTrie := mpt.NewProofDB()
	for _, node := range result.AccountProof {
		proofTrie.Put(crypto.Keccak256(node), node)
	}

	return proofTrie

}

func verifyStorageProof(proofTrie *mpt.ProofDB, account string, stateRootHash common.Hash) (bool, []byte) {
	log.Println("stateRootHash: ", stateRootHash)
	accountAddr := common.HexToAddress(account)
	// stateRootHash = common.HexToHash("0x7b127327de0842612844cf2c450a42270de0fefbf86ef75dc5d82777df0be300")

	// verify the proof against the stateRootHash
	validAccountState, err := mpt.VerifyProof(
		stateRootHash.Bytes(), crypto.Keccak256(accountAddr.Bytes()), proofTrie)

	if err != nil {
		log.Fatal(err)
		return false, nil
	}

	return true, validAccountState

	// double check the account state is identical with the account state in the result.
	// accountState, err := rlp.EncodeToBytes([]interface{}{
	// 	result.Nonce,
	// 	result.Balance.ToInt(),
	// 	result.StorageHash,
	// 	result.CodeHash,
	// })
	// log.Println(bytes.Equal(validAccountState, accountState))

	// fmt.Sprintf("%x!=%x", validAccountState, accountState)
}
