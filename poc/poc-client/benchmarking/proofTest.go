package benchmarking

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"poc-client/client"
	"poc-client/msg/request"
	"poc-client/protocol"
	"poc-client/utils/cryptoutil"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"poc-server/mpt"
	"poc-server/resmsg"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func ProofGenerationAndVerification(client *client.Client) {
	wsEndpoint := client.BcWsEndpoint
	bcClient, _ := ethclient.Dial(wsEndpoint)

	numRequests := 100

	// average time for generating a msg
	msg := sendRequestAverage(client, numRequests)
	return

	// pick the block (with a lot of transactions)
	// nr := 42 // 1
	// nr := 47 // 10
	// nr := 49 // 50
	// nr := 46 / 100
	nr := 512 // 200
	// nr := 52 //300
	// nr := 54 //400

	blockNr := big.NewInt(int64(nr))
	block, _ := bcClient.BlockByNumber(context.Background(), blockNr)
	txs := block.Transactions()

	// average time for generating a proof
	proofArray := generateProofAverage(numRequests, client, blockNr)

	// average time for verifying a proof
	verifyProofAverage(numRequests, client, blockNr, proofArray)

	// Request Size and Response size
	for i := 0; i < numRequests; i++ {
		ResponseAndRequestSize(client, msg, proofArray[i], txs[i].Hash(), uint32(i), blockNr)
	}
}

func generateTxRequest() {

}

func ResponseAndRequestSize(client *client.Client, msg []byte, proof mpt.Proof, txHash common.Hash, idx uint32, blockNr *big.Int) {

	log.Println("Size of Request: ", len(msg))
	wsEndpoint := client.BcWsEndpoint
	bcClient, _ := ethclient.Dial(wsEndpoint)

	parts := strings.SplitN(string(msg), ":", 2)
	if len(parts) < 2 {
		log.Println("Invalid message format")
		return
	}

	body := parts[1]
	var req request.RequestMsg
	err := json.Unmarshal([]byte(body), &req)

	if err != nil {
		log.Fatal("Unmarshal error: ", err)
	}
	// requestByte := req.ReqByte
	// log.Println("Size of Request: ", len(requestByte))

	txReceipt, err := bcClient.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	channelId := common.HexToHash("0xae3bd68e96cf44fbda03d80e8745125d007c61ff72d2937ad2266e4c853a6b53")

	responseBody := resmsg.ResponseBody{
		SignedReqBody: req.SignedReqBody,
		Proof:         proof.CustomSerialize(),
		TxHash:        txHash,
		TxIdx:         uint32(idx),
	}

	ttLength := 0
	for _, sublist := range proof.CustomSerialize() {
		ttLength += len(sublist)
	}
	log.Println("size of proof: ", ttLength)

	serverPrivKey, _ := crypto.HexToECDSA("bcd5c542c981dbb7cee1f3352fcee082581b4a323bf5cbff105aa84fa718f690")

	sig := cryptoutil.Sign(serverPrivKey, responseBody.HashBytes())
	responseMsg := resmsg.ResponseMsg{
		Type:               "response",
		ChannelId:          channelId,
		Amount:             req.Amount,
		SignedReqBody:      req.SignedReqBody,
		CurrentBlockHeight: blockNr,
		ReturnValue:        txReceipt.Bloom.Bytes(),
		Proof:              proof.CustomSerialize(),
		TxHash:             txHash,
		TxIdx:              uint32(idx),
		Signature:          sig,
	}
	log.Println("size of Response after json.Unmarshall: ", len(responseMsg.Bytes()))

	v := reflect.ValueOf(responseMsg)
	t := v.Type()
	responseSize := 0
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i).Type

		var length int
		switch field.Kind() {
		case reflect.String:
			length = len(field.String())
		case reflect.Slice:
			length = field.Len()
			if fieldType.Elem().Kind() == reflect.Slice {
				// If the element type is a slice (e.g., [][]byte), calculate the total length of inner slices
				length = 0
				for j := 0; j < field.Len(); j++ {
					length += field.Index(j).Len()
				}
			}
		case reflect.Ptr:
			if fieldType.Elem().Kind() == reflect.Struct && fieldType.Elem() == reflect.TypeOf(big.Int{}) {
				// Calculate the length of a *big.Int
				length = len(field.Interface().(*big.Int).Bytes())
			}
		case reflect.Array:
			length = field.Len()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			length = int(unsafe.Sizeof(field.Interface()))
		default:
			length = int(unsafe.Sizeof(field.Interface()))
		}

		fmt.Printf("%s: %d bytes\n", t.Field(i).Name, length)
		responseSize += length
	}
	log.Println("size of Response: ", responseSize)
}

