// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract HeaderDecoder {
    struct BasicHeader {
        bytes32 parentHash;
        bytes32 uncleHash;
        address coinbase;
        bytes32 stateRoot;
        bytes32 txRoot;
        bytes32 receiptRoot;
    }

    struct ExtendedHeader {
        bytes logsBloom;
        uint256 difficulty;
        uint256 number;
        uint256 gasLimit;
        uint256 gasUsed;
        uint256 timestamp;
        bytes extraData;
        bytes32 mixDigest;
        uint64 nonce;
    }

    struct OptionalFields {
        uint256 baseFee;
        bytes32 withdrawalsRoot;
        uint256 blobGasUsed;
        uint256 excessBlobGas;
        bytes32 parentBeaconRoot;
    }

    function generateBlockHash(
        BasicHeader memory basic,
        ExtendedHeader memory extended,
        OptionalFields memory optional
    ) public pure returns (bytes32) {
        // Encode each part
        bytes memory encodedBasic = abi.encodePacked(
            encodeBytes32(basic.parentHash),
            encodeBytes32(basic.uncleHash),
            encodeAddress(basic.coinbase),
            encodeBytes32(basic.stateRoot),
            encodeBytes32(basic.txRoot),
            encodeBytes32(basic.receiptRoot)
        );

        bytes memory encodedExtended = abi.encodePacked(
            encodeBytes(extended.logsBloom),
            encodeUint(extended.difficulty),
            encodeUint(extended.number),
            encodeUint(extended.gasLimit),
            encodeUint(extended.gasUsed),
            encodeUint(extended.timestamp),
            encodeBytes(extended.extraData),
            encodeBytes32(extended.mixDigest),
            encodeUint(extended.nonce)
        );

        bytes memory encodedOptional = abi.encodePacked(
            encodeUint(optional.baseFee),
            encodeBytes32(optional.withdrawalsRoot),
            encodeUint(optional.blobGasUsed),
            encodeUint(optional.excessBlobGas),
            encodeBytes32(optional.parentBeaconRoot)
        );

        // Concatenate all encoded parts
        bytes memory rlpEncoded = abi.encodePacked(
            encodedBasic,
            encodedExtended,
            encodedOptional
        );

        // Compute and return the Keccak256 hash
        return keccak256(rlpEncoded);
    }

    // Helper function to encode a uint256 in RLP
    function encodeUint(uint256 value) internal pure returns (bytes memory) {
        if (value == 0) {
            return hex"80"; // RLP for zero
        } else if (value < 0x80) {
            return abi.encodePacked(uint8(value)); // Single-byte value
        } else {
            bytes memory encoded = toBytes(value);
            return abi.encodePacked(uint8(0x80 + encoded.length), encoded);
        }
    }

    // Helper function to encode a bytes array in RLP
    function encodeBytes(bytes memory data) internal pure returns (bytes memory) {
        if (data.length == 0) {
            return hex"80"; // RLP for empty bytes
        } else if (data.length == 1 && uint8(data[0]) < 0x80) {
            return data; // Single-byte value
        } else {
            return abi.encodePacked(encodeLength(data.length, 0x80), data);
        }
    }

    // Helper function to encode a bytes32 in RLP
    function encodeBytes32(bytes32 data) internal pure returns (bytes memory) {
        return encodeBytes(abi.encodePacked(data));
    }

    // Helper function to encode an address in RLP
    function encodeAddress(address data) internal pure returns (bytes memory) {
        return encodeBytes(abi.encodePacked(data));
    }

    // Encode the length of the data
    function encodeLength(uint256 length, uint256 offset) internal pure returns (bytes memory) {
        if (length < 56) {
            return abi.encodePacked(uint8(length + offset));
        } else {
            bytes memory encoded = toBytes(length);
            return abi.encodePacked(uint8(encoded.length + offset + 55), encoded);
        }
    }

    // Convert uint256 to bytes
    function toBytes(uint256 value) internal pure returns (bytes memory) {
        bytes memory result = new bytes(32);
        assembly {
            mstore(add(result, 32), value)
        }
        uint256 offset = 0;
        while (offset < 32 && result[offset] == 0) {
            offset++;
        }
        bytes memory trimmed = new bytes(32 - offset);
        for (uint256 i = offset; i < 32; i++) {
            trimmed[i - offset] = result[i];
        }
        return trimmed;
    }
}

