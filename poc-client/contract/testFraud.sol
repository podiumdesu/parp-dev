// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "./rlp/Helper.sol";  // Assuming you have added an RLP decoding library

import { MerkleVerify as MerkleVerify} from "./merkleverify.sol";
import "./Types.sol";

// import "@polytope-labs/solidity-merkle-trees/MerklePatricia.sol";


    
contract FraudProofDecoder {
    using RLPReader for bytes;
    using RLPReader for RLPReader.RLPItem;


    address constant public fullNode = 0xC8A7ae3f6Ae079c20BA19164089143F48F7B965f;
    address constant public lightClient = 0xD25a31702b7b86B2e953Baf9ff88Ef716A5306Cc;

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

    event ValuesVerified(bytes[] values);

    event msgHash(bytes32 msgHash);
    event txHash(bytes32 txHassssh);
    event txIdx(bytes txIdx);
    event reqByte(bytes SignedReqBody);
    event emitProof(bytes[] Proof);

    function decodeAll(bytes memory res, bytes memory req) public returns (address, address) {

        ResponseMsg memory response = decodeResponse(res);
        RequestBody memory request = decodeRequest(req);


        bytes32 resHash = getMessageHash(response.Proof, response.TxHash, response.TxIdx, response.SignedReqBody);
        emit msgHash(resHash);
        address responseSigner = recoverSignature(resHash, response.Signature);
        require(responseSigner == fullNode, "It must be a vaid response from the full node.");

        bytes32 reqHash = getReqMessageHash(request.Amount, request.LocalBlockHash, request.ReqBytes);
        emit msgHash(reqHash);
        address requestSigner = recoverSignature(reqHash, response.SignedReqBody);
        require(requestSigner == lightClient, "It must be a valid request from the light client.");

        require(request.Amount == response.Amount, "The amount needs to be the same");


        // -------STEP: verify Merkle Proof
        bytes32 root = 0xbb345e208bda953c908027a45aa443d6cab6b8d2fd64e83ec52f1008ddeafa58;
        bytes[] memory key = new bytes[](1);
        key[0] = response.TxIdx;

        verifyProof(root, response.Proof, key);

        return (requestSigner, responseSigner);
    }

    event DecodeRequest(
        // uint32 channelId,
        uint Amount,
        bytes32 LocalBlockHash,
        bytes ReqBytes,
        // uint length
        uint length
    );


    function verifyRequester (bytes memory req, bytes memory signature) public returns (address){
        RequestBody memory request = decodeRequest(req);
        bytes32 reqHash = getReqMessageHash(request.Amount, request.LocalBlockHash, request.ReqBytes);
        emit msgHash(reqHash);
        address requestSigner = recoverSignature(reqHash, signature);
        require(requestSigner == lightClient, "It must be a valid request from the light client.");
        return requestSigner;
    }
    
    function decodeResponse(bytes memory res) public returns (ResponseMsg memory) {
        // Step 1: Parse the RLP encoded data
        RLPReader.RLPItem[] memory items = res.toRlpItem().toList();
        
        // Step 2: Check if the correct number of fields is present in the RLP-encoded message
        require(items.length == 11, "Incorrect number of fields in RLP encoded data");

        ResponseMsg memory response;
        // Decoding string Type
        response.Type = string(items[0].toBytes());
        
        // Decoding ChannelId (which is a bytes32)
        response.ChannelId = bytes32(items[1].toBytes());

        // Decoding Amount (which is uint)
        response.Amount = items[2].toUint();

        response.ReqBodyHash = bytes32(items[3].toBytes());
        // Decoding SignedReqBody (which is bytes)
        response.SignedReqBody = items[4].toBytes();

        // Decoding CurrentBlockHeight (which is uint)
        response.CurrentBlockHeight = items[5].toUint();

        // Decoding ReturnValue (which is bytes)
        response.ReturnValue = items[6].toBytes();

        // Decoding Proof as an array of strings
        RLPReader.RLPItem[] memory proofItems = items[7].toList();
        response.Proof = new bytes[](proofItems.length);
        for (uint256 i = 0; i < proofItems.length; i++) {
            response.Proof[i] = bytes(proofItems[i].toBytes());
        }

        // Decoding TxHash (which is bytes32)
        response.TxHash = bytes32(items[8].toBytes());

        // Decoding TxIdx (which is uint32)
        response.TxIdx = items[9].toBytes();

        // Decoding Signature (which is bytes)
        response.Signature = items[10].toBytes();

        // Emit the decoded data as an event to verify the correctness
        emit DecodedResponse(
            response.Type,
            response.ChannelId,
            response.Amount,
            response.ReqBodyHash,
            response.SignedReqBody,
            response.CurrentBlockHeight,
            response.ReturnValue,
            response.Proof,
            response.TxHash,
            response.TxIdx,
            response.Signature
        );

        bytes32 resHash = getMessageHash(response.Proof, response.TxHash, response.TxIdx, response.SignedReqBody);
        emit msgHash(resHash);

        return response;
    }

    function getMessageHash(
        bytes[] memory Proof,
        bytes32 TxHash,
        bytes memory TxIdx,
        bytes memory SignedReqBody
    ) public returns (bytes32) {
        // emit txHash(TxHash);
        // emit txIdx(TxIdx);
        // emit emitProof(Proof);
        emit reqByte(SignedReqBody);
        bytes memory proofBytes;
        for (uint256 i = 0; i < Proof.length; i++) {
            proofBytes = abi.encodePacked(proofBytes, Proof[i]);
        }
        return keccak256(abi.encodePacked(
            SignedReqBody,
            proofBytes,
            TxHash,
            TxIdx
        ));
    }

    struct RequestBody {
        // uint32 ChannelId;
        uint Amount;
        bytes32 LocalBlockHash;
        bytes ReqBytes;
    }



    event Debug(string level, uint number);
    function decodeRequest(bytes memory serializedReqBody) public returns(RequestBody memory) {
        RLPReader.RLPItem[] memory items = serializedReqBody.toRlpItem().toList();
        
        // Step 2: Check if the correct number of fields is present in the RLP-encoded message
        require(items.length == 3, "Incorrect number of fields in RLP encoded data");
        // emit Debug("Number of fields found", items.length); 
        
        RequestBody memory reqBody;

        // ----- NOTE: REMEMBER THE DEFINATION LOCALLY WILL CHANGE THE ITEM LENGTH


        // // Decoding ChannelId (which is a uint32)
        // reqBody.ChannelId = uint32(items[0].toUint());

        // // emit emitChannelId(reqBody.ChannelId);
        reqBody.Amount = items[0].toUint();
        emit emitAmount(reqBody.Amount);
        
        // // Decoding Proof as an array of strings
        
        reqBody.LocalBlockHash = bytes32(items[1].toBytes());
        emit emitLocalBlockHash(reqBody.LocalBlockHash);

        reqBody.ReqBytes = items[2].toBytes(); 

        emit DecodeRequest(reqBody.Amount, reqBody.LocalBlockHash, reqBody.ReqBytes, items.length);

        bytes32 messageHash = getReqMessageHash(reqBody.Amount, reqBody.LocalBlockHash, reqBody.ReqBytes);
        emit msgHash(messageHash);
        return reqBody;
        
    }

    event emitChannelId(uint32 channelId);
    event emitAmount(uint amount); 
    event emitReqByte(bytes[] reqByte);
    event emitLocalBlockHash(bytes32 LocalBlockHash);

    function getReqMessageHash(
        uint Amount,
        bytes32 LocalBlockHash,
        bytes memory ReqBytes
    ) public pure returns (bytes32) {
        return keccak256(abi.encodePacked(
            Amount,
            LocalBlockHash,
            ReqBytes
        ));
    } 

    function recoverSignature(
        bytes32 messageHash,
        bytes memory signature
    ) internal pure returns (address) {
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


    event emitProofValues(StorageValue[] proofReturnValues);

    function verifyProof(
        bytes32 root,
        bytes[] memory proof,
        bytes[] memory keys
    ) public pure returns (StorageValue[] memory) {
        StorageValue[] memory values = MerkleVerify.VerifyEthereumProof(root, proof, keys);
        // emit emitProofValues(values);
        return values;
    }

}
