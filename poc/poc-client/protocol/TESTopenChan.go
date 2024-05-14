package protocol

// package main

// import (
// 	"context"
// 	"crypto/ecdsa"
// 	"encoding/hex"
// 	"fmt"
// 	"log"
// 	"math/big"
// 	"strings"

// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	"github.com/ethereum/go-ethereum/rlp"
// )

// func main() {

// 	// Connect to the endpoint
// 	wsEndpoint := "ws://localhost:8100"
// 	client, err := ethclient.Dial(wsEndpoint)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// set up account
// 	privateKey, err := crypto.HexToECDSA("535468b2ddcd8fc2b87c3b825922880c0d9f546095908bb924f1053e39852d5a")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
// 	}

// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
// 	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	gasPrice, err := client.SuggestGasPrice(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	gasLimit := uint64(1000000)

// 	// Contract settings
// 	const contractABI = `[{"inputs":[{"indexed":true,"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"ChannelOpened","type":"event"},{"inputs":[{"internalType":"address","name":"addr","type":"address"}],"name":"balance","outputs":[{"internalType":"uint256","name":"bal","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"closeChan","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"confirmClosure","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"}],"name":"greeting","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"senderDeposit","type":"uint256"}],"name":"openChan","outputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"paychanCheck","outputs":[{"components":[{"internalType":"bytes32","name":"id","type":"bytes32"},{"internalType":"address payable","name":"sender","type":"address"},{"internalType":"address payable","name":"recipient","type":"address"},{"internalType":"uint256","name":"senderDeposit","type":"uint256"},{"internalType":"uint256","name":"startTime","type":"uint256"},{"internalType":"uint256","name":"status","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"},{"internalType":"uint256","name":"disputeStartTime","type":"uint256"},{"internalType":"uint256","name":"disputeDuration","type":"uint256"},{"internalType":"bool","name":"senderConfirm","type":"bool"},{"internalType":"bool","name":"recipientConfirm","type":"bool"}],"internalType":"struct paychan.PayChan","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"paychanSelectedArguments","outputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"rec","type":"address"},{"internalType":"uint256","name":"status","type":"uint256"},{"internalType":"uint256","name":"senderB","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"}],"stateMutability":"view","type":"function"}]`
// 	abiObj, _ := abi.JSON(strings.NewReader(contractABI))

// 	if err != nil {
// 		log.Fatalf("Failed to parse ABI: %v", err)
// 	}

// 	contractAddress := common.HexToAddress("0x094D6cd9dA692A4c490C1F8AD3E74D089E3492D6")
// 	msgValue := big.NewInt(200000)

// 	// OpenChannel
// 	fnAddr := common.HexToAddress("0xA2131E7503F7Dd11ff5dAAC09fa7c301e7Fe0f30")

// 	data, err := abiObj.Pack("openChan", fnAddr, msgValue)

// 	if err != nil {
// 		log.Fatalf("Failed to pack data for openChan: %v", err)
// 	}

// 	log.Println(data)

// 	// Create transaction
// 	tx := types.NewTransaction(nonce, contractAddress, msgValue, gasLimit, gasPrice, data)
// 	chainID, err := client.NetworkID(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
// 	log.Print(fmt.Sprintf("chainID: %s", chainID))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ts := types.Transactions{signedTx}
// 	rawTxBytes, _ := rlp.EncodeToBytes(ts[0])
// 	rawTxHex := hex.EncodeToString(rawTxBytes)

// 	fmt.Printf(rawTxHex) // f86...772

// 	// Send the tx to the network
// 	txSend := new(types.Transaction)
// 	rlp.DecodeBytes(rawTxBytes, &txSend)

// 	err = client.SendTransaction(context.Background(), txSend)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("tx sent: %s", txSend.Hash().Hex())

// 	// event subscribe

// 	query := ethereum.FilterQuery{
// 		FromBlock: big.NewInt(260),
// 		Addresses: []common.Address{contractAddress},
// 	}

// 	logs, err := client.FilterLogs(context.Background(), query)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// logs := make(chan types.Log)
// 	// sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// for {
// 	// 	select {
// 	// 	case err := <-sub.Err():
// 	// 		log.Fatal(err)
// 	// 	case vLog := <-logs:
// 	// 		fmt.Println(vLog) // pointer to event log
// 	// 	}
// 	// }

// 	contractAbi, err := abi.JSON(strings.NewReader(string(contractABI)))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, vLog := range logs {
// 		fmt.Println(vLog.BlockHash.Hex()) // 0x3404b8c050aa0aacd0223e91b5c32fee6400f357764771d0684fa7b3f448f1a8
// 		fmt.Println(vLog.BlockNumber)     // 2394201
// 		fmt.Println(vLog.TxHash.Hex())    // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6

// 		fmt.Println(vLog.Data)
// 		// event := struct {
// 		// 	Key   [32]byte
// 		// 	Value [32]byte
// 		// }{}
// 		event, err := contractAbi.Unpack("ChannelOpened", vLog.Data)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		// fmt.Println(string(event))          // foo
// 		log.Println("logs")
// 		fmt.Println(event) // foo
// 		// fmt.Println(string(event.Value[:])) // bar

// 		var topics [4]string
// 		for i := range vLog.Topics {
// 			topics[i] = vLog.Topics[i].Hex()
// 		}

// 		fmt.Println(topics[0]) // 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4
// 	}

// 	// Get Receipt
// 	// receipt, err := client.TransactionReceipt(context.Background(), txSend.Hash())
// 	// if err != nil {
// 	// 	log.Fatal(err) // Handle the error properly; it might be due to the receipt not being available yet.
// 	// }

// 	// if receipt.Status == types.ReceiptStatusFailed {
// 	// 	fmt.Println("Transaction failed")
// 	// 	// Optionally, proceed to retrieve more information on why it failed.
// 	// } else {
// 	// 	fmt.Println("Transaction succeeded")
// 	// }

// 	// simulate the reason
// 	// msg := ethereum.CallMsg{
// 	// 	From:     fromAddress,    // The sender of the transaction
// 	// 	To:       txSend.To(),    // The recipient (can be a contract)
// 	// 	GasPrice: gasPrice,       // Just reuse the gas price from the transaction
// 	// 	Gas:      gasLimit,       // Reuse the gas limit from the transaction
// 	// 	Value:    txSend.Value(), // The value transferred
// 	// 	Data:     txSend.Data(),  // The data sent in the transaction
// 	// }

// 	// result, err := client.CallContract(context.Background(), msg, nil) // `nil` means the latest block
// 	// if err != nil {
// 	// 	log.Println("Failed to retrieve revert reason:", err)
// 	// } else if len(result) > 0 {
// 	// 	log.Println("Revert reason:", string(result))
// 	// } else {
// 	// 	log.Println("Transaction failed without a revert reason")
// 	// }
// }
