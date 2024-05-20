package benchmarking

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"poc-client/client"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GethSyncTx(client *client.Client, contractAddress common.Address) {
	var totalDuration time.Duration
	totalNum := 10
	var successNum int
	for i := 0; i < totalNum; i++ {
		fmt.Printf("Executing transaction %d...\n", i+1)
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
		duration, err := sendOpenChanTxsToGeth(client, contractAddress, nonce)
		if duration != 0 {
			successNum += 1
		}
		if err != nil {
			log.Fatalf("Failed to send transaction %d: %v", i+1, err)
		}

		fmt.Printf("Transaction %d took %s\n", successNum, duration)
		totalDuration += duration
	}

	fmt.Printf("In total %d / %d successfully executed\n", successNum, totalNum)
	averageDuration := totalDuration / 10
	fmt.Printf("Average time taken for 10 transactions: %s\n", averageDuration)
}
