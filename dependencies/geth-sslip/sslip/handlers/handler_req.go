package handlers

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/poc-client/msg/request"
	"github.com/ethereum/go-ethereum/poc-client/utils/cryptoutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/sslip/manager"
	"github.com/ethereum/go-ethereum/sslip/mpt"
	"github.com/ethereum/go-ethereum/sslip/resmsg"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/sha3"
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

	// Log Helper: Print the balance in Ether
	// params := balanceRequest.Params
	// log.Println(params[0])
	// address := common.HexToAddress(params[0].(string))
	// balance, _ := bcClient.BalanceAt(context.Background(), address, nil)
	// balanceInEther := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(math.Pow10(18)))
	// fmt.Printf("Balance of %s: %s ETH\n", address.Hex(), balanceInEther.Text('f', 18)) // formatted to 18 decimal places

	currentBlockHeader, _ := bcClient.HeaderByNumber(context.Background(), nil)
	log.Println("Block number: ", currentBlockHeader.Number.Text(16))

	log.Println("block Header hash: ", currentBlockHeader.Hash().Hex())

	// Log Helper: for block information
	// ComputeBlockHash(currentBlockHeader)
	// PrintBlockInfo(currentBlockHeader)
	// log.Println(BhRlpBytes(currentBlockHeader))

	log.Println("*************************Block Header RLP encoded!*************************")

	channelId := m.GetClientChannelID(clientID)
	channelIdBytes := channelId.Bytes()
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

	// PrintProofToSubmit(storageProof, m.ContractAddress, currentBlockHeader.Root)
	log.Println(res)

	log.Println("*************************Generation & Verification DONE!*************************")

	log.Println()
	log.Printf("----------------Generate Response--------------------")

	txRootHash := currentBlockHeader.Root
	responseSPMsg := resmsg.ResponseSPMsg{
		Type:           "responseSP",
		ChannelId:      channelId,
		Amount:         req.Amount,
		ReqBodyHash:    reqHash,
		SignedReqBody:  req.SignedReqBody,
		CurrentBlockNr: currentBlockHeader.Number,
		ReturnValue:    []byte(validState),
		Proof:          storageProof.CustomRLPSerialize(),
		Address:        crypto.Keccak256(common.HexToAddress(m.ContractAddress).Bytes()),
		Signature:      nil,
		TxRootHash:     txRootHash,
		// BlockNr:            currentBlockHeader.Number,
	}
	PrintProofToSubmit(storageProof, responseSPMsg.Address, currentBlockHeader.Root)
	resBodyHash := responseSPMsg.Keccak256Hash()
	sig := cryptoutil.SignHash(m.PrivateKey, resBodyHash)
	responseSPMsg.Signature = sig

	log.Println("Root Hash: ", txRootHash)
	log.Println("SignedReqBody:", hex.EncodeToString(req.SignedReqBody))
	log.Println("Signature:", hex.EncodeToString(sig))
	log.Println("resHash: ", resBodyHash)
	log.Println("reqHash: ", reqHash)

	fmt.Println("-=-=-=-=-= Now print response message bytes -=-=-=-=-=-=")
	// log.Println(responseSPMsg.RlpBytes())

	// responseSPMsg.Proof[1][3] = '4'
	// responseSPMsg.Proof[0][10] = 'f'
	// resBodyHash = responseSPMsg.Keccak256Hash()
	// sig = cryptoutil.SignHash(m.PrivateKey, resBodyHash)
	// responseSPMsg.Signature = sig
	log.Println(responseSPMsg.RlpBytes())
	fmt.Println("*********************************************************************")

	reqBodyBytesString := req.RequestBodyRlpBytes()
	log.Println(reqBodyBytesString)

	fmt.Println("*********************************************************************")

	log.Println("-=-=-=-=-=-= Now print request payment bytes -=-=-=-=-=-=")
	log.Println("Payment body byte: ", req.PaymentBodyRlpBytes())
	log.Println("*********************************************************************")

	_ = conn.WriteMessage(mt, responseSPMsg.Bytes())

	// 		curl http://localhost:8100 \
	//   -X POST \
	//   -H "Content-Type: application/json" \
	//   -d '{"jsonrpc":"2.0","method":"eth_getProof","params":["0x50D69B935A828113Dd0E4c7Fc721105632014a1d",["0x0"], "0x2EE"],"id":1}'
	return nil
}

