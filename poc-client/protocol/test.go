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
// 	"time"

// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	"github.com/ethereum/go-ethereum/rlp"
// 	// "golang.org/x/net/context"
// 	// for demo
// )

// func main() {

// 	const contractABI = `[{"inputs":[{"internalType":"string","name":"_version","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"key","type":"string"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"ItemSet","type":"event"},{"inputs":[{"internalType":"string","name":"","type":"string"}],"name":"items","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"key","type":"string"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`
// 	// const contractABI = `[{"inputs":[{"internalType":"string","name":"key","type":"string"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"setItem","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

// 	abiObj, err := abi.JSON(strings.NewReader(contractABI))

// 	wsEndpoint := "ws://localhost:8100"

// 	// Connect to the Ethereum node
// 	client, err := ethclient.Dial(wsEndpoint)

// 	// client, err := ethclient.Dial("https://rinkeby.infura.io")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

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

// 	auth := bind.NewKeyedTransactor(privateKey)
// 	auth.Nonce = big.NewInt(int64(nonce))
// 	auth.Value = big.NewInt(0)     // in wei
// 	auth.GasLimit = uint64(300000) // in units
// 	auth.GasPrice = gasPrice
// 	// auth.ChainID = 32382
// 	// instance, err := store.NewStore(address, client)

// 	address := common.HexToAddress("0x5BFCC7cbbC120A119B5D0361014C4B1DE00cCB2E")

// 	value := big.NewInt(0)

// 	key := "exampleKey"
// 	keyV := big.NewInt(300)
// 	data, err := abiObj.Pack("setItem", key, keyV)

// 	log.Println(data)
// 	tx := types.NewTransaction(nonce, address, value, auth.GasLimit, auth.GasPrice, data)
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

// 	fmt.Printf("tx sent: %s", txSend.Hash().Hex())

// 	//event subscribe

// 	// query := ethereum.FilterQuery{
// 	// 	// FromBlock: big.NewInt(2394201),
// 	// 	// ToBlock:   big.NewInt(2394201),
// 	// 	Addresses: []common.Address{
// 	// 		address,
// 	// 	},
// 	// }

// 	// logs, err := client.FilterLogs(context.Background(), query)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	contractAbi, err := abi.JSON(strings.NewReader(string(contractABI)))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var txReceipt *types.Receipt
// 	for txReceipt == nil {
// 		// Query the transaction receipt
// 		txReceipt, err = client.TransactionReceipt(context.Background(), txSend.Hash())
// 		if err != nil {
// 			log.Println("Waiting for transaction to be mined...")
// 			time.Sleep(5 * time.Second) // Adjust the sleep duration based on expected block time
// 		}
// 	}
// 	log.Printf("Transaction mined in block %d", txReceipt.BlockNumber.Uint64())

// 	// By directly monitor the mempool

// 	// receiptBytes, err := json.Marshal(txReceipt)
// 	// if err != nil {
// 	// 	log.Fatal("Failed to marshal receipt: ", err)
// 	// }
// 	// _ = conn.WriteMessage(mt, []byte(receiptBytes))

// 	for _, vLog := range txReceipt.Logs {
// 		fmt.Printf("Log Address: %s\n", vLog.Address.Hex())
// 		if len(vLog.Topics) > 0 {
// 			eventName, err := contractAbi.EventByID(vLog.Topics[0])
// 			if err != nil {
// 				log.Println("Error finding event name:", err)
// 				continue
// 			}
// 			fmt.Printf("Event Name: %s\n", eventName.Name)
// 			log.Println("Event Data: ", vLog.Data)

// 			// var results []interface{}

// 			testEvent := struct {
// 				Key   string
// 				Value *big.Int
// 			}{}
// 			err = contractAbi.UnpackIntoInterface(&testEvent, eventName.Name, vLog.Data)
// 			if err != nil {
// 				log.Println("Error unpacking log data:", err)
// 				continue
// 			}

// 			fmt.Printf("Event Data: %+v\n", testEvent)
// 		}
// 	}

// 	// for _, vLog := range logs {
// 	// 	fmt.Println(vLog.BlockHash.Hex()) // 0x3404b8c050aa0aacd0223e91b5c32fee6400f357764771d0684fa7b3f448f1a8
// 	// 	fmt.Println(vLog.BlockNumber)     // 2394201
// 	// 	fmt.Println(vLog.TxHash.Hex())    // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6

// 	// 	fmt.Println(vLog.Data)

// 	// 	// event := struct {
// 	// 	// 	Key   string
// 	// 	// 	Value int
// 	// 	// }{}
// 	// 	event, err := contractAbi.Unpack("ItemSet", vLog.Data)
// 	// 	if err != nil {
// 	// 		log.Fatal("err", err)
// 	// 	}
// 	// 	log.Println("log")
// 	// 	fmt.Println(event) // foo
// 	// 	// fmt.Println(string(event.Value[:])) // bar

// 	// 	testEvent := struct {
// 	// 		Key   string
// 	// 		Value *big.Int
// 	// 	}{}

// 	// 	_ = contractAbi.UnpackIntoInterface(&testEvent, "ItemSet", vLog.Data)
// 	// 	fmt.Println(testEvent.Key) // foo
// 	// 	fmt.Println(testEvent.Value)

// 	// 	// var topics [4]string
// 	// 	// for i := range vLog.Topics {
// 	// 	// 	topics[i] = vLog.Topics[i].Hex()
// 	// 	// }

// 	// 	// fmt.Println(topics[0]) // 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4
// 	// }

// 	eventSignature := []byte("ItemSet(bytes32,bytes32)")
// 	hash := crypto.Keccak256Hash(eventSignature)
// 	fmt.Println(hash.Hex()) // 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4
// }
