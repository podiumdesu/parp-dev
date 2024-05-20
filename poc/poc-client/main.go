package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"poc-server/mpt"
	"poc-server/resmsg"
	"sync"
	"time"

	"poc-client/benchmarking"
	"poc-client/client"
	"poc-client/connection"
	"poc-client/hub"

	"poc-client/hub/wsClient"
	"poc-client/msg/request"
	"poc-client/protocol"
	"poc-client/utils/cryptoutil"

	// "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	// "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	// "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
)

var wg sync.WaitGroup

type Config struct {
	PrivateKeyFilePath string `json:"privateKeyFilePath"`
	BcWsEndpoint       string `json:"bcWsEndpoint"`
	BcRpcEndpoint      string `json:"bcRpcEndpoint"`
	ServerEndpoint     string `json:"serverEndpoint"`
	ContractAddress    string `json:"contractAddress"`
}

var startTime time.Time
var startTimeQuerySslip time.Time
var startTimeVerifySslip time.Time
var durationQueryTimes []time.Duration
var durationQueryVerifyTimes []time.Duration

type Result struct {
	Success       bool
	ProofDuration time.Duration
	TotalDuration time.Duration
}

var results []Result        // This will hold the results from each response
var resultsMutex sync.Mutex // Mutex to protect access to the results slice

