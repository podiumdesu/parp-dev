pragma solidity ^0.8.17;

import "./trie/Node.sol";
import "./trie/Option.sol";
import "./trie/NibbleSlice.sol";
import "./trie/TrieDB.sol";

import "./trie/EthereumTrieDB.sol";
import "./Types.sol";

// SPDX-License-Identifier: Apache2

/**
 * @title A Merkle Patricia library
 * @author Polytope Labs
 * @dev Use this library to verify merkle patricia proofs
 * @dev refer to research for more info. https://research.polytope.technology/state-(machine)-proofs
 */
library MerkleVerify{

    // /**
    //  * @notice Verifies ethereum specific merkle patricia proofs as described by EIP-1188.
    //  * @param root hash of the merkle patricia trie
    //  * @param proof a list of proof nodes
    //  * @param keys a list of keys to verify
    //  * @return bytes[] a list of values corresponding to the supplied keys.
    //  */
    function VerifyEthereumProof(
        bytes32 root,
        bytes[] memory proof,
        bytes[] memory keys
    ) public returns (bool success, StorageValue[] memory) {
        StorageValue[] memory values = new StorageValue[](keys.length);
        TrieNode[] memory nodes = new TrieNode[](proof.length);

        for (uint256 i = 0; i < proof.length; i++) {
            nodes[i] = TrieNode(keccak256(proof[i]), proof[i]);
        }

        for (uint256 i = 0; i < keys.length; i++) {
            values[i].key = keys[i];
            NibbleSlice memory keyNibbles = NibbleSlice(keys[i], 0);
            (bool found, bytes memory nodeData) = TrieDB.get(nodes, root);
            if (!found) {
            return (false, values); // If a required node is missing, return failure
            }

            NodeKind memory node = EthereumTrieDB.decodeNodeKind(nodeData);

            // This loop is unbounded so that an adversary cannot insert a deeply nested key in the trie
            // and successfully convince us of it's non-existence, if we consume the block gas limit while
            // traversing the trie, then the transaction should revert.
            for (uint256 j = 1; j > 0; j++) {
                NodeHandle memory nextNode;

                if (TrieDB.isLeaf(node)) {
                    Leaf memory leaf = EthereumTrieDB.decodeLeaf(node);
                    // Let's retrieve the offset to be used
                    uint256 offset = keyNibbles.offset % 2 == 0
                        ? keyNibbles.offset / 2
                        : keyNibbles.offset / 2 + 1;
                    // Let's cut the key passed as input
                    keyNibbles = NibbleSlice(
                        NibbleSliceOps.bytesSlice(keyNibbles.data, offset),
                        0
                    );
                    if (NibbleSliceOps.eq(leaf.key, keyNibbles)) {
                        (found, values[i].value) = TrieDB.load(nodes, leaf.value);
                        if (!found) {
                        return (false, values); // Failure on loading leaf value
                        }
                    }
                    break;
                } else if (TrieDB.isExtension(node)) {
                    Extension memory extension = EthereumTrieDB.decodeExtension(
                        node
                    );
                    if (NibbleSliceOps.startsWith(keyNibbles, extension.key)) {
                        // Let's cut the key passed as input
                        uint256 cutNibble = keyNibbles.offset +
                            NibbleSliceOps.len(extension.key);
                        keyNibbles = NibbleSlice(
                            NibbleSliceOps.bytesSlice(
                                keyNibbles.data,
                                cutNibble / 2
                            ),
                            cutNibble % 2
                        );
                        nextNode = extension.node;
                    } else {
                        break;
                    }
                } else if (TrieDB.isBranch(node)) {
                    Branch memory branch = EthereumTrieDB.decodeBranch(node);
                    if (NibbleSliceOps.isEmpty(keyNibbles)) {
                        if (Option.isSome(branch.value)) {
                            (found, values[i].value) = TrieDB.load(nodes, branch.value.value);
                            if (!found) {
                            return (false, values); // Failure on loading leaf value
                            }
                        }
                        break;
                    } else {
                        NodeHandleOption memory handle = branch.children[
                            NibbleSliceOps.at(keyNibbles, 0)
                        ];
                        if (Option.isSome(handle)) {
                            keyNibbles = NibbleSliceOps.mid(keyNibbles, 1);
                            nextNode = handle.value;
                        } else {
                            break;
                        }
                    }
                } else if (TrieDB.isEmpty(node)) {
                    break;
                }

                bytes memory nextNodeData;
                (found, nextNodeData) = TrieDB.load(nodes, nextNode);
                if (!found) {
                    return (false, values); // Failure on loading leaf value
                }
                node = EthereumTrieDB.decodeNodeKind(
                    nextNodeData
                );
            }
        }

        return (true, values);
    }
}