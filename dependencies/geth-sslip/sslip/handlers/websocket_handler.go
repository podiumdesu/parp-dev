package handlers

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/poc-client/msg/handshake"
	"github.com/ethereum/go-ethereum/poc-client/msg/request"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/sslip/manager"
	"github.com/ethereum/go-ethereum/sslip/mpt"
	"github.com/ethereum/go-ethereum/sslip/resmsg"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

// type FullMemDB struct {
// 	*memorydb.Database
// }

// func (db *FullMemDB) Ancient(kind string, number uint64) ([]byte, error) {
// 	// Implement or stub as necessary; returning nil or an error based on your scenario
// 	return nil, nil
// }

// func (db *FullMemDB) AncientDatadir() (string, error) {
// 	// Implement or stub as necessary; returning nil or an error based on your scenario
// 	return "", nil
// }
// func (db *FullMemDB) AncientRange(string, uint64, uint64, uint64) ([][]byte, error) {
// 	// Implement or stub as necessary; returning nil or an error based on your scenario
// 	return nil, nil
// }
// func (db *FullMemDB) AncientSize(string) (uint64, error) {
// 	// Implement or stub as necessary; returning nil or an error based on your scenario
// 	return 0, nil
// }
// func (db *FullMemDB) Ancients() (uint64, error) {
// 	// Implement or stub as necessary; returning nil or an error based on your scenario
// 	return 0, nil
// }
// func (db *FullMemDB) HasAncient(string, uint64) (bool, error) {
// 	// Implement or stub as necessary; returning nil or an error based on your scenario
// 	return false, nil
// }
// func (db *FullMemDB) MigrateTable(string, func([]byte) ([]byte, error)) {
// 	// Implement or stub as necessary; returning nil or an error based on your scenario
// }

func HandleWebSocket(m *manager.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var upgrader = websocket.Upgrader{
		// 	CheckOrigin: func(r *http.Request) bool {
		// 		return true
		// 	},
		// }
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer conn.Close()

		clientID := (strings.Split(r.URL.Path, "/"))[2]
		// if clientID == "" {
		// 	http.Error(w, "Client ID not provided", http.StatusBadRequest)
		// 	return
		// }
		log.Print("Client ID when handling: ", clientID)

		m.AddClient(clientID, conn)

		m.PrintClientsMap()

		client := m.GetClient(clientID)

		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Printf("Client %s disconnected", clientID)
				} else {
					log.Println("Read:", err)
				}
				break
			}

			parts := strings.SplitN(string(msg), ":", 2)
			if len(parts) < 2 {
				log.Println("Invalid message format")
				return
			}

			header := parts[0]
			body := parts[1]

			switch header {

			case "HANDSHAKE":
				log.Println("----------------HANDSHAKE from clientID ", clientID, "---------------------")
				log.Println("Recv HANDSHAKE from clientID ", clientID)

				var hsMsg handshake.Msg
				err := json.Unmarshal([]byte(body), &hsMsg)

				if err != nil {
					log.Println("Unmarshal error: ", err)
					break
				}
				jsonHsMsg, _ := json.Marshal(hsMsg)
				log.Println("Handshake msg: ", string(jsonHsMsg))

				m.SetClientPubKB(clientID, hsMsg.Body.PubKB)
				m.SetContractAddress(hsMsg.Body.ContractAddress)
				// log.Println("pubk has been set: ", client.PubKeyB)

				go func() {
					handshakeMsg := resmsg.HandshakeMsg{
						Type:             "HANDSHAKE-CONFIRMED",
						ServerPublicKeyB: m.PubKeyBytes(),
					}
					client.Send(handshakeMsg.Bytes())
				}()
				fmt.Println("------------------------------------------------------\n")
				// log.Println("---------Connection with clientID ", clientID, " established---------")

			case "SIG":
				log.Println("----------------Request from clientID ", clientID, "---------------------")

				log.Println("Recv SIG from clientID ", clientID)

				// Unmarshal request body
				var req request.RequestMsg
				err := json.Unmarshal([]byte(body), &req)
				if err != nil {
					log.Println("Unmarshal error: ", err)
					break
				}
				jsonReq, _ := json.Marshal(req)
				log.Println("Request: ", string(jsonReq))
				requestByte := req.ReqByte
				reqString := string(requestByte)
				log.Println(reqString)

				// form request
				resp, err := http.Get(reqString)
				if err != nil {
					log.Println("Error: ", err)
					return
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				log.Println("body: ", string(body))

				sigFlag, _ := m.VerifyRequestWithSig(clientID, req)
				var msg resmsg.ServerMsg
				if sigFlag {
					log.Println("PASS: Signature verified")
					msg = resmsg.ServerMsg{
						Type: "SigCheck",
						Info: []byte("Passed"),
					}
				} else {
					msg = resmsg.ServerMsg{
						Type: "SigCheck",
						Info: []byte("WRONG signature"),
					}
				}
				client.Send(msg.Bytes())
				fmt.Println("------------------------------------------------------\n")

			case "TX":
				handler_tx(clientID, body, m, conn, mt)

			case "REQ":
				handler_req(clientID, body, m)
			default:

				log.Println("\n----------------Message from clientID ", clientID, "---------------------")

				log.Printf("Recv from clientID %s: %s", clientID, msg)
				serverMsg := resmsg.ServerMsg{
					Type: "info",
					Info: []byte("server:" + string(msg)),
				}
				client.Send(serverMsg.Bytes())
				if err != nil {
					log.Println("Write: ", err)
					break
				}
				fmt.Println("------------------------------------------------------\n")
				// log.Println("-------------------------------------------------------------------------")

			}

		}

		m.RemoveClient(clientID)
		m.PrintClientsMap()
	}
}