func verifyProofAverage(numRequests int, client *client.Client, blockNr *big.Int, proofArray []mpt.Proof) {

	wsEndpoint := client.BcWsEndpoint
	bcClient, _ := ethclient.Dial(wsEndpoint)
	var totalDuration time.Duration

	block, _ := bcClient.BlockByNumber(context.Background(), blockNr)

	txs := block.Transactions()
	// txs := TransactionsJSON()
	log.Println("in total proof number: ", numRequests)
	for i := 0; i < numRequests; i++ {
		log.Println("Verify Proof Average Transaction ", i)
		tx := txs[i]
		start := time.Now()
		verifyProof(client, tx.Hash(), proofArray[i], blockNr, uint32(i))
		duration := time.Since(start)
		totalDuration += duration

	}
	averageDuration := totalDuration / time.Duration(numRequests)
	fmt.Printf("Average proof verification time: %s\n", averageDuration)
	// return msg
}

func generateProofAverage(numRequests int, c *client.Client, blockNr *big.Int) []mpt.Proof {
	wsEndpoint := c.BcWsEndpoint
	bcClient, _ := ethclient.Dial(wsEndpoint)

	block, _ := bcClient.BlockByNumber(context.Background(), blockNr)

	txs := block.Transactions()
	// txs := TransactionsJSON()
	log.Println("in total transactions number: ", numRequests)

	var totalDuration time.Duration

	var proofArray []mpt.Proof
	for i := 0; i < numRequests; i++ {
		log.Println("Generate Proof Average Transaction ", i)
		tx := txs[i]
		start := time.Now()
		proof, _ := generateProof(block, tx.Hash())
		log.Println(tx.Hash())
		duration := time.Since(start)
		totalDuration += duration

		proofArray = append(proofArray, proof)

		ttLength := 0
		for _, sublist := range proof.CustomSerialize() {
			ttLength += len(sublist)
		}
		log.Println("size of proof: ", ttLength)
	}

	averageDuration := totalDuration / time.Duration(numRequests)
	fmt.Printf("Average proof generation: %s\n", averageDuration)

	return proofArray
}

func sendRequestAverage(client *client.Client, numRequests int) []byte {
	var totalDuration time.Duration

	var msg []byte
	for i := 0; i < numRequests; i++ {
		log.Println("Transaction ", i)
		start := time.Now()
		msg = sendRequests(client, 300)
		duration := time.Since(start)
		totalDuration += duration
	}

	// Average time to verify the message

	var sigVeriDuration time.Duration

	clientPrivKey, _ := crypto.HexToECDSA("f39985d76a9bf831c3a3fe19cfbe7d038ad25b5de2ceffd4a1cf15191808b396")
	publicKey := clientPrivKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	for i := 0; i < numRequests; i++ {
		veriTimer := time.Now()
		parts := strings.SplitN(string(msg), ":", 2)
		if len(parts) < 2 {
			log.Println("Invalid message format")
			return nil
		}

		header := parts[0]
		body := parts[1]
		switch header {
		case "REQ":
			var req request.RequestMsg
			err := json.Unmarshal([]byte(body), &req)
			if err != nil {
				log.Fatal("Unmarshal error: ", err)
				break
			}
			// jsonReq, _ := json.Marshal(req)

			// requestByte := req.ReqByte

			// serverPrivKey, _ := crypto.HexToECDSA("bcd5c542c981dbb7cee1f3352fcee082581b4a323bf5cbff105aa84fa718f690")
			// placeholder for amount checks
			var tmp bool
			if req.Amount < 0 {
				log.Println("Invalid amount")
				tmp = false
			}
			tmp = true

			//
			pubKey := crypto.FromECDSAPub(publicKeyECDSA)
			// verify the signature
			sigFlag := VerifyRequestWithSig(req, pubKey)

			res := sigFlag && tmp
			if res {
				sigVeriDuration += time.Since(veriTimer)
			}
		}

	}

	averageDuration := totalDuration / time.Duration(numRequests)
	fmt.Printf("Average request time: %s\n", averageDuration)

	averageVeriDuration := sigVeriDuration / time.Duration(numRequests)
	fmt.Printf("Average request verification time: %s\n", averageVeriDuration)
	return msg
}

