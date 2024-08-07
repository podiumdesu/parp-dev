package cryptoutil

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

func PubkeyToHexAddr(pubKeyByte *ecdsa.PublicKey) string {
	return crypto.PubkeyToAddress(*pubKeyByte).Hex()
}
