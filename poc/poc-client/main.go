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

    "os"
	"poc-client/client"
	"poc-client/connection"
	"poc-client/hub"
	"poc-client/gethbenchmarking"
	"poc-client/hub/wsClient"
	"poc-client/msg/request"
	"poc-client/protocol"
	"poc-client/utils/cryptoutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
				log.Println(string(msg))
				var resMsg resmsg.ResponseMsg

				log.Println("Size of the Tx response: %d bytes", len(msg))
				err := json.Unmarshal(msg, &resMsg)
				if err != nil {
					log.Println("Unmarshal error: ", err)
					break
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
			case "responseSP":
				log.Println(string(msg))
				var resMsg resmsg.ResponseSPMsg
				err := json.Unmarshal(msg, &resMsg)
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

			default:
				log.Println("Unrecognized message type: ", serverMsg.Type)
			}
		}

	}()

	// wg.Add(1)

	// go func() {
	// 	log.Println("\n------------------Handshake-------------------")

	// 	defer wg.Done()
	// 	msg := client.InitHandshakeMsg(config.ContractAddress, 10, big.NewInt(100000), big.NewInt(200))
	// 	b := append([]byte("HANDSHAKE:"), msg.Bytes()...)

	// 	// b := []byte("HANDSHAKE:" + string(client.PubKeyBytes()))
	// 	log.Println("Sending: ", b)
	// 	hub.Send_fn <- b //[]byte("SIG:qwerrtqreqwrqwerqwrwerqwrewqtqwetqwrewqrqwerwqewqrqwer")
	// 	fmt.Println("----------------------------------------------\n")
	// 	hub.Send_fn <- []byte("FE: Connected to server")
	// }()

	// wg.Wait()

	// go func() {
	// 	for i := 0; i < 30; i++ {
	// 		log.Println("\n------------------Send OpenChan Tx request--------------------")
	// 		OpenChanTx := sendOpenChanTxs(client, common.HexToAddress(config.ContractAddress))
	// 		select {
	// 		case hub.Send_fn <- OpenChanTx:
	// 			log.Println("Request message sent successfully")
	// 		default:
	// 			log.Println("Failed to send request message: channel is full or closed")
	// 		}
	// 		fmt.Println("------------------------------------------------------\n")
	// 		time.Sleep(30 * time.Second)
	// 	}
	// }()

	
	go func() {
		currentDateTime := time.Now().Format("2006-01-02 15:04:05")

		// Open a file. Create the file if it does not exist and truncate it if it does.
		file, err := os.OpenFile("date.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close() // Ensure the file is closed after the function exits
	
		// Write the current date to the file
		_, err = fmt.Fprintf(file, "Start: Current Date: %s\n", currentDateTime)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	
		fmt.Println("Date written successfully")
		// logic for requests
		// for i := 0; i < 240; i++ {
		// 	log.Println("\n------------------Send BalanceChecking request--------------------")
		// 	balanceCheckingReq := sendRequests(client, client.Amount+20)
		// 	select {
		// 	case hub.Send_fn <- balanceCheckingReq:
		// 		log.Println("Request message sent successfully")
		// 	default:
		// 		log.Println("Failed to send request message: channel is full or closed")
		// 	}
		// 	fmt.Println("------------------------------------------------------\n")
		// 	time.Sleep(500 * time.Millisecond)
		// }

		gethbenchmarking.GethSyncQuery(client)

		currentDateTime = time.Now().Format("2006-01-02 15:04:05")
		_, err = fmt.Fprintf(file, "End: Current Date: %s\n", currentDateTime)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	
	}()
	// // Setup webpage for front end
	// http.HandleFunc("/", web.HomeHandler)

	// // Resolve websocket connection from the front end
	// http.HandleFunc("/ws-client-fb-connect", connection.ConnectToFE(hub))

	// // Send connection request from the front end to the server

	// go func() {
	// 	log.Println("Client server is running on port 8081")
	// 	log.Fatal(http.ListenAndServe(":8081", nil))
	// }()

	// // http.HandleFunc("/start-handshaking", connection.Handshaking(hub))

	// // Stop using front-end, just let the server generates the necessary data
	// time.Sleep(2 * time.Second)

	// go func() {
	// 	// hub.Send_fn <- []byte("Reee:")
	// 	// hub.Send_fn <- []byte("Send request:")

	// 	i := protocol.GenerateRequest(client, 20, big.NewInt(333))
	// 	log.Println("Request: ", string(i))
	// 	b := append([]byte("SIG:"), i...)
	// 	log.Println("Sending: ", string(b))
	// 	// log.Println("Sending: ", string(b))
	// 	// log.Println("Payload size: ", len(b))
	// 	hub.Send_fn <- b
	// 	// hub.Send_fn <- []byte("FFFFFF:")
	// }()

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
