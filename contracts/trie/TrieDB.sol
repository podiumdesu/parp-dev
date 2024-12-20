// SPDX-License-Identifier: Apache2
pragma solidity ^0.8.17;

import "./Node.sol";

event DebugHash(bytes32 hash);

library TrieDB {
    function get(
        TrieNode[] memory nodes,
        bytes32 hash
    ) internal returns (bool found, bytes memory) {
        for (uint256 i = 0; i < nodes.length; i++) {
            emit DebugHash(hash);

            if (nodes[i].hash == hash) {
                return (true, nodes[i].node);
            }
        }
        // revert("Incomplete Proof!");
        return (false, bytes(""));
    }

    function load(
        TrieNode[] memory nodes,
        NodeHandle memory node
    ) internal returns (bool found, bytes memory) {
        if (node.isInline) {
            return (true, node.inLine);
        } else if (node.isHash) {
            return get(nodes, node.hash);
        }

        return (true, bytes(""));
    }

    function isNibbledBranch(
        NodeKind memory node
    ) internal pure returns (bool) {
        return (node.isNibbledBranch ||
            node.isNibbledHashedValueBranch ||
            node.isNibbledValueBranch);
    }

    function isExtension(NodeKind memory node) internal pure returns (bool) {
        return node.isExtension;
    }

    function isBranch(NodeKind memory node) internal pure returns (bool) {
        return node.isBranch;
    }

    function isLeaf(NodeKind memory node) internal pure returns (bool) {
        return (node.isLeaf || node.isHashedLeaf);
    }

    function isEmpty(NodeKind memory node) internal pure returns (bool) {
        return node.isEmpty;
    }

    function isHash(NodeHandle memory node) internal pure returns (bool) {
        return node.isHash;
    }

    function isInline(NodeHandle memory node) internal pure returns (bool) {
        return node.isInline;
    }
}