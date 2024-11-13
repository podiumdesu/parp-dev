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
	signature, err := crypto.Sign(data, privateKey)
	if err != nil {
		return nil, err
	}

	// Adjust v value, Ethereum uses 27 or 28
	if signature[64] != 27 && signature[64] != 28 {
		signature[64] += 27
	}

	return signature, nil
}

func (c *Client) Verify(hash []byte, sig []byte) bool {
	signatureNoRecoverID := sig[:len(sig)-1]
	pubKeyByte := crypto.FromECDSAPub(c.PublicKey)
	return crypto.VerifySignature(pubKeyByte, hash, signatureNoRecoverID)
}

func (c *Client) PrivKeyBytes() []byte {
	return crypto.FromECDSA(c.PrivateKey)
}

func (c *Client) PubKeyBytes() []byte {
	return crypto.FromECDSAPub(&c.PrivateKey.PublicKey)
}

func (c *Client) AddrHex() string {
	return c.Address.Hex()
}
