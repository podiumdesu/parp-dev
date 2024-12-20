// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract SigHelper {
    function recoverSignature(
        bytes32 messageHash,
        bytes memory signature
    ) external pure returns (address) {
        require(signature.length == 65, "invalid signature length");

        bytes32 r;
        bytes32 s;
        uint8 v;

        // Extract r, s and v from the signature
        assembly {
            r := mload(add(signature, 0x20))
            s := mload(add(signature, 0x40))
            v := byte(0, mload(add(signature, 0x60)))
        }
        // Adjust the v value if necessary
        if (v < 27) {
            v += 27;
        }

        // Ensure v is valid
        require(v == 27 || v == 28, "Invalid v value");

        // Use ecrecover to recover the address from the signature
        address signer = ecrecover(messageHash, v, r, s);

        // Check if the recovered signer is the expected signer
        return signer;
    }

}