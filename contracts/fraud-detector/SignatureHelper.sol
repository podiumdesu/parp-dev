// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

library SignatureHelper {
    function recover(bytes32 hash, bytes memory signature) external pure returns (address) {
        return ECDSA.recover(hash, signature);
    }
}