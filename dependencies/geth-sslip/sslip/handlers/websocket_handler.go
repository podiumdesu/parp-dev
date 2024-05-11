package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/poc-client/msg/handshake"
	"github.com/ethereum/go-ethereum/poc-client/msg/request"
	"github.com/ethereum/go-ethereum/poc-client/utils/cryptoutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/sslip/manager"
	"github.com/ethereum/go-ethereum/sslip/mpt"
	"github.com/ethereum/go-ethereum/sslip/resmsg"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

const wsEndpointPort = 8100

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
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
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
				log.Println("\n----------------HANDSHAKE from clientID ", clientID, "---------------------")
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
				log.Println("pubk has been set: ", m.GetClient(clientID).PubKeyB)

				go func() {
					handshakeMsg := resmsg.HandshakeMsg{
						Type:             "HANDSHAKE-CONFIRMED",
						ServerPublicKeyB: m.PubKeyBytes(),
					}
					conn.WriteMessage(mt, handshakeMsg.Bytes())
				}()
				fmt.Println("------------------------------------------------------\n")
				// log.Println("---------Connection with clientID ", clientID, " established---------")

			case "SIG":
				log.Println("\n----------------Request from clientID ", clientID, "---------------------")

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

				sigFlag := m.VerifyRequestWithSig(clientID, req)
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
				conn.WriteMessage(mt, msg.Bytes())
				fmt.Println("------------------------------------------------------\n")

			case "TX":
				log.Println("\n----------------SignTx from clientID ", clientID, "---------------------")

				var req request.RequestMsg
				err := json.Unmarshal([]byte(body), &req)
				if err != nil {
					log.Fatal("Unmarshal error: ", err)
					break
				}
				jsonReq, _ := json.Marshal(req)
				log.Println("Request: ", string(jsonReq))

				requestByte := req.ReqByte
				// reqString := string(requestByte)
				log.Println(requestByte)

				// TODO: delete later

				// log.Println(hex.EncodeToString(req.SignedPaymentBody))
				// test := &request.PaymentBody{
				// 	ChannelID: req.ChannelID,
				// 	Amount:    req.Amount,
				// }
				// log.Println(hex.EncodeToString(test.HashByte()))

				// log.Println(len(req.SignedPaymentBody))
				// sig := req.SignedPaymentBody
				// v := sig[64]
				// if v == 0 || v == 1 {
				// 	v += 27
				// }

				// r := new(big.Int).SetBytes(sig[:32])
				// s := new(big.Int).SetBytes(sig[32:64])

				// fmt.Println("r: ", r)
				// fmt.Println("s: ", s)
				// fmt.Println("v: ", v)

				///////

				// verify the signature
				sigFlag := m.VerifyRequestWithSig(clientID, req)
				var msg resmsg.ServerMsg
				if sigFlag {
					log.Println("PASS: Signature verified")
					msg = resmsg.ServerMsg{
						Type: "info",
						Info: []byte("SigCheck: Passed"),
					}
				} else {
					msg = resmsg.ServerMsg{
						Type: "info",
						Info: []byte("SigCheck: WRONG signature"),
					}
				}
				conn.WriteMessage(mt, msg.Bytes())

				wsEndpoint := "ws://localhost:8100"
				client, err := ethclient.Dial(wsEndpoint)
				if err != nil {
					log.Fatal(err)
				}

				// forward to the network
				// -----
				// log.Println("Forward to the network")
				// connect to the network

				// // Send the tx to the network

				txSend := new(types.Transaction)
				rlp.DecodeBytes(requestByte, &txSend)

				err = client.SendTransaction(context.Background(), txSend)
				if err != nil {
					log.Fatal(err)
				}

				// by directly inject the transaction in the mempool
				// ethBackend, ok := backend.(*eth.EthAPIBackend)
				// if !ok {
				// 	utils.Fatalf("Ethereum service not running")
				// }
				// Set the gas price to the limits from the CLI and start mining

				// ethBackend.TxPool().SetGasTip(gasprice)
				// tx := new(types.Transaction)
				// rlp.DecodeBytes(requestByte, &tx)
				// if err := ethBackend.TxPool().AddLocal(tx); err != nil {
				// 	log.Fatal(err)
				// }
				// log.Println("tx submitted: %s", tx.Hash().Hex())

				log.Println("tx submitted: %s", txSend.Hash().Hex())
				msg = resmsg.ServerMsg{
					Type: "info-hex",
					Info: txSend.Hash().Bytes(),
				}
				conn.WriteMessage(mt, msg.Bytes())

				var txReceipt *types.Receipt
				for txReceipt == nil {
					// Query the transaction receipt
					txReceipt, err = client.TransactionReceipt(context.Background(), txSend.Hash())
					if err != nil {
						log.Println("Waiting for transaction to be mined...")
						time.Sleep(5 * time.Second) // Adjust the sleep duration based on expected block time
					}
				}
				log.Printf("Transaction mined in block %d", txReceipt.BlockNumber.Uint64())

				// By directly monitor the mempool

				// receiptBytes, err := json.Marshal(txReceipt)
				// if err != nil {
				// 	log.Fatal("Failed to marshal receipt: ", err)
				// }
				// _ = conn.WriteMessage(mt, []byte(receiptBytes))

				// Retrieve result
				const contractABI = `[{"inputs":[{"indexed":true,"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"ChannelOpened","type":"event"},{"inputs":[{"internalType":"address","name":"addr","type":"address"}],"name":"balance","outputs":[{"internalType":"uint256","name":"bal","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"closeChan","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"},{"internalType":"uint256","name":"value","type":"uint256"}],"name":"confirmClosure","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"}],"name":"greeting","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"senderDeposit","type":"uint256"}],"name":"openChan","outputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"paychanCheck","outputs":[{"components":[{"internalType":"bytes32","name":"id","type":"bytes32"},{"internalType":"address payable","name":"sender","type":"address"},{"internalType":"address payable","name":"recipient","type":"address"},{"internalType":"uint256","name":"senderDeposit","type":"uint256"},{"internalType":"uint256","name":"startTime","type":"uint256"},{"internalType":"uint256","name":"status","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"},{"internalType":"uint256","name":"disputeStartTime","type":"uint256"},{"internalType":"uint256","name":"disputeDuration","type":"uint256"},{"internalType":"bool","name":"senderConfirm","type":"bool"},{"internalType":"bool","name":"recipientConfirm","type":"bool"}],"internalType":"struct paychan.PayChan","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"channelId","type":"bytes32"}],"name":"paychanSelectedArguments","outputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"rec","type":"address"},{"internalType":"uint256","name":"status","type":"uint256"},{"internalType":"uint256","name":"senderB","type":"uint256"},{"internalType":"uint256","name":"fee","type":"uint256"}],"stateMutability":"view","type":"function"}]`
				contractAbi, err := abi.JSON(strings.NewReader(string(contractABI)))
				var channelId common.Hash
				for _, vLog := range txReceipt.Logs {
					fmt.Printf("Log Address: %s\n", vLog.Address.Hex())
					if len(vLog.Topics) > 0 {
						eventName, err := contractAbi.EventByID(vLog.Topics[0])
						if err != nil {
							log.Println("Error finding event name:", err)
							continue
						}
						fmt.Printf("Event Name: %s\n", eventName.Name)

						// var results []interface{}
						if len(vLog.Topics) > 1 {

							channelId = common.BytesToHash(vLog.Topics[1].Bytes())
							fmt.Printf("Channel ID: %s\n", channelId.Hex())
							msg = resmsg.ServerMsg{
								Type: "info-hex",
								Info: channelId.Bytes(),
							}
							conn.WriteMessage(mt, msg.Bytes())
						}
					}
				}

				// Generate the proof of the transaction
				blockHash := txReceipt.BlockHash
				block, _ := client.BlockByHash(context.Background(), blockHash)
				blockNr := block.Number()
				txHash := txReceipt.TxHash
				txIndex := txReceipt.TransactionIndex
				txRootHash := block.TxHash()

				log.Println("Block Hash: ", blockHash.Hex())
				log.Println("Tx Hash: ", txHash.Hex())
				log.Println("Tx Index: ", txIndex)
				log.Println("TxRoot Hash: ", txRootHash)

				proof, idx := generateProof(block, txHash)
				if err != nil {
					log.Println("Error: ", err)
				}

				if proof == nil {
					log.Println("Error: unable to generate the proof")
				} else {
					res := verifyProof(txHash, proof, blockNr, uint32(idx))

					log.Println("Proof Verification: ", res)
				}
				// serializedProof := proof.Serialize()
				// var buffer bytes.Buffer
				// for _, part := range serializedProof {
				// 	buffer.Write(part)
				// }

				// Send the response to the client
				responseBody := resmsg.ResponseBody{
					SignedReqBody: req.SignedReqBody,
					Proof:         proof.CustomSerialize(),
					TxHash:        txHash,
					TxIdx:         uint32(idx),
				}
				log.Println(proof)
				sig := cryptoutil.Sign(m.PrivateKey, responseBody.HashBytes())
				responseMsg := resmsg.ResponseMsg{
					Type:               "response",
					ChannelId:          channelId,
					Amount:             req.Amount,
					SignedReqBody:      req.SignedReqBody,
					CurrentBlockHeight: blockNr,
					ReturnValue:        txReceipt.Bloom.Bytes(),
					Proof:              proof.CustomSerialize(),
					TxHash:             txHash,
					TxIdx:              uint32(idx),
					Signature:          sig,
				}

				log.Println(responseMsg)
				_ = conn.WriteMessage(mt, responseMsg.Bytes())

				fmt.Println("------------------------------------------------------\n")
			default:
				log.Println("\n----------------Message from clientID ", clientID, "---------------------")

				log.Printf("Recv from clientID %s: %s", clientID, msg)
				serverMsg := resmsg.ServerMsg{
					Type: "info",
					Info: []byte("server:" + string(msg)),
				}
				err = conn.WriteMessage(mt, serverMsg.Bytes())
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

func generateProof(block *types.Block, txHash common.Hash) (mpt.Proof, int) {
	txs := block.Transactions()
	idx := -1
	for index, tx := range txs {
		if tx.Hash() == txHash {
			idx = index
		}
	}
	if idx < 0 {
		return nil, -1
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
	return proof, idx
}

func verifyProof(txHash common.Hash, proof mpt.Proof, blockNr *big.Int, idx uint32) bool {
	wsEndpoint := "ws://localhost:8100"
	client, err := ethclient.Dial(wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	// query the block information
	block, _ := client.HeaderByNumber(context.Background(), blockNr)
	txRootHash := block.TxHash
	tx, _, _ := client.TransactionByHash(context.Background(), txHash)
	txRLP, _ := rlp.EncodeToBytes(tx)
	key, _ := rlp.EncodeToBytes(uint32(idx))
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
