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

func FromEthTransaction(t *types.Transaction) *mpt.Transaction {
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

func FuncTestTransactionRootAndProof() (mpt.Proof, []byte, common.Hash) {

	trie := mpt.NewTrie()

	txs := TransactionsJSONRead()

	log.Println("Tx: ", txs)

	for i, tx := range txs {
		// key is the encoding of the index as the unsigned integer type
		key, _ := rlp.EncodeToBytes(uint(i))
		// log.Println("key, value: ", i)
		// log.Println(key)
		// require.NoError(t, err)
		transaction := FromEthTransaction(tx)

		// value is the RLP encoding of a transaction
		rlp, _ := transaction.GetRLP()
		// log.Println(rlp)
		// require.NoError(t, err)

		trie.Put(key, rlp)
	}

	txHash := txs[30].Hash()
	log.Println("txHash: ", txHash)
	// the transaction root for block 10467135
	// https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0x9fb73f&boolean=true&apikey=YourApiKeyToken
	transactionRoot, _ := hex.DecodeString("bb345e208bda953c908027a45aa443d6cab6b8d2fd64e83ec52f1008ddeafa58")
	// require.NoError(t, err)

	// t.Run("merkle root hash should match with 10467135's transactionRoot", func(t *testing.T) {
	// transaction root should match with block 10467135's transactionRoot
	// require.Equal(t, transactionRoot, trie.Hash())
	// })

	// t.Run("a merkle proof for a certain transaction can be verified by the offical trie implementation", func(t *testing.T) {
	key, _ := rlp.EncodeToBytes(uint(30))
	// require.NoError(t, err)

	proof, _ := trie.Prove(key)
	// require.Equal(t, true, found)

	// txRLP, err := VerifyProof(transactionRoot, key, proof)
	// require.NoError(t, err)

	// // verify that if the verification passes, it returns the RLP encoded transaction
	// rlp, err := FromEthTransaction(txs[30]).GetRLP()
	// require.NoError(t, err)
	// require.Equal(t, rlp, txRLP)

	log.Println("______-------________")
	log.Println(proof)
	log.Println("______-------________")

	serializedProof := proof.CustomRLPSerialize()
	deserialProof, _ := mpt.DeserializeProof(serializedProof)
	log.Println(deserialProof)
	log.Println("______-------________dafs")
	serializedProofStr := proof.CustomRLPSerializeString()

	// Use to check the proof
	txRLP, _ := mpt.VerifyProof(transactionRoot, key, proof)
	// require.NoError(t, err)
	// rlp, err := FromEthTransaction(txs[30]).GetRLP()
	// require.NoError(t, err)
	// require.Equal(t, rlp, txRLP)
	// ----end----

	fmt.Println("txRLP: ", hex.EncodeToString(txRLP))
	fmt.Printf("Root Hash: %x\n", transactionRoot)

	fmt.Println(serializedProofStr)
	// 2. Print the serialized proof as a bytes array (proof is []string).
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
	// 3. Print the keys as a bytes array.
	// RLP encode the index as the key (assuming resMsg.TxIdx is the index)
	// Convert the RLP-encoded key to hex format.
	keyHex := "0x" + hex.EncodeToString(key)
	// ["0xf90131a0a07c378dc14ba3b7c4dbd27b00f3512d8e211103daa335f308c3420c1aeba3caa057625d7b5df6fe76d0420dfe57514fe347b642144d1531f8514495d530a85253a0559174a95832abf4370c3ade6c8167230dd4f27c4e8a01eaeb72f260d409297ea0fc032970cced211d96d7ad41fc28ac7d0d7f6bb08adcb7f489e73a4c308e64e8a094d6e64babce6d4ad7fa22cc17b64256826314879b9918073456a4cdb4a3f878a063fb3460b72a87f24e0e60c52f11d78ee2abe43e1a355db70415d2d80dea6e25a09f4bcadeebaabd049060cbedbe3a347b2ea5e3216a4421e7c0fdb6cb1f5bf8afa0bbd77f432dcf12f1253c3d23bbcfeefca41b1a73a50e5a0c4afab1033bf3ad7ba09fc760313e0298268573d4ff1934cd26769811a328605a83f406909ec87386808080808080808080", "0xf90211a02a5225c0003862cbcd182f6e939c42c10fc3318de84e9fc880dd275214368ef2a0367346a5453f846c40ff20a58b24d97506149656fe664138db6a5ab1b4f0135ea02c48bce71b884571810c2d9ec68cb7d11308d115e90570744c8ee8c803e1857ba0fa237768f28ab0446aa1dce8eb74545e646fdb949d3d0c4f7357b6e96c1c672aa05d8f7ef95e9b1080aa538133243643c9b49246912dbbca1ded882884d012ce7da01678b089698e01ac22535302d9e95f0716538f45a459a4a56f3942fc6fd0dcbfa0b625cd08dad5b5a151eb2d04da332a25c1e018f57bf903f3cd148ec1cc6b6775a0656207acf984bc2fb4cd2c750aa3179b3432a196cc6d8bf110206fa731be49dea00d16ee630b1a360b69e74f6065387464b66adf191efc81eda1544d85cba7754aa07e89f99a1784dce5045f0d78a1da77ceca74e532fd1f94aef0c499bf72498ea9a05fbbe68419da961d263664408b02305a819df9c1efc56da9abcb0b4e45957676a06b092a40e835e5149adf07049e3fca79e562597b4e3b6a553f4217dc82cdad1da05060297997665ff3ab669e54fe73b672511a894ba0880664471ca99282b8f379a02d59e2921f0f7210a9732f2a2ed1cd445650cfb23dc6465a34f2243128367c2ea043b980a6d0bc5ed554796e169500ee8160feaf937ba8eb5ad251db89ff33dc92a05a490e876b754b9d675f70a0a524d2ff399aa079cbece960ccef2c512062799980", "0xf87420b871f86f84019a41db843b9aca0082c35094783745372c230805071b81305eed750489d0ea7f8801723ed8a9cf24dc8026a04041df1015a8e8ed240a57f4dde96a0614ec5073484ad634e25b8cd7ac9b92f5a02f518cc656e4effb425b314c389172b9ccc01998b0c398d73a43e0ffaa6bf094"]
	// Print the key as bytes[] for Remix.
	fmt.Printf("Keys (bytes[] keys): [\"%s\"]\n", keyHex)

	// })

	return proof, key, txHash
}
