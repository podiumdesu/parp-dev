package mpt

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
)

type Proof interface {
	// Put inserts the given value into the key-value data store.
	Put(key []byte, value []byte) error

	// Delete removes the key from the key-value data store.
	Delete(key []byte) error

	// Has retrieves if a key is present in the key-value data store.
	Has(key []byte) (bool, error)

	// Get retrieves the given key if it's present in the key-value data store.
	Get(key []byte) ([]byte, error)

	// Serialize returns the serialized proof
	Serialize() [][]byte

	// DeserializeProof(proof [][]byte) *ProofDB

	CustomSerialize() [][]byte

	CustomRLPSerialize() [][]byte

	CustomRLPSerializeString() []string

	CustomRLPSerializeBytes() []interface{}
}

type ProofDB struct {
	kv map[string][]byte
}

func NewProofDB() *ProofDB {
	return &ProofDB{
		kv: make(map[string][]byte),
	}
}

func (w *ProofDB) Put(key []byte, value []byte) error {
	keyS := fmt.Sprintf("%x", key)
	w.kv[keyS] = value
	fmt.Printf("put key: %x, value: %x\n", key, value)
	return nil
}

func (w *ProofDB) Delete(key []byte) error {
	keyS := fmt.Sprintf("%x", key)
	delete(w.kv, keyS)
	return nil
}
func (w *ProofDB) Has(key []byte) (bool, error) {
	keyS := fmt.Sprintf("%x", key)
	_, ok := w.kv[keyS]
	return ok, nil
}

func (w *ProofDB) Get(key []byte) ([]byte, error) {
	keyS := fmt.Sprintf("%x", key)
	val, ok := w.kv[keyS]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return val, nil
}

func (w *ProofDB) Serialize() [][]byte {
	nodes := make([][]byte, 0, len(w.kv))
	for _, value := range w.kv {
		nodes = append(nodes, value)
	}
	return nodes
}

func (w *ProofDB) CustomSerialize() [][]byte {
	nodes := make([][]byte, 0, len(w.kv)*2)
	for key, value := range w.kv {
		keyBytes, _ := hex.DecodeString(key)
		nodes = append(nodes, keyBytes, value)
	}
	return nodes
}

// func DeserializeProof(data [][]byte) (*ProofDB, error) {

// }

func (w *ProofDB) CustomRLPSerialize() [][]byte {
	var serializedProof [][]byte
	for _, value := range w.kv {
		// decodedKey, _ := hex.DecodeString(key)
		// serializedProof = append(serializedProof, decodedKey, value)
		serializedProof = append(serializedProof, value)
	}
	return serializedProof
}

func DeserializeProof(data [][]byte) (*ProofDB, error) {
	proofDB := NewProofDB()
	log.Println("len(data): ", len(data))
	for i := 0; i < len(data); i += 1 {
		// if i+1 >= len(data) {
		// 	return nil, fmt.Errorf("odd number of elements in serialized data, expected key-value pairs")
		// }
		// key := data[i]
		key := Keccak256(data[i])
		value := data[i]
		log.Println("key: ", key, "value: ", value)
		if err := proofDB.Put([]byte(key), value); err != nil {
			return nil, fmt.Errorf("failed to put key-value pair in proof DB: %v", err)
		}
	}
	return proofDB, nil
}

// the old one for []string
func (w *ProofDB) CustomRLPSerializeString() []string {
	var serializedProof []string

	// Iterate over each TrieNode and convert its serialized content back to hex string
	for _, node := range w.kv {
		// Convert the serialized byte slice to a hex string
		serializedHex := hex.EncodeToString(node)

		// fmt.Println("serializedHex: ", serializedHex)
		// Add the "0x" prefix to indicate hexadecimal format
		serializedProof = append(serializedProof, "0x"+serializedHex)
	}

	return serializedProof
}

func (w *ProofDB) CustomRLPSerializeBytes() []interface{} {
	var serializedProof [][]byte

	for _, node := range w.kv {
		serializedProof = append(serializedProof, node)
	}

	// return serializedProof

	var proofAsInterface []interface{}
	for _, proofNode := range serializedProof {
		proofAsInterface = append(proofAsInterface, proofNode)
	}

	return proofAsInterface
}

// func SerializeProofNodes(nodes [][]byte) ([]string, error) {
// 	var serializedProof []string

// 	// Iterate over each TrieNode and convert its serialized content back to hex string
// 	for _, node := range nodes {
// 		// Convert the serialized byte slice to a hex string
// 		serializedHex := hex.EncodeToString(node)

// 		// Add the "0x" prefix to indicate hexadecimal format
// 		serializedProof = append(serializedProof, "0x"+serializedHex)
// 	}

// 	return serializedProof, nil
// }

func DeserializeProofNodes(serializedProof []string) (*ProofDB, error) {
	proofDB := NewProofDB()

	// Iterate over each serialized hex string in the proof array
	for _, hexNode := range serializedProof {
		// Remove the "0x" prefix if present
		if len(hexNode) > 1 && hexNode[:2] == "0x" {
			hexNode = hexNode[2:]
		}

		// Decode the hex string into bytes
		serializedBytes, err := hex.DecodeString(hexNode)
		if err != nil {
			return nil, fmt.Errorf("failed to decode hex string: %w", err)
		}

		// Calculate the keccak256 hash of the serialized node bytes (same as Solidity's keccak256)
		nodeHash := crypto.Keccak256(serializedBytes)

		// Create the TrieNode and store the hash and serialized content
		proofDB.Put(nodeHash, serializedBytes)

	}

	return proofDB, nil
}

// Prove returns the merkle proof for the given key, which is
func (t *Trie) Prove(key []byte) (Proof, bool) {
	proof := NewProofDB()
	node := t.root
	nibbles := FromBytes(key)

	for {
		proof.Put(Hash(node), Serialize(node))

		if IsEmptyNode(node) {
			return nil, false
		}

		if leaf, ok := node.(*LeafNode); ok {
			matched := PrefixMatchedLen(leaf.Path, nibbles)
			if matched != len(leaf.Path) || matched != len(nibbles) {
				return nil, false
			}

			return proof, true
		}

		if branch, ok := node.(*BranchNode); ok {
			if len(nibbles) == 0 {
				return proof, branch.HasValue()
			}

			b, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining
			node = branch.Branches[b]
			continue
		}

		if ext, ok := node.(*ExtensionNode); ok {
			matched := PrefixMatchedLen(ext.Path, nibbles)
			// E 01020304
			//   010203
			if matched < len(ext.Path) {
				return nil, false
			}

			nibbles = nibbles[matched:]
			node = ext.Next
			continue
		}

		panic("not found")
	}
}

// VerifyProof verify the proof for the given key under the given root hash using go-ethereum's VerifyProof implementation.
// It returns the value for the key if the proof is valid, otherwise error will be returned
func VerifyProof(rootHash []byte, key []byte, proof Proof) (value []byte, err error) {
	return trie.VerifyProof(common.BytesToHash(rootHash), key, proof)
}

// func VerifySerializedProof(rootHash []byte, key []byte, proof []string) (value []byte, err error) {

// 	deserializedProofNodes, _ := DeserializeProofNodes(proof)

// 	return trie.VerifyProof(common.BytesToHash(rootHash), key, deserializedProofNodes)
// }
