package client

import (
	"crypto/ecdsa"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (c *Client) Init(p string, bcWsEndpoint string, bcRpcEndpoint string, serverEndpoint string, contractAddress string) *Client {
	// generate the necessary account information
	privateKey := &ecdsa.PrivateKey{}

	if _, err := os.Stat(p); os.IsNotExist(err) {
		log.Println("generating a new private key")
		privateKey, err = crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		err = crypto.SaveECDSA(p, privateKey)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("loading existing private key")
		privateKey, err = crypto.LoadECDSA(p)
		if err != nil {
			log.Fatal(err)
		}
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// TODO: REMOVE
	privateKeyBytes := crypto.FromECDSA(privateKey)
	log.Println("privateKey: ", privateKeyBytes)
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	log.Println("publicKey: ", publicKeyBytes)
	addressHex := address.Hex()
	log.Println("address: ", addressHex)
	// TODO: REMOVE ABOVE

	return &Client{
		BcWsEndpoint:    bcWsEndpoint,
		BcRpcEndpoint:   bcRpcEndpoint,
		ServerEndpoint:  serverEndpoint,
		PrivateKey:      privateKey,
		PublicKey:       publicKeyECDSA,
		ServerPublicKey: nil,
		ContractAddress: contractAddress,
		Address:         address,
		ConnectFN:       common.Address{},
		FeeStandards:    0,
		Secret:          nil,
		ChannelID:       nil,
		StartTime:       time.Now(),
		Amount:          0,
		Duration:        0,
		LastSignedProof: nil,
		BlockHeader:     nil,
		Checkpoint:      nil,
		Step:            0,
	}
}
