package cryptoutil

import (
	"bytes"
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func Sign(privateKey *ecdsa.PrivateKey, hashByte []byte) []byte {
	sig, err := crypto.Sign(hashByte, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	return sig
}
func SignHash(privateKey *ecdsa.PrivateKey, hash common.Hash) []byte {
	signature, err := crypto.Sign(hash[:], privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// Adjust v value, Ethereum uses 27 or 28
	if signature[64] != 27 && signature[64] != 28 {
		signature[64] += 27
	}

	return signature
}

func Verify(pubKeyByte []byte, hashByte []byte, sig []byte) bool {
	// 1. Need to remove the recovery id
	// sigNoRecoverID := sig[:len(sig)-1]
	// return crypto.VerifySignature(pubKeyByte, hashByte, sigNoRecoverID)

	// // 2. crypto.Ecrecover
	// sigPubKeyB, _ := crypto.Ecrecover(hashByte, sig)
	// return bytes.Equal(sigPubKeyB, pubKeyByte)

	// 3. crypto.SigToPub

	sigCopy := append([]byte{}, sig...)

	if sigCopy[64] == 27 || sigCopy[64] == 28 {
		sigCopy[64] -= 27
	}

	sigPubKeyECDSA, _ := crypto.SigToPub(hashByte, sigCopy)
	sigPubKeyB := crypto.FromECDSAPub(sigPubKeyECDSA)
	return bytes.Equal(sigPubKeyB, pubKeyByte)
}

func VerifyHash(pubKeyByte []byte, hash common.Hash, sig []byte) bool {
	return Verify(pubKeyByte, hash[:], sig)
}
