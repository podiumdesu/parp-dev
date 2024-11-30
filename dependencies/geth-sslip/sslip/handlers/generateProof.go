package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/sslip/mpt"
)

func TransactionsJSONRead() []*types.Transaction {
	path, err := filepath.Abs("./transactions.json")
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file at %s: %v", path, err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var txs []*types.Transaction
	json.Unmarshal(byteValue, &txs)
	return txs
}

func FuncTestTransactionRootAndProof() (mpt.Proof, []byte, common.Hash, *types.Transaction) {

	trie := mpt.NewTrie()
	txs := TransactionsJSONRead()

	for i, tx := range txs {
		// key is the encoding of the index as the unsigned integer type
		key, _ := rlp.EncodeToBytes(uint(i))
		transaction := mpt.FromEthTransaction(tx)

		// value is the RLP encoding of a transaction
		rlp, _ := transaction.GetRLP()
		trie.Put(key, rlp)
	}

	txHash := txs[30].Hash()

	// the transaction root for block 10467135
	// https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0x9fb73f&boolean=true&apikey=YourApiKeyToken

	key, _ := rlp.EncodeToBytes(uint(30))

	proof, _ := trie.Prove(key)

	// Check Serialization and Deserialization
	// serializedProof := proof.CustomRLPSerialize()
	// deserialProof, _ := mpt.DeserializeProof(serializedProof)

	return proof, key, txHash, txs[30]
}

func PrintProofToSubmit(proof mpt.Proof, key []byte, rootHash common.Hash) {

	log.Println("root Hash: ", rootHash)
	// 2. Print the serialized proof as a bytes array (proof is []string).
	serializedProofStr := proof.CustomRLPSerializeString()
	fmt.Print("Proof (bytes[] proof): [")
	for i, node := range serializedProofStr {
		if i > 0 {
			fmt.Print(", ")
		}
		// Add "0x" prefix only once for each proof element.
		if len(node) >= 2 && node[:2] == "0x" {
			// If the proof element already has "0x" prefix, use it directly.
			fmt.Printf("\"%s\"", node)
		} else {
			// Otherwise, add "0x" prefix.
			fmt.Printf("\"0x%s\"", node)
		}
	}
	fmt.Println("]")

	keyHex := "0x" + hex.EncodeToString(key)
	// Print the key as bytes[] for Remix.
	fmt.Printf("Keys (bytes[] keys): [\"%s\"]\n", keyHex)
}
