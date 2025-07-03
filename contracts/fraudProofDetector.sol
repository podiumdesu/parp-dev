// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "./fraud-detector/interfaces/IPayChan.sol";
import "./fraud-detector/interfaces/IDeposit.sol";
import "./fraud-detector/FraudProofTxProcessor.sol";
import "./fraud-detector/FraudProofSPProcessor.sol";

import "./fraud-detector/libs/rlp/Helper.sol";  // Assuming you have added an RLP decoding library

import "./fraud-detector/libs/decode/msgDecoding.sol";

import "./fraud-detector/newHeader.sol";

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";


contract FraudProofDecoder is FraudProofTxProcessor, FraudProofSPProcessor {

    using RLPReader for bytes;
    using RLPReader for RLPReader.RLPItem;

    using FraudProofDecoderLibrary for bytes;
    using HeaderDecoder for bytes;


    struct Context {
        bytes signedReqBody;
        bytes resSignature;
        uint256 blockNr;
        bytes32 reqHash;
        bytes32 resHash;
        bool isSP;
        bytes32 channelId;
        bytes proofKey;
        bytes[] proof;
    }

    constructor(address _paychanContractAddress, address _depositContract)
        FraudProofBase(_paychanContractAddress, _depositContract)
    {}
    // address constant public fullNode = 0xC8A7ae3f6Ae079c20BA19164089143F48F7B965f;
    // address constant public lightClient = 0xD25a31702b7b86B2e953Baf9ff88Ef716A5306Cc;


    function _decodeReqResAndChanValid (
        bytes memory res,
        bytes memory req
    ) internal returns (Context memory ctx) {
        // Step 1: Decode Request
        (RequestBody memory request, bytes32 reqHash) = FraudProofDecoderLibrary.decodeRequest(req);
        ctx.reqHash = reqHash;

        // Setep 2: Decode Response
        ctx.isSP = (keccak256(abi.encodePacked(getType(res)))) != keccak256(abi.encodePacked("response"));
        if (!ctx.isSP) {
            (ResponseMsg memory response, bytes32 resHash) = FraudProofDecoderLibrary.decodeResponse(res);
            require(request.ChannelId == response.ChannelId, "Channel ID must be the same.");
            ctx.channelId = request.ChannelId;
            ctx.signedReqBody = response.SignedReqBody;
            ctx.resSignature = response.Signature;
            ctx.blockNr = response.CurrentBlockHeight;
            ctx.proofKey = response.TxIdx;
            ctx.proof = response.Proof;
            ctx.resHash = resHash;
        } else {
            (ResponseSPMsg memory responseSP, bytes32 resHash) = FraudProofDecoderLibrary.decodeResponseSP(res);
            require(request.ChannelId == responseSP.ChannelId, "Channel ID must be the same.");
            ctx.channelId = request.ChannelId;
            ctx.signedReqBody = responseSP.SignedReqBody;
            ctx.resSignature = responseSP.Signature;
            ctx.blockNr = responseSP.CurrentBlockHeight;
            ctx.proofKey = responseSP.Address;
            ctx.proof = responseSP.Proof;
            ctx.resHash = resHash;
        }
    }

    function fraudProofDetector(bytes memory res, bytes memory req, bytes memory blockHeaderInfo, address witness) public {        

        Context memory ctx = _decodeReqResAndChanValid(res, req);

        // Step 3: Get Channel Information
        ChannelInfo memory channelInfo;
        (channelInfo.sender, channelInfo.recipient, channelInfo.status, ) = paychanContract.paychanSelectedArguments(ctx.channelId);
        require(channelInfo.status != 0, "The channel must not be closed.");

        // Step 4: (based on Step 3) verify request and sender from channel information
        address requestSigner = ECDSA.recover(ctx.reqHash, ctx.signedReqBody);
        require(requestSigner == channelInfo.sender, "It must be a valid request from the light client.");

        // Step 4: (based on Step 3) verify response and sender from channel information
        address responseSigner = ECDSA.recover(ctx.resHash, ctx.resSignature);
        require(responseSigner == channelInfo.recipient, "It must be a vaid response from the full node.");


        // Step 5: Fraud proof detect

        HeaderDecoder.HeaderResults memory header = HeaderDecoder.decodeHeader(blockHeaderInfo);

        bytes32 blockHash = blockhash(ctx.blockNr);
        require(blockHash == header.headerHash, "Cant trust your root values");
        
        bool proofStatus;
        proofStatus = verifyFraudDetection(ctx.isSP, header, ctx.proofKey, ctx.proof);

        // If proofStatus is true, it means the merkle proof is not a fraud
        require(proofStatus == false, "Fraud proof is valid. Full node is honest.");
        slashWithAddresses(channelInfo.sender, channelInfo.recipient, witness);
    }
    
    function verifyFraudDetection(
        bool isSP,
        HeaderDecoder.HeaderResults memory header, 
        bytes memory proofKey,
        bytes[] memory proof
    ) internal returns (bool) {
        bool proofStatus;
        bytes[] memory key = new bytes[](1);

        key[0] = proofKey;
        if (!isSP) {
            proofStatus = FraudProofHelper.verifyProof(header.txRoot, proof, key);
        } else {
            proofStatus = FraudProofHelper.verifyProof(header.stateRoot, proof, key);
        }
        emit LogBool(proofStatus);
        return proofStatus;
    }

    function slashWithAddresses(address lc, address fn, address witness) internal {
        depositContract.slash(fn, lc, witness);
    }

    function getType(bytes memory res) internal pure returns (string memory) {
        RLPReader.RLPItem[] memory items = res.toRlpItem().toList();

        // Ensure the correct number of fields
        require(items.length > 1, "Incorrect number of fields in RLP encoded data");

        ResponseMsg memory responseSP;

        // Decode fields
        responseSP.Type = string(items[0].toBytes());
        
        return responseSP.Type;
    }


    receive() external payable {}

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
            emit LogBool(true);
            return true;
        } else {
            return false;
        }
    }

}

