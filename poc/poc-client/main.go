package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"poc-server/resmsg"
	"sync"

	"poc-client/client"
	"poc-client/connection"
	"poc-client/handlers"
	"poc-client/hub"
	"poc-client/hub/wsClient"
	"poc-client/msg/request"
	"poc-client/protocol"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
	config, _ := loadConfig("localConfig.json")

	// Generate a new client from the local configuration file
	log.Println("\n-----------------Generate a new client-------------------")
	client, err := client.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	fmt.Println("---------------------------------------------------------")

	// Start the hub
	hub := hub.NewHub()
	go hub.Run()

	// Setup connection with the PoC server
	server_wsConn, err := connection.ConnectToServer(client)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	server := &wsClient.Client{Conn: server_wsConn, Send: make(chan []byte, 10000), Receive: make(chan []byte, 10000)}
	hub.Set_fn <- server

	go server.WritePump()
	go server.ReadPump()

	go func() {
		for msg := range server.Receive {
			err := handlers.HandleMesssage(msg, client)
			if err != nil {
				log.Fatalf("Failed to handle message: %v", err)
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
		log.Println("\n------------------Send OpenChan Tx request--------------------")
		OpenChanTx := sendOpenChanTxs(client, common.HexToAddress(config.ContractAddress))
		select {
		case hub.Send_fn <- OpenChanTx:
			log.Println("Request message sent successfully")
		default:
			log.Println("Failed to send request message: channel is full or closed")
		}
		fmt.Println("------------------------------------------------------\n")

		for i := 0; i < 1; i++ {
			log.Println("\n------------------Send BalanceChecking request--------------------")
			balanceCheckingReq := sendRequests(client, client.Amount+20)
			select {
			case hub.Send_fn <- balanceCheckingReq:
				log.Println("Request message sent successfully")
			default:
				log.Println("Failed to send request message: channel is full or closed")
			}
			fmt.Println("------------------------------------------------------\n")
		}

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

func VerifyFraudProof(resMsg resmsg.ResponseMsg) bool {

	return true

}

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