func PrintBlockInfo(header *types.Header) {
	if header == nil {
		log.Println("Header is nil.")
		return
	}

	fmt.Println("Block Header Info:")
	fmt.Printf("Parent Hash:            %s\n", header.ParentHash.Hex())
	fmt.Printf("Uncle Hash:             %s\n", header.UncleHash.Hex())
	fmt.Printf("Coinbase (Miner):       %s\n", header.Coinbase.Hex())
	fmt.Printf("State Root:             %s\n", header.Root.Hex())
	fmt.Printf("Transaction Root:       %s\n", header.TxHash.Hex())
	fmt.Printf("Receipt Root:           %s\n", header.ReceiptHash.Hex())
	fmt.Printf("Bloom:                  %x\n", header.Bloom)
	if header.Difficulty != nil {
		fmt.Printf("Difficulty:             %s\n", header.Difficulty.String())
	}
	if header.Number != nil {
		fmt.Printf("Block Number:           %s\n", header.Number.String())
	}
	fmt.Printf("Gas Limit:              %d\n", header.GasLimit)
	fmt.Printf("Gas Used:               %d\n", header.GasUsed)
	fmt.Printf("Timestamp:              %d\n", header.Time)
	fmt.Printf("Extra Data:             %x\n", header.Extra)
	fmt.Printf("Mix Digest:             %s\n", header.MixDigest.Hex())
	fmt.Printf("Nonce:                  %x\n", header.Nonce)

	// Optional Fields
	if header.BaseFee != nil {
		fmt.Printf("Base Fee Per Gas:       %s\n", header.BaseFee.String())
	}
	if header.WithdrawalsHash != nil {
		fmt.Printf("Withdrawals Root:       %s\n", header.WithdrawalsHash.Hex())
	}
	if header.BlobGasUsed != nil {
		fmt.Printf("Blob Gas Used:          %d\n", *header.BlobGasUsed)
	}
	if header.ExcessBlobGas != nil {
		fmt.Printf("Excess Blob Gas:        %d\n", *header.ExcessBlobGas)
	}
	if header.ParentBeaconRoot != nil {
		fmt.Printf("Parent Beacon Root:     %s\n", header.ParentBeaconRoot.Hex())
	}
}

// Function to encode the block header and compute the hash
func ComputeBlockHash(header *types.Header) common.Hash {
	// Create a buffer to hold the RLP-encoded header
	var encodedHeader bytes.Buffer

	// RLP encode the header
	err := rlp.Encode(&encodedHeader, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
		header.MixDigest,
		header.Nonce,
		header.BaseFee, // Optional field
		header.WithdrawalsHash,
		header.BlobGasUsed,
		header.ExcessBlobGas,
		header.ParentBeaconRoot,
	})
	if err != nil {
		log.Fatalf("Failed to RLP encode block header: %v", err)
	}

	// Compute the Keccak256 hash of the encoded header
	hash := sha3.NewLegacyKeccak256()
	hash.Write(encodedHeader.Bytes())

	log.Printf("RLP encoded header: %x\n", encodedHeader.Bytes())
	log.Println("Result")
	log.Println(common.BytesToHash(hash.Sum(nil)))
	return common.BytesToHash(hash.Sum(nil))

}
func generateStorageProof(address string, pos string, blockN *big.Int) *mpt.ProofDB {
	// blockHex := blockN.Text(16)
	rpcClient, _ := rpc.Dial("http://localhost:8000")
	var resultProof json.RawMessage
	// var paramsGenerated = []interface{}{
	// 	address, []string{pos}, blockHex,
	// }

	var paramsGenerated = []interface{}{
		address, []string{pos}, "0x" + blockN.Text(16),
	}

	err := rpcClient.Call(&resultProof, "eth_getProof", paramsGenerated...)
	if err != nil {
		log.Fatalf("Failed to execute request: %v", err)
	}
	// log.Println(resultProof)

	extendedResultProof := append([]byte(`{"jsonrpc":"2.0","id":1,"result":`), resultProof...)
	extendedResultProof = append(extendedResultProof, []byte(`}`)...)

	// log.Println(string(extendedResultProof))
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
	accountAddr := common.HexToAddress(account)

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
