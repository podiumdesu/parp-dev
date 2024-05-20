package benchmarking

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"poc-client/client"
	"poc-client/protocol"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GethAsyncTx(client *client.Client, contractAddress common.Address) {

<<<<<<< HEAD
	totalNum := 300
=======
	totalNum := 400
>>>>>>> 2c759414002d851fc28e6fb6e2128aaa7093295a
	var totalDuration time.Duration
	var successNum int
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	// Channel to collect durations and errors
	durations := make(chan time.Duration, totalNum)
	errors := make(chan error, totalNum)

	privateKey := client.PrivateKey

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	wsEndpoint := client.BcWsEndpoint
	bcClient, err := ethclient.Dial(wsEndpoint)

	nonce, err := bcClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	// Sending multiple transactions concurrently
	for i := 0; i < totalNum; i++ {
		wg.Add(1)
		go func(txNum int) {
			defer wg.Done()
			fmt.Printf("Executing transaction %d...\n", txNum+1)
			log.Println(txNum)
			// log.Println(nonce+uint64(txNum))
			duration, err := sendOpenChanTxsToGeth(client, bcClient, contractAddress, nonce+uint64(txNum))
			mu.Lock()
			if duration != 0 {
				successNum++
			}
			mu.Unlock()
			durations <- duration
			errors <- err

		}(i)
	}

	// Wait for all transactions to complete
	wg.Wait()
	close(durations)
	close(errors)

	// Collect results
	for duration := range durations {
		totalDuration += duration
	}

	for err := range errors {
		if err != nil {
			log.Println("err", err)
		}
	}

	fmt.Printf("In total %d / %d successfully executed\n", successNum, totalNum)
	if successNum > 0 {
		averageDuration := totalDuration / time.Duration(successNum)
		fmt.Printf("Average time taken for successful transactions: %s\n", averageDuration)
	} else {
		fmt.Println("No successful transactions.")
	}
}

func sendOpenChanTxsToGeth(client *client.Client, bcClient *ethclient.Client, contractAddress common.Address, nonce uint64) (time.Duration, error) {
	log.Println("CALLED")
	fnAddr := common.HexToAddress("0xA2131E7503F7Dd11ff5dAAC09fa7c301e7Fe0f30")
	deposit := big.NewInt(200000)
	log.Println(nonce)
	signedTx := protocol.OpenChanTxToGeth(client, fnAddr, deposit, contractAddress, nonce)
	// // calculate the tx Size
	// signedTxBytes, err := signedTx.MarshalBinary()
	// if err != nil {
	// 	log.Fatalf("Failed to serialize the transaction: %v", err)
	// }
	// log.Println("Size of pure Tx: ", len(signedTxBytes))

	startTime := time.Now()
	err := bcClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send the transaction: %v", err)
		return 0, err
	}

	// rawTxBytes, err := signedTx.MarshalBinary()
	// if err != nil {
	//     log.Fatalf("Failed to serialize the transaction: %v", err)
	// 	return 0, err
	// }
	// rawTxHex := hexutil.Encode(rawTxBytes)
	// fmt.Println("Raw Transaction Hex:", rawTxHex)

	// fmt.Println("Transaction sent successfully. Tx Hash:", signedTx.Hash().Hex())

	receipt, err := waitForReceipt(bcClient, signedTx.Hash())
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
		return 0, err
	}

	// Record the end time
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("Transaction receipt: %+v\n", receipt)
	fmt.Printf("Time taken to get receipt: %s\n", duration)
	return duration, nil
}
func waitForReceipt(bcClient *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()
	for {
		receipt, err := bcClient.TransactionReceipt(ctx, txHash)
		if err == ethereum.NotFound {
			time.Sleep(100 * time.Millisecond) // Wait for 100 milliseconds before retrying
			continue
		} else if err != nil {
			return nil, err
		}
		return receipt, nil
	}
}
