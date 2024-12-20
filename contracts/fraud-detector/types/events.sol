// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "./proof.sol";

    event ValuesVerified(bytes[] values);

    event msgHash(bytes32 msgHash);
    event txHash(bytes32 txHassssh);
    event txIdx(bytes txIdx);
    event reqByte(bytes SignedReqBody);
    event emitProof(bytes[] Proof);
    event emitAddress(address signer);
    event LogBool(bool);
    event LogBytes(bytes);
    event emitProofValues(StorageValue[] proofReturnValues);


    event DecodeRequest(
        bytes32 channelId,
        uint Amount,
        bytes32 LocalBlockHash,
        bytes ReqBytes,
        uint length
    );
    event emitChannelId(uint32 channelId);
    event emitAmount(uint amount); 
    event emitReqByte(bytes[] reqByte);
    event emitLocalBlockHash(bytes32 LocalBlockHash);