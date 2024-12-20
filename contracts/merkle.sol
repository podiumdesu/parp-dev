// // SPDX-License-Identifier: GPL-3.0

// pragma solidity ^0.8.0;

// import "@polytope-labs/solidity-merkle-trees/MerklePatricia.sol";

// contract CheckFraud {

//     event ValuesVerified(bytes[] values);

//     struct ResponseMsg {
//         string Type;
//         bytes32 ChannelId;
//         uint256 Amount; 
//         bytes SignedReqBody;
//         uint256 CurrentBlockHeight;
//         bytes ReturnValue;
//         string[] Proof;
//         bytes32 TxHash;
//         uint32 TxIdx;
//         bytes Signature;
//     }

//     event DecodedResponse(
//         string Type
//     );

//     function decodeResponse(bytes memory serializedData) public {
// uint256 offset = 0;

//         // Step 1: Decode each component using helper functions
//         (string memory Type, uint256 newOffset) = _decodeString(serializedData, offset);
//         offset = newOffset;

//         (bytes32 ChannelId) = _decodeBytes32(serializedData, offset);
//         offset += 32;

//         (uint256 Amount) = _decodeUint64(serializedData, offset);
//         offset += 8;

//         (bytes memory SignedReqBody, uint256 newOffset2) = _decodeBytes(serializedData, offset);
//         offset = newOffset2;

//         (uint256 CurrentBlockHeight, uint256 newOffset3) = _decodeUint256(serializedData, offset);
//         offset = newOffset3;

//         (bytes memory ReturnValue, uint256 newOffset4) = _decodeBytes(serializedData, offset);
//         offset = newOffset4;

//         (string[] memory Proof, uint256 newOffset5) = _decodeStringArray(serializedData, offset);
//         offset = newOffset5;

//         (bytes32 TxHash) = _decodeBytes32(serializedData, offset);
//         offset += 32;

//         (uint32 TxIdx) = _decodeUint32(serializedData, offset);
//         offset += 4;

//         (bytes memory Signature, uint256 newOffset6) = _decodeBytes(serializedData, offset);
//         offset = newOffset6;

//         // Emit event to verify decoding
//         emit DecodedResponse(
//             Type
//         );
//     }

//     function _decodeString(bytes memory data, uint256 offset) internal pure returns (string memory, uint256) {
//         uint32 length;
//         assembly {
//             length := mload(add(data, add(offset, 0x20)))
//         }
//         offset += 4;

//         bytes memory strBytes = new bytes(length);
//         for (uint256 i = 0; i < length; i++) {
//             strBytes[i] = data[offset + i];
//         }
//         offset += length;

//         return (string(strBytes), offset);
//     }

//     function _decodeBytes32(bytes memory data, uint256 offset) internal pure returns (bytes32) {
//         bytes32 result;
//         assembly {
//             result := mload(add(data, add(offset, 0x20)))
//         }
//         return result;
//     }

//     function _decodeUint64(bytes memory data, uint256 offset) internal pure returns (uint64) {
//         uint64 result;
//         assembly {
//             result := mload(add(data, add(offset, 0x8)))
//         }
//         return result;
//     }

//     function _decodeUint32(bytes memory data, uint256 offset) internal pure returns (uint32) {
//         uint32 result;
//         assembly {
//             result := mload(add(data, add(offset, 0x20)))
//         }
//         return result;
//     }

//     function _decodeUint256(bytes memory data, uint256 offset) internal pure returns (uint256, uint256) {
//         uint256 result;
//         assembly {
//             result := mload(add(data, add(offset, 0x20)))
//         }
//         offset += 32;
//         return (result, offset);
//     }

//     function _decodeBytes(bytes memory data, uint256 offset) internal pure returns (bytes memory, uint256) {
//         uint32 length;
//         assembly {
//             length := mload(add(data, add(offset, 0x20)))
//         }
//         offset += 4;

//         bytes memory result = new bytes(length);
//         for (uint256 i = 0; i < length; i++) {
//             result[i] = data[offset + i];
//         }
//         offset += length;

//         return (result, offset);
//     }

//     function _decodeStringArray(bytes memory data, uint256 offset) internal pure returns (string[] memory, uint256) {
//         uint32 count;
//         assembly {
//             count := mload(add(data, add(offset, 0x20)))
//         }
//         offset += 4;

//         string[] memory result = new string[](count);
//         for (uint256 i = 0; i < count; i++) {
//             (string memory element, uint256 newOffset) = _decodeString(data, offset);
//             result[i] = element;
//             offset = newOffset;
//         }

//         return (result, offset);
//     }

//     // function 
// 	// Type               string
// 	// ChannelId          common.Hash
// 	// Amount             uint
// 	// SignedReqBody      []byte
// 	// CurrentBlockHeight *big.Int
// 	// ReturnValue        []byte
// 	// Proof              []string
// 	// TxHash             common.Hash
// 	// TxIdx              uint32
// 	// Signature          []byte




// }