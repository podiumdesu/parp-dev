package client

import "github.com/ethereum/go-ethereum/crypto"

func (c *Client) PrivKeyBytes() []byte {
	return crypto.FromECDSA(c.PrivateKey)
}

func (c *Client) PubKeyBytes() []byte {
	return crypto.FromECDSAPub(&c.PrivateKey.PublicKey)
}

func (c *Client) AddrHex() string {
	return c.Address.Hex()
}
