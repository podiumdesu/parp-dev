package main

import (
	"bytes"
	"context"
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
var durationQueryProofVerifyTimes []time.Duration
var durationQueryRestVerifyTimes []time.Duration

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
				durationRes := time.Since(startTime)
				verifyTimer := time.Now()
				log.Printf("Duration for transaction: %s\n", durationRes)

				// log.Println(string(msg))
				var resMsg resmsg.ResponseMsg

				// log.Println("Size of the Tx response: %d bytes", len(msg))
				err := json.Unmarshal(msg, &resMsg)
				if err != nil {
					log.Println("Unmarshal error: ", err)
					break
				}

				resMsgBodyHash := resMsg.BodyHashBytes()
				res := cryptoutil.Verify(crypto.FromECDSAPub(client.ServerPublicKey), resMsgBodyHash, resMsg.Signature)
				// log.Println("Verify Response signature:", res)
				log.Println(resMsg.Proof)
				ttLength := 0
				for _, sublist := range resMsg.Proof {
					ttLength += len(sublist)
				}
				log.Println("size of proof for transactions: ", ttLength)
				proof, err := mpt.DeserializeProof(resMsg.Proof)

				// log.Println("size of deserializedProof: ", len(proof))
				if err != nil {
					log.Println("Error deserializing proof: ", err)
				}
				res = verifyProof(client, resMsg.TxHash, proof, resMsg.CurrentBlockHeight, uint32(resMsg.TxIdx))

				durationVeriTotal := time.Since(startTime)
				durationVeri := time.Since(verifyTimer)
				if res {
					log.Printf("Duration for transaction: %s\n", durationVeriTotal)
					log.Printf("Duration for transaction verification: %s\n", durationVeri)
				}

				log.Println("Proof Verification: ", res)
			case "responseSP":
				durationGetQueryBackSSLIP := time.Since(startTimeQuerySslip)
				durationQueryTimes = append(durationQueryTimes, durationGetQueryBackSSLIP)

				startTimeVerifySslip = time.Now()
				// log.Println(string(msg))
				var resMsg resmsg.ResponseSPMsg
				err := json.Unmarshal(msg, &resMsg)

				// log.Println("Size of the Request response: %d bytes", len(msg))

				if err != nil {
					log.Println("Unmarshal error: ", err)
					break
				}
				_ = cryptoutil.Verify(crypto.FromECDSAPub(client.ServerPublicKey), resMsg.BodyHashBytes(), resMsg.Signature)
				// log.Println("Verify Response signature:", res)

				proofVerifyTimer := time.Now()
				proof, err := mpt.DeserializeProof(resMsg.Proof)
				if err != nil {
					log.Println("Error deserializing proof: ", err)
				}
				result, validState := verifySPProof(client, proof, resMsg.BlockNr, resMsg.Address)
				proofVerifyDuration := time.Since(proofVerifyTimer)


				log.Println(result, validState)
				if result {
					durationQueryVeriTime := time.Since(startTimeVerifySslip)
					durationQueryVerifyTimes = append(durationQueryVerifyTimes, durationQueryVeriTime)
					durationQueryRestTime := durationQueryVeriTime - proofVerifyDuration
					durationQueryRestVerifyTimes = append(durationQueryRestVerifyTimes, durationQueryRestTime)
					durationQueryProofVerifyTimes = append(durationQueryProofVerifyTimes, proofVerifyDuration)
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

	go func() {
		benchmarking.Greeting()
		// log.Println("\n------------------Send OpenChan Tx request--------------------")
		// OpenChanTx := sendOpenChanTxs(client, common.HexToAddress(config.ContractAddress))
		// startTime = time.Now()

		// select {
		// case hub.Send_fn <- OpenChanTx:
		// 	log.Println("Request message sent successfully")
		// default:
		// 	log.Println("Failed to send request message: channel is full or closed")
		// }
		// fmt.Println("------------------------------------------------------\n")

		// // Benchmarking for Geth nodes
		// log.Println("\n------------------Send OpenChan Tx request to geth--------------------")
		// benchmarking.GethSyncTx(client, common.HexToAddress(config.ContractAddress))
		// fmt.Println("------------------------------------------------------\n")

		// log.Println("\n------------------Send OpenChan Tx request to geth async--------------------")
		// benchmarking.GethAsyncTx(client, common.HexToAddress(config.ContractAddress))
		// fmt.Println("------------------------------------------------------\n")

		var testDuration time.Duration
		const numRequest = 100
		for i := 0; i < numRequest; i++ {
			log.Println("\n------------------Send BalanceChecking request--------------------")
			reqGeneTimer := time.Now()
			balanceCheckingReq := sendRequests(client, client.Amount+20)
			reqTime := time.Since(reqGeneTimer)
			testDuration += reqTime
			
			startTimeQuerySslip = time.Now()
			select {
			case hub.Send_fn <- balanceCheckingReq:
				log.Println("Request message sent successfully")
			default:
				log.Println("Failed to send request message: channel is full or closed")
			}
			fmt.Println("------------------------------------------------------\n")
			time.Sleep(500 * time.Millisecond) // Sleep for 100 milliseconds
		}


		var totalDurationResponse time.Duration
		var totalDurationVerification time.Duration
		var totalProofVerification time.Duration
		var totalRestVerification time.Duration
		for _, d := range durationQueryTimes {
			totalDurationResponse += d
		}
		for _, d := range durationQueryVerifyTimes {
			totalDurationVerification += d
		}
		for _, d := range durationQueryProofVerifyTimes {
			totalProofVerification += d
		}

		for _, d := range durationQueryRestVerifyTimes {
			totalRestVerification += d
		}

		averageDurationResponse := totalDurationResponse / time.Duration(len(durationQueryTimes))
		averageDurationVerify := totalDurationVerification / time.Duration(len(durationQueryVerifyTimes))
		log.Println("Average request generation time: ", testDuration / time.Duration(numRequest))
		log.Printf("Average duration for %d requests: %s", len(durationQueryTimes), averageDurationResponse)
		log.Printf("Average duration for %d requests got verified: %s", len(durationQueryVerifyTimes), averageDurationVerify)
		log.Println("Proof Verification time: ", totalProofVerification / time.Duration(len(durationQueryProofVerifyTimes)))
		log.Println("Rest Verification: ", totalRestVerification / time.Duration(len(durationQueryRestVerifyTimes)))
		// Benchmarking: geth requests

		// benchmarking.GethSyncQuery(client)

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

func sendOpenChanTxs(client *client.Client, contractAddress common.Address) []byte {
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
	openChanSignTx := protocol.OpenChanTx(client, fnAddr, deposit, contractAddress)

	msg := protocol.GenerateRequest(client, 20, client.Amount+100, openChanSignTx, blockHeader.Hash())

	client.Amount += 100
	msgWType := append([]byte("TX:"), msg...)

	log.Println("Sending: ", string(msgWType))

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
	log.Println("txProofRLP: ", txProofRLP)
	log.Println("txRLP: ", txRLP)
	log.Println("proof: ", proof.Serialize())

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
