// SPDX-License-Identifier: Apache2
pragma solidity ^0.8.17;

import "../types/proof.sol";
import "../types/events.sol";

import { MerkleVerify as MerkleVerify} from "./merkleverify.sol";

library FraudProofHelper {

    function verifyProof(
        bytes32 root,
        bytes[] memory proof,
        bytes[] memory keys
    ) public returns (bool) {
        // verifies ethereum specific merkle patricia proofs as described by EIP-1188.
        // can be used to verify the receipt trie, transaction trie and state trie
        // contributed by @ripa1995
        (bool success, StorageValue[] memory values) = MerkleVerify.VerifyEthereumProof(root, proof, keys);
        // do something with the verified values.
        // Emit the event to log the values
        if (success) {
            emit emitProofValues(values);
            return true;
        } else {
            return false;
        }
    }



}
