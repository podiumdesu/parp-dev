package manager

import (
	"crypto/ecdsa"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	mClient "github.com/ethereum/go-ethereum/parp/clients"

	"github.com/gorilla/websocket"
)

type Manager struct {
	clientsMap      map[string]*mClient.Client
	mu              sync.Mutex
	PrivateKey      *ecdsa.PrivateKey
	PublicKey       *ecdsa.PublicKey
	ContractAddress string
}

type ClientManager interface {
	AddClient(id string, conn *websocket.Conn)
	RemoveClient(id string)
	GetClient(id string) *mClient.Client
	SetClientPubK(id string, pubK []byte)
	SetClientChannel(id string, channelID common.Hash)
	PrintClientsMap()
}

func NewManager() *Manager {
	privateKeyECDSA, _ := crypto.HexToECDSA("bcd5c542c981dbb7cee1f3352fcee082581b4a323bf5cbff105aa84fa718f690")
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	return &Manager{
		clientsMap:      make(map[string]*mClient.Client),
		PrivateKey:      privateKeyECDSA,
		PublicKey:       publicKeyECDSA,
		ContractAddress: "",
	}
}
