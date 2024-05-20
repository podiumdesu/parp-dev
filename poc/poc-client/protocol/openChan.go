package protocol

import (
	"context"
	"crypto/ecdsa"

	// "encoding/hex"
	// "fmt"
	"log"
	"math/big"
	"math/rand"
	pocClient "poc-client/client"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func OpenChanTx(c *pocClient.Client, fnAddr common.Address, deposit *big.Int, contractAddress common.Address) []byte {
	// Setup Account
	wsEndpoint := c.BcWsEndpoint
	client, err := ethclient.Dial(wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	privateKey := c.PrivateKey
	// privateKey, err := crypto.HexToECDSA("535468b2ddcd8fc2b87c3b825922880c0d9f546095908bb924f1053e39852d5a")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasLimit := uint64(1000000)

	// Contract settings
	const contractABI = `[{"inputs":[{"indexed":true,"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"ChannelOpened","type":"event"},{"inputs":[{"internalType":"address","name":"addr","type":"address"}],"name":"balance","outputs":[{"internalType":"uint256","name":"bal","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"closeChan","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"confirmClosure","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"}],"name":"greeting","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"senderDeposit","type":"uint256"}],"name":"openChan","outputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"paychanCheck","outputs":[{"components":[{"internalType":"bytes32","name":"id","type":"bytes32"},{"internalType":"address payable","name":"sender","type":"address"},{"internalType":"address payable","name":"recipient","type":"address"},{"internalType":"uint256","name":"senderDeposit","type":"uint256"},{"internalType":"uint256","name":"startTime","type":"uint256"},{"internalType":"uint256","name":"status","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"},{"internalType":"uint256","name":"disputeStartTime","type":"uint256"},{"internalType":"uint256","name":"disputeDuration","type":"uint256"},{"internalType":"bool","name":"senderConfirm","type":"bool"},{"internalType":"bool","name":"recipientConfirm","type":"bool"}],"internalType":"struct paychan.PayChan","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"paychanSelectedArguments","outputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"rec","type":"address"},{"internalType":"uint256","name":"status","type":"uint256"},{"internalType":"uint256","name":"senderB","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	abiObj, _ := abi.JSON(strings.NewReader(contractABI))

	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// contractAddress := common.HexToAddress("0xA6f5B646d18665f4c86aCD0e86659b44A1367af9")
	// msgValue := big.NewInt(200000)

	// OpenChannel
	data, err := abiObj.Pack("openChan", fnAddr, deposit)

	if err != nil {
		log.Fatalf("Failed to pack data for openChan: %v", err)
	}

	// Create transaction
	tx := types.NewTransaction(nonce, contractAddress, deposit, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	// log.Print(fmt.Sprintf("chainID: %s", chainID))
	if err != nil {
		log.Fatal(err)
	}

	ts := types.Transactions{signedTx}
	rawTxBytes, _ := rlp.EncodeToBytes(ts[0])
	// rawTxHex := hex.EncodeToString(rawTxBytes)

	// fmt.Printf(rawTxHex) // f86...772
	return rawTxBytes
}

func OpenChanTxToGeth(c *pocClient.Client, fnAddr common.Address, deposit *big.Int, contractAddress common.Address, nonce uint64) *types.Transaction {
	// Setup Account
	wsEndpoint := c.BcWsEndpoint
	client, err := ethclient.Dial(wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	privateKey := c.PrivateKey
	// privateKey, err := crypto.HexToECDSA("535468b2ddcd8fc2b87c3b825922880c0d9f546095908bb924f1053e39852d5a")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	// }
	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	randN := rand.Intn(1000)
	deposit.Add(deposit, big.NewInt(int64(randN)))
	gasLimit := uint64(1000000)

	// Contract settings
	const contractABI = `[{"inputs":[{"indexed":true,"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"ChannelOpened","type":"event"},{"inputs":[{"internalType":"address","name":"addr","type":"address"}],"name":"balance","outputs":[{"internalType":"uint256","name":"bal","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"closeChan","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"confirmClosure","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"}],"name":"greeting","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"senderDeposit","type":"uint256"}],"name":"openChan","outputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"paychanCheck","outputs":[{"components":[{"internalType":"bytes32","name":"id","type":"bytes32"},{"internalType":"address payable","name":"sender","type":"address"},{"internalType":"address payable","name":"recipient","type":"address"},{"internalType":"uint256","name":"senderDeposit","type":"uint256"},{"internalType":"uint256","name":"startTime","type":"uint256"},{"internalType":"uint256","name":"status","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"},{"internalType":"uint256","name":"disputeStartTime","type":"uint256"},{"internalType":"uint256","name":"disputeDuration","type":"uint256"},{"internalType":"bool","name":"senderConfirm","type":"bool"},{"internalType":"bool","name":"recipientConfirm","type":"bool"}],"internalType":"struct paychan.PayChan","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"paychanSelectedArguments","outputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"rec","type":"address"},{"internalType":"uint256","name":"status","type":"uint256"},{"internalType":"uint256","name":"senderB","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	abiObj, _ := abi.JSON(strings.NewReader(contractABI))

	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// contractAddress := common.HexToAddress("0xA6f5B646d18665f4c86aCD0e86659b44A1367af9")
	// msgValue := big.NewInt(200000)

	// OpenChannel
	data, err := abiObj.Pack("openChan", fnAddr, deposit)

	if err != nil {
		log.Fatalf("Failed to pack data for openChan: %v", err)
	}

	// Create transaction
	tx := types.NewTransaction(nonce, contractAddress, deposit, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

	if err != nil {
		log.Fatal(err)
	}

	return signedTx
	// ts := types.Transactions{signedTx}
	// return ts
	// rawTxBytes, _ := rlp.EncodeToBytes(ts[0])
	// rawTxHex := hex.EncodeToString(rawTxBytes)

	// fmt.Printf(rawTxHex) // f86...772
	// return rawTxBytes
}