func main() {
	// The hub is responsible for managing all the websocket connections
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	log.Println("Loading config....")

	log.Println("\n-----------------Generate a new client-------------------")
	client := (&client.Client{}).Init(config.PrivateKeyFilePath, config.BcWsEndpoint, config.BcRpcEndpoint, config.ServerEndpoint)

	fmt.Println("---------------------------------------------------------\n")
	hub := hub.NewHub()
	go hub.Run()

	// Setup connection with the PoC server

	server_wsConn := connection.ConnectToServer(client)

	server := &wsClient.Client{Conn: server_wsConn, Send: make(chan []byte, 10000), Receive: make(chan []byte, 10000)}
	hub.Set_fn <- server

	// TOREMOVE: Benchmarking
	// responses := make(chan resmsg.ResponseMsg, 1) // Buffer size can be adjusted

	go server.WritePump()
	go server.ReadPump()

	go func() {
		for msg := range server.Receive {
			var serverMsg resmsg.ServerMsg
			err := json.Unmarshal(msg, &serverMsg)

			log.Println(serverMsg.Type)
			if err != nil {
				log.Println("Unmarshal error: ", err)
				break
			}

			switch serverMsg.Type {
			case "info":
				var infoMsg resmsg.ServerMsg
				err := json.Unmarshal(msg, &infoMsg)
				if err != nil {
					log.Println("Unmarshal error: ", err)
					break
				}
				log.Println(string(infoMsg.Info))
			case "info-hex":
				log.Println(hex.EncodeToString(serverMsg.Info))
			case "HANDSHAKE-CONFIRMED":
				var hsMsg resmsg.HandshakeMsg
				err := json.Unmarshal(msg, &hsMsg)
				if err != nil {
					log.Println("Unmarshal error: ", err)
					break
				}
				serverPubKeyByte := hsMsg.ServerPublicKeyB
				serverPubKeyECDSA, _ := crypto.UnmarshalPubkey(serverPubKeyByte)
				client.ServerPublicKey = serverPubKeyECDSA
				serverPubKeyAddress := crypto.PubkeyToAddress(*serverPubKeyECDSA)
				log.Println("Server public key address has been set: ", serverPubKeyAddress.Hex())
			case "SigCheck":
				log.Println(msg)
			case "response":
				log.Println("response")
				// benchmarking.HandleResponses(responses, client)
				// log.Println(string(msg))
				var bresMsg resmsg.ResponseMsg

				// log.Println("Size of the Tx response: %d bytes", len(msg))
				err := json.Unmarshal(msg, &bresMsg)
				if err != nil {
					log.Println("Unmarshal error: ", err)
					break
				}
				wg.Add(1)
				go func(resMsg resmsg.ResponseMsg) {
					defer wg.Done()
					verifyTimer := time.Now()

					proof, err := mpt.DeserializeProof(resMsg.Proof)
					if err != nil {
						log.Printf("Error deserializing proof: %v\n", err)
						return
					}
					proofTimer := time.Now()
					proofRes := verifyProof(client, resMsg.TxHash, proof, resMsg.CurrentBlockHeight, uint32(resMsg.TxIdx))
					proofDuration := time.Since(proofTimer)

					resMsgBodyHash := resMsg.BodyHashBytes()
					_ = cryptoutil.Verify(crypto.FromECDSAPub(client.ServerPublicKey), resMsgBodyHash, resMsg.Signature)

					totalDuration := time.Since(verifyTimer)
					restVeriDuration := totalDuration - proofDuration
					if proofRes {
						log.Printf("Duration for proof verification: %s\n", proofDuration)
						log.Printf("Duration for the rest verification: %s\n", restVeriDuration)
						log.Printf("Duration for request verification: %s\n", totalDuration)
					}
					// Protect the slice with a mutex
					resultsMutex.Lock()
					results = append(results, Result{
						Success:       proofRes,
						ProofDuration: proofDuration,
						TotalDuration: restVeriDuration,
					})
					resultsMutex.Unlock()
					log.Println(results)
				}(bresMsg)
				wg.Wait()

				log.Println("Done")
				calculateResults(results, resultsMutex)

				break
				// durationRes := time.Since(startTime)
				verifyTimer := time.Now()
				// log.Printf("Duration for transaction: %s\n", durationRes)

				// log.Println(string(msg))
				var resMsg resmsg.ResponseMsg

				// log.Println("Size of the Tx response: %d bytes", len(msg))
				err = json.Unmarshal(msg, &resMsg)
				if err != nil {
					log.Println("Unmarshal error: ", err)
					break
				}

				proofTimer := time.Now()
				proof, err := mpt.DeserializeProof(resMsg.Proof)
				if err != nil {
					log.Println("Error deserializing proof: ", err)
				}
				proofRes := verifyProof(client, resMsg.TxHash, proof, resMsg.CurrentBlockHeight, uint32(resMsg.TxIdx))
				durationProof := time.Since(proofTimer)

				// log.Println("size of proof for transactions: ", len(resMsg.Proof))
				resMsgBodyHash := resMsg.BodyHashBytes()
				_ = cryptoutil.Verify(crypto.FromECDSAPub(client.ServerPublicKey), resMsgBodyHash, resMsg.Signature)
				// log.Println("Verify Response signature:", res)

				durationVeri := time.Since(verifyTimer)
				durationVeriRest := durationVeri - durationProof

				if proofRes {
					log.Printf("Duration for proof verification: %s\n", durationProof)
					log.Printf("Duration for the rest verification: %s\n", durationVeri)
					log.Printf("Duration for request verification: %s\n", durationVeriRest)
				}

				// return durationProof, durationVeri, durationVeriRest
			case "responseSP":
				durationGetQueryBackSSLIP := time.Since(startTimeQuerySslip)
				durationQueryTimes = append(durationQueryTimes, durationGetQueryBackSSLIP)
				startTimeVerifySslip = time.Now()
				log.Println(string(msg))
				var resMsg resmsg.ResponseSPMsg
				err := json.Unmarshal(msg, &resMsg)

				log.Println("Size of the Request response: %d bytes", len(msg))

				if err != nil {
					log.Println("Unmarshal error: ", err)
					break
				}
				res := cryptoutil.Verify(crypto.FromECDSAPub(client.ServerPublicKey), resMsg.BodyHashBytes(), resMsg.Signature)
				log.Println("Verify Response signature:", res)

				proof, err := mpt.DeserializeProof(resMsg.Proof)
				if err != nil {
					log.Println("Error deserializing proof: ", err)
				}
				result, validState := verifySPProof(client, proof, resMsg.BlockNr, resMsg.Address)
				log.Println(result, validState)
				if result {
					durationQueryVeriTime := time.Since(startTimeVerifySslip)
					durationQueryVerifyTimes = append(durationQueryVerifyTimes, durationQueryVeriTime)
				}
			default:
				log.Println("Unrecognized message type: ", serverMsg.Type)
			}
		}

	}()

	wg.Add(1)

	go func() {
		log.Println("\n------------------Handshake-------------------")

		defer wg.Done()
		msg := client.InitHandshakeMsg(config.ContractAddress, 10, big.NewInt(100000), big.NewInt(200))
		b := append([]byte("HANDSHAKE:"), msg.Bytes()...)

		// b := []byte("HANDSHAKE:" + string(client.PubKeyBytes()))
		log.Println("Sending: ", b)

		hub.Send_fn <- b //[]byte("SIG:qwerrtqreqwrqwerqwrwerqwrewqtqwetqwrewqrqwerwqewqrqwer")
		fmt.Println("----------------------------------------------\n")
		hub.Send_fn <- []byte("FE: Connected to server")
	}()

	wg.Wait()

	// go func() {
	// 	log.Println("\n------------------Send OpenChan Tx request--------------------")
	// 	OpenChanTx := sendOpenChanTxs(client, common.HexToAddress(config.ContractAddress))
	// 	// startTime = time.Now()

	// 	select {
	// 	case hub.Send_fn <- OpenChanTx:
	// 		log.Println("Request message sent successfully")
	// 	default:
	// 		log.Println("Failed to send request message: channel is full or closed")
	// 	}
	// 	fmt.Println("------------------------------------------------------\n")
	// }()

	go func() {
		benchmarking.Greeting()

		const requestNum = 5
		wsEndpoint := client.BcWsEndpoint
		bcClient, err := ethclient.Dial(wsEndpoint)

		privateKey := client.PrivateKey

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

		nonce, err := bcClient.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatal(err)
		}

		var totalReqGenTime time.Duration
		for i := 0; i < requestNum; i++ {
			// go func(i int) {
			requestGenTimer := time.Now()
			OpenChanTx := sendOpenChanTxs(client, common.HexToAddress(config.ContractAddress), nonce+uint64(i))

			duration := time.Since(requestGenTimer)
			totalReqGenTime += duration
			select {
			case hub.Send_fn <- OpenChanTx:
				log.Printf("Request %d sent successfully\n", i)
			default:
				log.Printf("Failed to send request %d: channel is full or closed\n", i)
			}
			// }(i)
		}
		log.Println("Average request generation time: ", totalReqGenTime/time.Duration(requestNum))

		log.Println("All requests have been sent")

		// // Benchmarking for Geth nodes
		// log.Println("\n------------------Send OpenChan Tx request to geth--------------------")
		// benchmarking.GethSyncTx(client, common.HexToAddress(config.ContractAddress))
		// fmt.Println("------------------------------------------------------\n")

		// log.Println("\n------------------Send OpenChan Tx request to geth async--------------------")
		// benchmarking.GethAsyncTx(client, common.HexToAddress(config.ContractAddress))
		// fmt.Println("------------------------------------------------------\n")

		// for i := 0; i < 10; i++ {
		// 	log.Println("\n------------------Send BalanceChecking request--------------------")
		// 	balanceCheckingReq := sendRequests(client, client.Amount+20)

		// 	startTimeQuerySslip = time.Now()
		// 	select {
		// 	case hub.Send_fn <- balanceCheckingReq:
		// 		log.Println("Request message sent successfully")
		// 	default:
		// 		log.Println("Failed to send request message: channel is full or closed")
		// 	}
		// 	fmt.Println("------------------------------------------------------\n")
		// 	time.Sleep(400 * time.Millisecond) // Sleep for 100 milliseconds
		// }

		// var totalDurationResponse time.Duration
		// var totalDurationVerification time.Duration
		// for _, d := range durationQueryTimes {
		// 	totalDurationResponse += d
		// }
		// for _, d := range durationQueryVerifyTimes {
		// 	totalDurationVerification += d
		// }

		// averageDurationResponse := totalDurationResponse / time.Duration(len(durationQueryTimes))
		// averageDurationVerify := totalDurationVerification / time.Duration(len(durationQueryVerifyTimes))

		// log.Printf("Average duration for %d requests: %s", len(durationQueryTimes), averageDurationResponse)
		// log.Printf("Average duration for %d requests got verified: %s", len(durationQueryVerifyTimes), averageDurationVerify)

		// Benchmarking: geth requests

		// benchmarking.GethSyncQuery(client)  //20ms

	}()

	go func() {
		// benchmarking.ProofGenerationAndVerification(client)
	}()
	select {}

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

