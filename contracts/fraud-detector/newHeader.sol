// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "./libs/rlp/RLPReader.sol";
import "./libs/rlp/RLPEncoder.sol";

library HeaderDecoder {
    using RLPReader for bytes;
    using RLPReader for RLPReader.RLPItem;
    using RLPEncoder for bytes;

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
    event BasicHeaderDecoded(
        bytes32 parentHash,
        bytes32 uncleHash,
        address coinbase,
        bytes32 stateRoot,
        bytes32 txRoot,
        bytes32 receiptRoot
    );

    event ExtendedHeaderDecoded(
        bytes logsBloom,
        uint256 difficulty,
        uint256 number,
        uint256 gasLimit,
        uint256 gasUsed,
        uint256 timestamp,
        bytes extraData,
        bytes32 mixDigest,
        uint64 nonce
    );

    event OptionalFieldsDecoded(
        uint256 baseFee,
        bytes32 withdrawalsRoot,
        uint256 blobGasUsed,
        uint256 excessBlobGas,
        bytes32 parentBeaconRoot
    );

    event HeaderHashGenerated(bytes32 headerHash);
    event LogUint(uint);

    function decodeHeader(bytes memory rlpEncodedHeader) public returns (bytes32, bytes32, bytes32) {
        RLPReader.RLPItem[] memory items = rlpEncodedHeader.toRlpItem().toList();

        // Decode and emit basic fields
        BasicHeader memory basic = _decodeAndEmitBasicHeader(items);

        // Decode and emit extended fields
        ExtendedHeader memory extended = _decodeAndEmitExtendedHeader(items);

        // Decode and emit optional fields
        OptionalFields memory optional = _decodeAndEmitOptionalFields(items);

        bytes32 headerHash = generateBlockHash(basic, extended, optional);
        emit HeaderHashGenerated(headerHash);
        emit HeaderHashGenerated(basic.txRoot);
        emit HeaderHashGenerated(basic.stateRoot);
        return (headerHash, basic.txRoot, basic.stateRoot);

        // bytes32 blockHash = blockhash(extended.number);
        // emit LogUint(extended.number); // Debug the block number
        // emit LogBytes32(blockHash);
    }

    function _decodeAndEmitBasicHeader(RLPReader.RLPItem[] memory items) internal returns (BasicHeader memory) {
        bytes32 parentHash = bytes32(items[0].toBytes());
        bytes32 uncleHash = bytes32(items[1].toBytes());
        address coinbase = address(bytes20(items[2].toAddress()));
        bytes32 stateRoot = bytes32(items[3].toBytes());
        bytes32 txRoot = bytes32(items[4].toBytes());
        bytes32 receiptRoot = bytes32(items[5].toBytes());
        emit BasicHeaderDecoded(parentHash, uncleHash, coinbase, stateRoot, txRoot, receiptRoot);
        return BasicHeader({
            parentHash: bytes32(items[0].toBytes()),
            uncleHash: bytes32(items[1].toBytes()),
            coinbase: address(bytes20(items[2].toAddress())),
            stateRoot: bytes32(items[3].toBytes()),
            txRoot: bytes32(items[4].toBytes()),
            receiptRoot: bytes32(items[5].toBytes())
        });

    }

    function _decodeAndEmitExtendedHeader(RLPReader.RLPItem[] memory items) internal returns (ExtendedHeader memory) {
        bytes memory logsBloom = items[6].toBytes();
        uint256 difficulty = items[7].toUint();
        uint256 number = items[8].toUint();
        uint256 gasLimit = items[9].toUint();
        uint256 gasUsed = items[10].toUint();
        uint256 timestamp = items[11].toUint();
        bytes memory extraData = items[12].toBytes();
        bytes32 mixDigest = bytes32(items[13].toBytes());
        uint64 nonce = uint64(items[14].toUint());


        emit ExtendedHeaderDecoded(
            logsBloom,
            difficulty,
            number,
            gasLimit,
            gasUsed,
            timestamp,
            extraData,
            mixDigest,
            nonce
        );
        return ExtendedHeader({
            logsBloom: items[6].toBytes(),
            difficulty: items[7].toUint(),
            number: items[8].toUint(),
            gasLimit: items[9].toUint(),
            gasUsed: items[10].toUint(),
            timestamp: items[11].toUint(),
            extraData: items[12].toBytes(),
            mixDigest: bytes32(items[13].toBytes()),
            nonce: uint64(items[14].toUint())
        });
    }

    function _decodeAndEmitOptionalFields(RLPReader.RLPItem[] memory items) internal returns (OptionalFields memory) {
        uint256 baseFee = 0;
        bytes32 withdrawalsRoot = bytes32(0);
        uint256 blobGasUsed = 0;
        uint256 excessBlobGas = 0;
        bytes32 parentBeaconRoot = bytes32(0);

        if (items.length > 15) {
            baseFee = items[15].toUint();
        }
        if (items.length > 16) {
            withdrawalsRoot = bytes32(items[16].toBytes());
        }
        if (items.length > 17) {
            blobGasUsed = items[17].toUint();
        }
        if (items.length > 18) {
            excessBlobGas = items[18].toUint();
        }
        if (items.length > 19) {
            parentBeaconRoot = bytes32(items[19].toBytes());
        }

        emit OptionalFieldsDecoded(baseFee, withdrawalsRoot, blobGasUsed, excessBlobGas, parentBeaconRoot);
        return OptionalFields({
            baseFee: baseFee,
            withdrawalsRoot: withdrawalsRoot,
            blobGasUsed: blobGasUsed,
            excessBlobGas: excessBlobGas,
            parentBeaconRoot: parentBeaconRoot
        });
    }


    event LogRLPEncoded(bytes);
    event LogBytes32(bytes32);
    function generateBlockHash(
        BasicHeader memory basic,
        ExtendedHeader memory extended,
        OptionalFields memory optional
    ) internal pure returns (bytes32) {

        bytes[] memory elements = new bytes[](20);
        elements[0] = RLPEncoder.encodeBytes32(basic.parentHash);
        elements[1] = RLPEncoder.encodeBytes32(basic.uncleHash);
        elements[2] = RLPEncoder.encodeAddress(basic.coinbase); 
        elements[3] = RLPEncoder.encodeBytes32(basic.stateRoot);   
        elements[4] = RLPEncoder.encodeBytes32(basic.txRoot);  
        elements[5] = RLPEncoder.encodeBytes32(basic.receiptRoot);  

        elements[6] = RLPEncoder.encodeBytes(extended.logsBloom);
        elements[7] = RLPEncoder.encodeUint(extended.difficulty);
        elements[8] = RLPEncoder.encodeUint(extended.number);
        elements[9] = RLPEncoder.encodeUint(extended.gasLimit);
        elements[10] = RLPEncoder.encodeUint(extended.gasUsed);
        elements[11] = RLPEncoder.encodeUint(extended.timestamp);
        elements[12] = RLPEncoder.encodeBytes(extended.extraData);
        elements[13] = RLPEncoder.encodeBytes32(extended.mixDigest);
        elements[14] = RLPEncoder.encodeBytes(abi.encodePacked(uint64(extended.nonce)));

        elements[15] = optional.baseFee > 0 ? RLPEncoder.encodeUint(optional.baseFee) : RLPEncoder.encodeEmpty();
        elements[16] = optional.withdrawalsRoot != bytes32(0) ? RLPEncoder.encodeBytes32(optional.withdrawalsRoot) : RLPEncoder.encodeEmpty();
        elements[17] = optional.blobGasUsed > 0 ? RLPEncoder.encodeUint(optional.blobGasUsed) : RLPEncoder.encodeEmpty();
        elements[18] = optional.excessBlobGas > 0 ? RLPEncoder.encodeUint(optional.excessBlobGas) : RLPEncoder.encodeEmpty();
        elements[19] = optional.parentBeaconRoot != bytes32(0) ? RLPEncoder.encodeBytes32(optional.parentBeaconRoot) : RLPEncoder.encodeEmpty();

        bytes memory rlpRes = RLPEncoder.encodeList(elements);

        // emit LogRLPEncoded(rlpRes); // Emit the encoded result for comparison
        return keccak256(rlpRes);
    }
}