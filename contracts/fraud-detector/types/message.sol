// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

struct RequestBody {
    bytes32 ChannelId;
    uint Amount;
    bytes32 LocalBlockHash;
    bytes ReqBytes;
}

struct ResponseMsg {
    string Type;
    bytes32 ChannelId;
    uint256 Amount;
    bytes32 ReqBodyHash;
    bytes SignedReqBody;
    uint256 CurrentBlockHeight;
    bytes ReturnValue;
    bytes[] Proof;
    bytes32 TxHash;
    bytes TxIdx;
    bytes Signature;
    bytes32 TxRootHash;
}

struct ResponseSPMsg {
    string Type;
    bytes32 ChannelId;
    uint256 Amount;
    bytes32 ReqBodyHash;
    bytes SignedReqBody;
    uint256 CurrentBlockHeight;
    bytes ReturnValue;
    bytes[] Proof;
    bytes Address;
    bytes Signature;
    bytes32 TxRootHash;
}

event DecodedResponse(
    string Type,
    bytes32 ChannelId,
    uint256 Amount,
    bytes32 ReqBodyHash,
    bytes SignedReqBody,
    uint256 CurrentBlockHeight,
    bytes ReturnValue,
    bytes[] Proof,
    bytes32 TxHash,
    bytes TxIdx,
    bytes Signature
);

event DecodedResponseSP(
    string Type,
    bytes32 ChannelId,
    uint256 Amount,
    bytes32 ReqBodyHash,
    bytes SignedReqBody,
    uint256 CurrentBlockHeight,
    bytes ReturnValue,
    bytes[] Proof,
    bytes Address,
    bytes Signature
);