func loadConfig(filename string) (*Config, error) {
	// Read configuration file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse configuration from JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
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
	// log.Println("txProofRLP: ", txProofRLP)
	// log.Println("txRLP: ", txRLP)
	// log.Println("proof: ", proof.Serialize())

	return bytes.Equal(txRLP, txProofRLP)
}

func VerifyFraudProof(resMsg resmsg.ResponseMsg) bool {

	return true

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

func calculateResults(results []Result, resultsMutex sync.Mutex) {
	log.Println("Calculate results")
	var totalProofDuration, totalDuration time.Duration
	var successCount int

	resultsMutex.Lock()
	for _, result := range results {
		if result.Success {
			successCount++
			totalProofDuration += result.ProofDuration
			totalDuration += result.TotalDuration
		}
	}
	resultsMutex.Unlock()

	averageProofDuration := totalProofDuration / time.Duration(len(results))
	averageTotalDuration := totalDuration / time.Duration(len(results))
	successRate := float64(successCount) / float64(len(results)) * 100

	log.Printf("Average Proof Verification Duration: %s\n", averageProofDuration)
	log.Printf("Average Total Verification Duration: %s\n", averageTotalDuration)
	log.Printf("Success Rate: %.2f%%\n", successRate)
}

func sendOpenChanTxs(client *client.Client, contractAddress common.Address, idx uint64) []byte {
	// in @openChan.go as well
	wsEndpoint := client.BcWsEndpoint
	bcClient, err := ethclient.Dial(wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	blockHeader, _ := bcClient.HeaderByNumber(context.Background(), nil)
	log.Println("Block Hash: ", blockHeader.Hash().Hex())

	// log.Println("\n------------------Send OpenChan request--------------------")
	fnAddr := common.HexToAddress("0xA2131E7503F7Dd11ff5dAAC09fa7c301e7Fe0f30")
	deposit := big.NewInt(200000)

	signedTx := protocol.OpenChanTxToGeth(client, fnAddr, deposit, contractAddress, uint64(idx))

	ts := types.Transactions{signedTx}
	rawTxBytes, _ := rlp.EncodeToBytes(ts[0])

	msg := protocol.GenerateRequest(client, 20, client.Amount+100, rawTxBytes, blockHeader.Hash())

	client.Amount += 100
	msgWType := append([]byte("TX:"), msg...)

	log.Println("Sending: ", string(msgWType))

	return msgWType
}
