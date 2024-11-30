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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func TransferTokenTx(bcClient *ethclient.Client, privateKey *ecdsa.PrivateKey, recipient common.Address, amount *big.Int, contractAddress common.Address) ([]byte, common.Hash) {
	// Extract public key and derive sender address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get the nonce for the sender address
	nonce, err := bcClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Get the suggested gas price
	gasPrice, err := bcClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Define gas limit for the transaction
	gasLimit := uint64(100000) // Adjust as necessary for the contract

	// ERC20 transfer function ABI
	const contractABI = `[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"}]`
	abiObj, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Pack the function call with parameters
	data, err := abiObj.Pack("transfer", recipient, amount)
	if err != nil {
		log.Fatalf("Failed to pack data for transfer: %v", err)
	}

	// Create the transaction
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, data)
	chainID, err := bcClient.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Encode the transaction to RLP format
	rawTxBytes, _ := rlp.EncodeToBytes(signedTx)

	txHash := signedTx.Hash()

	return rawTxBytes, txHash
}
