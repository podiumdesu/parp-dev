package mpt

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"os"

// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/rlp"
// )

// func main() {
// 	// TestTransactionRootAndProof()
// 	// Read Transactions
// 	trieWithOneTx()
// 	trieWithBlockTxs()

// }

// func checkOneTransaction() *types.Transaction {
// 	jsonFile, _ := os.Open("transaction.json")
// 	defer jsonFile.Close()
// 	byteValue, _ := ioutil.ReadAll(jsonFile)
// 	var tx types.Transaction
// 	json.Unmarshal(byteValue, &tx)
// 	return &tx
// }
// func fromEthTransaction(t *types.Transaction) *Transaction {
// 	v, r, s := t.RawSignatureValues()
// 	return &Transaction{
// 		AccountNonce: t.Nonce(),
// 		Price:        t.GasPrice(),
// 		GasLimit:     t.Gas(),
// 		Recipient:    t.To(),
// 		Amount:       t.Value(),
// 		Payload:      t.Data(),
// 		V:            v,
// 		R:            r,
// 		S:            s,
// 	}
// }

// func trieWithOneTx() {
// 	key, _ := rlp.EncodeToBytes(uint(0))

// 	tx := checkOneTransaction()

// 	transaction := fromEthTransaction(tx)
// 	rlp, _ := transaction.GetRLP()

// 	trie := NewTrie()
// 	trie.Put(key, rlp)

// 	txRootHash := fmt.Sprintf("%x", types.DeriveSha(types.Transactions{tx}))
// 	fmt.Printf("txRootHash: %v\n", txRootHash)
// 	fmt.Printf("root: %x\n", trie.Hash())
// 	// require.Equal(t, txRootHash, fmt.Sprintf("%x", trie.Hash()))
// }

// func transactionsJSON() []*types.Transaction {
// 	jsonFile, _ := os.Open("transactions.json")
// 	defer jsonFile.Close()
// 	byteValue, _ := ioutil.ReadAll(jsonFile)
// 	var txs []*types.Transaction
// 	json.Unmarshal(byteValue, &txs)
// 	return txs
// }

// func trieWithBlockTxs() {

// 	txs := transactionsJSON()

// 	trie := NewTrie()
// 	for i, tx := range txs {
// 		key, _ := rlp.EncodeToBytes(uint(i))

// 		transaction := fromEthTransaction(tx)

// 		rlp, _ := transaction.GetRLP()

// 		trie.Put(key, rlp)
// 	}

// 	txRootHash := fmt.Sprintf("%x", types.DeriveSha(types.Transactions(txs)))
// 	fmt.Printf("txRootHash: %v\n", txRootHash)
// 	fmt.Printf("%x", trie.Hash())

// 	// generate the proof and verify it
// 	key, _ := rlp.EncodeToBytes(uint(23))
// 	proof, found := trie.Prove(key)

// 	fmt.Printf("proof: %x, found: %v\n", proof, found)

// 	txRLP, _ := VerifyProof(trie.Hash(), key, proof)
// 	rlp, _ := fromEthTransaction(txs[30]).GetRLP()

// 	fmt.Println(txRLP)
// 	fmt.Println(rlp)
// }