function _generateHeaderHash(
    BasicHeader memory basic,
    ExtendedHeader memory extended,
    OptionalFields memory optional
) internal pure returns (bytes32) {
    // Manually encode each field following RLP rules

    // Encode basic fields
    bytes memory rlpBasic = abi.encodePacked(
        encodeBytes32(basic.parentHash),
        encodeBytes32(basic.uncleHash),
        encodeAddress(basic.coinbase),
        encodeBytes32(basic.stateRoot),
        encodeBytes32(basic.txRoot),
        encodeBytes32(basic.receiptRoot)
    );

    // Encode extended fields
    bytes memory rlpExtended = abi.encodePacked(
        encodeBytes(extended.logsBloom),
        encodeUint(extended.difficulty),
        encodeUint(extended.number),
        encodeUint(extended.gasLimit),
        encodeUint(extended.gasUsed),
        encodeUint(extended.timestamp),
        encodeBytes(extended.extraData),
        encodeBytes32(extended.mixDigest),
        encodeUint(extended.nonce)
    );

    // Encode optional fields
    bytes memory rlpOptional = abi.encodePacked(
        optional.baseFee > 0 ? encodeUint(optional.baseFee) : encodeEmpty(),
        optional.withdrawalsRoot != bytes32(0) ? encodeBytes32(optional.withdrawalsRoot) : encodeEmpty(),
        optional.blobGasUsed > 0 ? encodeUint(optional.blobGasUsed) : encodeEmpty(),
        optional.excessBlobGas > 0 ? encodeUint(optional.excessBlobGas) : encodeEmpty(),
        optional.parentBeaconRoot != bytes32(0) ? encodeBytes32(optional.parentBeaconRoot) : encodeEmpty()
    );

    // Concatenate all encoded parts into an RLP list
    bytes memory rlpEncoded = abi.encodePacked(rlpBasic, rlpExtended, rlpOptional);

    // Compute the Keccak256 hash of the RLP-encoded header
    return keccak256(rlpEncoded);
}

function encodeUint(uint256 value) internal pure returns (bytes memory) {
    if (value == 0) {
        return hex"80"; // RLP encoding for zero
    } else if (value < 0x80) {
        return abi.encodePacked(uint8(value)); // Single byte encoding
    } else {
        bytes memory encoded = toBytes(value);
        return abi.encodePacked(uint8(0x80 + encoded.length), encoded); // Multibyte encoding
    }
}

function encodeBytes(bytes memory data) internal pure returns (bytes memory) {
    if (data.length == 0) {
        return hex"80"; // RLP encoding for empty bytes
    } else if (data.length == 1 && uint8(data[0]) < 0x80) {
        return data; // Single byte optimization
    } else {
        return abi.encodePacked(encodeLength(data.length, 0x80), data);
    }
}

function encodeBytes32(bytes32 data) internal pure returns (bytes memory) {
    return encodeBytes(abi.encodePacked(data));
}

function encodeAddress(address data) internal pure returns (bytes memory) {
    return encodeBytes(abi.encodePacked(data));
}

function encodeEmpty() internal pure returns (bytes memory) {
    return hex"80"; // RLP encoding for empty fields
}

function encodeLength(uint256 length, uint256 offset) internal pure returns (bytes memory) {
    if (length < 56) {
        return abi.encodePacked(uint8(length + offset));
    } else {
        bytes memory encodedLength = toBytes(length);
        return abi.encodePacked(uint8(encodedLength.length + offset + 55), encodedLength);
    }
}

function toBytes(uint256 value) internal pure returns (bytes memory) {
    bytes memory result = new bytes(32);
    assembly {
        mstore(add(result, 32), value)
    }
    uint256 offset = 0;
    while (offset < 32 && result[offset] == 0) {
        offset++;
    }
    bytes memory trimmed = new bytes(32 - offset);
    for (uint256 i = offset; i < 32; i++) {
        trimmed[i - offset] = result[i];
    }
    return trimmed;
}
