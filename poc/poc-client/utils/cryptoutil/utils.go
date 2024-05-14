package cryptoutil

import (
	"bytes"
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func Sign(privateKey *ecdsa.PrivateKey, hashByte []byte) []byte {
	sig, err := crypto.Sign(hashByte, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	return sig
}

func Verify(pubKeyByte []byte, hashByte []byte, sig []byte) bool {
	// 1. Need to remove the recovery id
	// sigNoRecoverID := sig[:len(sig)-1]
	// return crypto.VerifySignature(pubKeyByte, hashByte, sigNoRecoverID)

	// // 2. crypto.Ecrecover
	// sigPubKeyB, _ := crypto.Ecrecover(hashByte, sig)
	// return bytes.Equal(sigPubKeyB, pubKeyByte)

	// 3. crypto.SigToPub
	sigPubKeyECDSA, _ := crypto.SigToPub(hashByte, sig)
	sigPubKeyB := crypto.FromECDSAPub(sigPubKeyECDSA)
	return bytes.Equal(sigPubKeyB, pubKeyByte)
}