func printResponseMsg(msg resmsg.ResponseMsg) {
	fmt.Println("ResponseMsg:")
	fmt.Println("-------------")

	// Print the Type
	fmt.Printf("Type: %s\n", msg.Type)

	// Print the ChannelId (as a hexadecimal string for readability)
	fmt.Printf("ChannelId: %s\n", "0x"+hex.EncodeToString(msg.ChannelId[:]))

	// Print the Amount
	fmt.Printf("Amount: %d\n", msg.Amount)

	// Print the SignedReqBody
	fmt.Printf("SignedReqBody: %s\n", string(msg.SignedReqBody))

	// Print the CurrentBlockHeight
	fmt.Printf("CurrentBlockHeight: %d\n", msg.CurrentBlockHeight)

	// Print the ReturnValue
	fmt.Printf("ReturnValue: %s\n", string(msg.ReturnValue))

	// Print the Proof array
	fmt.Println("Proof:")
	for i, proof := range msg.Proof {
		fmt.Printf("  Proof[%d]: %s\n", i, proof)
	}

	// Print the TxHash (as a hexadecimal string for readability)
	fmt.Printf("TxHash: %s\n", "0x"+hex.EncodeToString(msg.TxHash[:]))

	// Print the TxIdx
	fmt.Printf("TxIdx: %d\n", msg.TxIdx)

	// Print the Signature
	fmt.Printf("Signature: %s\n", string(msg.Signature))

	fmt.Println()
}

func fromEthTransaction(t *types.Transaction) *mpt.Transaction {
	v, r, s := t.RawSignatureValues()
	return &mpt.Transaction{
		AccountNonce: t.Nonce(),
		Price:        t.GasPrice(),
		GasLimit:     t.Gas(),
		Recipient:    t.To(),
		Amount:       t.Value(),
		Payload:      t.Data(),
		V:            v,
		R:            r,
		S:            s,
	}
}

func trieWithBlockTxs(txs []*types.Transaction, txRootHash common.Hash, txHash common.Hash, block *types.Block) {

	// txs := transactionsJSON()

	log.Println(len(txs))
	mptTrie := mpt.NewTrie()
	for i, tx := range txs {
		key, _ := rlp.EncodeToBytes(uint(i))

		transaction := fromEthTransaction(tx)

		rlp, _ := transaction.GetRLP()

		mptTrie.Put(key, rlp)
	}

	// hasher := trie.NewStackTrie(nil)

	// txRootHash := fmt.Sprintf("%x", types.DeriveSha(types.Transactions(txs), hasher))
	fmt.Printf("txRootHash: %v\n", txRootHash)
	fmt.Printf("%x", mptTrie.Hash())

	// generate the proof and verify it
	key, _ := rlp.EncodeToBytes(uint(0))
	proof, found := mptTrie.Prove(key)

	fmt.Printf("proof: %x, found: %v\n", proof, found)

	txRLP, _ := mpt.VerifyProof(mptTrie.Hash(), key, proof)
	rlp, _ := fromEthTransaction(txs[0]).GetRLP()

	fmt.Println(txRLP)
	fmt.Println(rlp)
}

func generateProof(block *types.Block, txHash common.Hash) (mpt.Proof, []byte) {
	txs := block.Transactions()
	idx := -1
	for index, tx := range txs {
		if tx.Hash() == txHash {
			idx = index
		}
	}
	if idx < 0 {
		return nil, nil
	}

	mptTrie := mpt.NewTrie()
	for i, tx := range txs {
		key, _ := rlp.EncodeToBytes(uint(i))

		transaction := fromEthTransaction(tx)

		rlp, _ := transaction.GetRLP()

		mptTrie.Put(key, rlp)
	}

	// generate the proof and verify it
	key, _ := rlp.EncodeToBytes(uint(idx))
	proof, found := mptTrie.Prove(key)
	proofSize := len(proof.Serialize()[0])
	log.Println("proofSize: ", proofSize)

	fmt.Printf("proof: %x, found: %v\n", proof, found)

	return proof, key
}

func verifyProof(txHash common.Hash, proof mpt.Proof, blockNr *big.Int, key []byte) bool {
	wsEndpoint := "ws://localhost:8100"
	bcClient, err := ethclient.Dial(wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	// query the block information
	block, _ := bcClient.HeaderByNumber(context.Background(), blockNr)
	txRootHash := block.TxHash
	tx, _, _ := bcClient.TransactionByHash(context.Background(), txHash)
	txRLP, _ := rlp.EncodeToBytes(tx)
	txProofRLP, _ := mpt.VerifyProof(txRootHash[:], key, proof)
	log.Println("txProofRLP: ", txProofRLP)
	log.Println("txRLP: ", txRLP)
	log.Println("proof: ", proof.Serialize())

	return bytes.Equal(txRLP, txProofRLP)

	// log.Println("Block Hash: ", block)
	// txRootHash common.Hash, idx uint, proof mpt.Proof
	// key, _ := rlp.EncodeToBytes(uint(idx))
	// txRLP, _ := mpt.VerifyProof(txRootHash[:], key, proof)
	// rlp, _ := rlp.EncodeToBytes(tx)
	// log.Println(txRLP)
	// log.Println(rlp)
	// return bytes.Equal(txRLP, rlp)
}

// txRLP, _ := mpt.VerifyProof(mptTrie.Hash(), key, proof)
// rlp, _ := fromEthTransaction(txs[0]).GetRLP()

// fmt.Println(txRLP)
// fmt.Println(rlp)
