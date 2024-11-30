package utils

import (
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func BhRlpBytes(blockHeader *types.Header) string {
	// Encode the ResponseMsg to RLP bytes
	serializedBytes, err := rlp.EncodeToBytes(blockHeader)
	if err != nil {
		log.Fatalf("Failed to RLP encode ResponseMsg: %v", err)
	}

	// Convert the serialized RLP bytes to a hex string, prefixed with "0x"
	hexString := "0x" + hex.EncodeToString(serializedBytes)

	return hexString
}
