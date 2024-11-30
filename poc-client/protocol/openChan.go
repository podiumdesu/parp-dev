package protocol

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func OpenChanTx(bcClient *ethclient.Client, privateKey *ecdsa.PrivateKey, fnAddr common.Address, deposit *big.Int, nonce uint64, contractAddress common.Address) ([]byte, common.Hash) {

	gasPrice, err := bcClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasLimit := uint64(10000000)

	// Contract settings
	contractABI := getAbi()
	abiObj, _ := abi.JSON(strings.NewReader(contractABI))

	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// OpenChannel
	data, err := abiObj.Pack("openChan", fnAddr, deposit)

	if err != nil {
		log.Fatalf("Failed to pack data for openChan: %v", err)
	}

	// Create transaction
	// Note: NewTransaction is deprecated, https://github.com/nnqq/geth-tx-hash-bug/blob/master/main.go

	tx := types.NewTransaction(nonce, contractAddress, deposit, gasLimit, gasPrice, data)
	chainID, err := bcClient.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

	if err != nil {
		log.Fatal(err)
	}

	ts := types.Transactions{signedTx}
	rawTxBytes, _ := rlp.EncodeToBytes(ts[0])

	txHash := signedTx.Hash()
	return rawTxBytes, txHash
}
