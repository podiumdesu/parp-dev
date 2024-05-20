package client

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func (c *Client) Sign(hash []byte) []byte {
	// sig, err := crypto.Sign(hash, c.PrivateKey)
	sig, err := SignEthereumMessage(hash, c.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	return sig
}

func SignEthereumMessage(data []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	// Hash the data with the Ethereum prefix
	// hash := crypto.Keccak256Hash(data).Bytes()

	// Prefix the hash as per Ethereum's requirement
	// prefixedHash := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(hash)) + string(hash))).Bytes()

	// Sign the prefixed hash
	signature, err := crypto.Sign(data, privateKey)
	if err != nil {
		return nil, err
	}

	// Adjust v value, Ethereum uses 27 or 28
	if signature[64] != 27 && signature[64] != 28 {
		signature[64] += 27
	}

	// log.Println("--------------")
	// log.Println(hex.EncodeToString(data))
	// log.Println(hex.EncodeToString(signature))

	return signature, nil
}

func (c *Client) Verify(hash []byte, sig []byte) bool {
	signatureNoRecoverID := sig[:len(sig)-1]
	pubKeyByte := crypto.FromECDSAPub(c.PublicKey)
	return crypto.VerifySignature(pubKeyByte, hash, signatureNoRecoverID)
}
