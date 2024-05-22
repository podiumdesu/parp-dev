package gethbenchmarking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"poc-client/client"
	"time"

	"poc-client/msg/request"

	"github.com/ethereum/go-ethereum/rpc"
)

var durationTotal time.Duration

func GethSyncQuery(client *client.Client) {
	totalNum := 240
	// Benchmarking: sslip requests
	wsEndpoint := client.BcWsEndpoint
	rpcClient, err := rpc.Dial(wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}


	request := request.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBalance",
		Params:  []interface{}{"0xA2131E7503F7Dd11ff5dAAC09fa7c301e7Fe0f30", "latest"},
		ID:      1,
	}
	for i := 0; i < totalNum; i++ {
		// log.Println("\n------------------Send BalanceChecking request --------------------")
		log.Println(i)
		var result json.RawMessage
		startTime := time.Now()
		err = rpcClient.CallContext(context.Background(), &result, request.Method, request.Params...)
		if err != nil {
			log.Fatal("Failed to send JSON-RPC request: ", err)
		}
		// Handle the response
		var balance string
		err = json.Unmarshal(result, &balance)
		if err != nil {
			log.Fatal("Failed to unmarshal response: ", err)
		}

		duration := time.Since(startTime)
		durationTotal += duration
		fmt.Printf("Balance: %s\n", balance)

		fmt.Println("------------------------------------------------------\n")
		time.Sleep(500 * time.Millisecond) // Sleep for 100 milliseconds
	}

	// log.Println(durationTotal)
	// log.Println(totalNum)
	// log.Println("Average query time: ", durationTotal/time.Duration(totalNum))
}
