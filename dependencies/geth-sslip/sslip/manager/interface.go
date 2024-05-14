package manager

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	mClient "github.com/ethereum/go-ethereum/sslip/clients"

	"github.com/gorilla/websocket"
)

func (m *Manager) AddClient(id string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clientsMap[id] = mClient.New(id, conn)
}

func (m *Manager) RemoveClient(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clientsMap, id)
}

func (m *Manager) SetClientPubKB(id string, pubKB []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clientsMap[id].PubKeyB = pubKB
}

func (m *Manager) SetClientChannelID(id string, channelID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clientsMap[id].ChannelID = channelID
}

func (m *Manager) GetClientChannelID(id string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clientsMap[id].ChannelID
}
func (m *Manager) GetClient(id string) *mClient.Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clientsMap[id]
}

func (m *Manager) PrintClientsMap() {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Print("\n\nClients Map: \n")
	for id, client := range m.clientsMap {
		log.Print(id, client)
	}
}

func (m *Manager) PrivKeyBytes() []byte {
	return crypto.FromECDSA(m.PrivateKey)
}

func (m *Manager) PubKeyBytes() []byte {
	return crypto.FromECDSAPub(m.PublicKey)
}

func (m *Manager) Address() common.Address {
	address := crypto.PubkeyToAddress(*m.PublicKey)
	return address
}

func (m *Manager) Sign(hash []byte) []byte {
	sig, err := crypto.Sign(hash, m.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	return sig
}
func (m *Manager) Verify(hash []byte, sig []byte) bool {
	signatureNoRecoverID := sig[:len(sig)-1]
	pubKeyByte := crypto.FromECDSAPub(m.PublicKey)
	return crypto.VerifySignature(pubKeyByte, hash, signatureNoRecoverID)
}