func VerifyRequestWithSig(req request.RequestMsg, pubKB []byte) bool {
	// Have to verify both signatures

	requestBody := request.ReqBody{
		ChannelID:      req.ChannelID,
		Amount:         req.Amount,
		ReqByte:        req.ReqByte,
		LocalBlockHash: req.LocalBlockHash,
	}

	paymentBody := request.PaymentBody{
		ChannelID: req.ChannelID,
		Amount:    req.Amount,
	}

	reqBSig := req.SignedReqBody
	payBSig := req.SignedPaymentBody

	if reqBSig[64] == 27 || reqBSig[64] == 28 {
		reqBSig[64] -= 27
	}

	if payBSig[64] == 27 || payBSig[64] == 28 {
		payBSig[64] -= 27
	}

	rbFlag := cryptoutil.Verify(pubKB, requestBody.PreHashByte(), req.SignedReqBody)
	pbFlag := cryptoutil.Verify(pubKB, paymentBody.PreHashByte(), req.SignedPaymentBody)
	log.Println("rbFlag: ", rbFlag, " pbFlag: ", pbFlag)
	return rbFlag && pbFlag
}

func sendRequests(client *client.Client, amount uint) []byte {
	wsEndpoint := client.BcWsEndpoint
	bcClient, err := ethclient.Dial(wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	blockHeader, _ := bcClient.HeaderByNumber(context.Background(), nil)

	request := request.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBalance",
		Params:  []interface{}{"0xA2131E7503F7Dd11ff5dAAC09fa7c301e7Fe0f30", "latest"},
		ID:      1,
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
	}
	msg := protocol.GenerateRequest(client, 20, amount, jsonRequest, blockHeader.Hash())
	client.Amount = amount
	msgWType := append([]byte("REQ:"), msg...)
	return msgWType
}

func TransactionsJSON() []*types.Transaction {
	jsonFile, err := os.Open("/Users/weihong/github/PoC/ethereum-pos-testnet/poc/poc-client/benchmarking/transactions.json")
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var txs []*types.Transaction
	json.Unmarshal(byteValue, &txs)
	log.Println(txs)
	return txs
}

func generateProof(block *types.Block, txHash common.Hash) (mpt.Proof, int) {
	txs := block.Transactions()
	// txs := TransactionsJSON()
	idx := -1
	for index, tx := range txs {
		if tx.Hash() == txHash {
			idx = index
		}
	}
	if idx < 0 {
		return nil, -1
	}

	mptTrie := mpt.NewTrie()
	for i, tx := range txs {
		key, _ := rlp.EncodeToBytes(uint(i))

		transaction := fromEthTransaction(tx)

		rlp, _ := transaction.GetRLP()

		mptTrie.Put(key, rlp)
	}

	// generate the proof and verify it
	key, _ := rlp.EncodeToBytes(uint(idx))
	proof, _ := mptTrie.Prove(key)
	// proofSize := len(proof.Serialize()[0])
	// log.Println("proofSize: ", proofSize)

	// fmt.Printf("proof: %x, found: %v\n", proof, found)
	return proof, idx
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
	// txRootHash, _ := hex.DecodeString("bb345e208bda953c908027a45aa443d6cab6b8d2fd64e83ec52f1008ddeafa58")
	tx, _, _ := bcClient.TransactionByHash(context.Background(), txHash)
	// txs := TransactionsJSON()
	// tx := txs[0]
	txRLP, _ := rlp.EncodeToBytes(tx)
	key, _ := rlp.EncodeToBytes(uint32(idx))
	txProofRLP, _ := mpt.VerifyProof(txRootHash[:], key, proof)
	// log.Println("txProofRLP: ", txProofRLP)
	// log.Println("txRLP: ", txRLP)
	// log.Println("proof: ", proof.Serialize())

	log.Println(bytes.Equal(txRLP, txProofRLP))
	return bytes.Equal(txRLP, txProofRLP)
}

func fromEthTransaction(t *types.Transaction) *mpt.Transaction {
	v, r, s := t.RawSignatureValues()
	return &mpt.Transaction{
		AccountNonce: t.Nonce(),
		Price:        t.GasPrice(),
		GasLimit:     t.Gas(),
		Recipient:    t.To(),
		Amount:       t.Value(),
		Payload:      t.Data(),
		V:            v,
		R:            r,
		S:            s,
	}
}
