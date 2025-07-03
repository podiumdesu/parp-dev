package handlers

import (
	"bytes"
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/parp/mpt"
	"github.com/ethereum/go-ethereum/rlp"
)

func verifyProof(txHash common.Hash, proof mpt.Proof, blockNr *big.Int, key []byte) bool {
	wsEndpoint := "ws://localhost:8100"
	bcClient, err := ethclient.Dial(wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	// query the block information
	block, _ := bcClient.HeaderByNumber(context.Background(), blockNr)
	txRootHash := block.TxHash
	tx, _, _ := bcClient.TransactionByHash(context.Background(), txHash)
	txRLP, _ := rlp.EncodeToBytes(tx)
	txProofRLP, _ := mpt.VerifyProof(txRootHash[:], key, proof)

	return bytes.Equal(txRLP, txProofRLP)
}
