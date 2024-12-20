// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "./interfaces/IPayChan.sol";
import "./interfaces/IDeposit.sol";
import "./FraudProofTxProcessor.sol";
import "./FraudProofSPProcessor.sol";

import "./libs/rlp/Helper.sol";  // Assuming you have added an RLP decoding library

import "./libs/decode/msgDecoding.sol";

import "./newHeader.sol";

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

contract FraudProofDecoder is FraudProofTxProcessor, FraudProofSPProcessor {
    // constructor(address _paychanContractAddress, address _depositContract)
    //     FraudProofTxProcessor(_paychanContractAddress, _depositContract)
    //     FraudProofSPProcessor(_paychanContractAddress, _depositContract)
    // {}

    using RLPReader for bytes;
    using RLPReader for RLPReader.RLPItem;

    using FraudProofDecoderLibrary for bytes;
    using HeaderDecoder for bytes;

    constructor(address _paychanContractAddress, address _depositContract)
        FraudProofBase(_paychanContractAddress, _depositContract)
    {}
    // address constant public fullNode = 0xC8A7ae3f6Ae079c20BA19164089143F48F7B965f;
    // address constant public lightClient = 0xD25a31702b7b86B2e953Baf9ff88Ef716A5306Cc;

    function decodeAllSP(bytes memory res, bytes memory req, bytes memory blockHeaderInfo, address witness) public {        
        string memory msgType = getType(res);
        bool fraudDetected;

        bytes memory signedReqBody;
        bytes memory resSignature;
        uint256 blockNr;

        // Step 1: Decode Request
        (RequestBody memory request, bytes32 reqHash) = FraudProofDecoderLibrary.decodeRequest(req);

        // Step 2: Decode Response
        ResponseMsg memory response;
        ResponseSPMsg memory responseSP;

        bytes32 resHash;
        if (keccak256(abi.encodePacked(msgType)) == keccak256(abi.encodePacked("response"))) {
            (response, resHash) = FraudProofDecoderLibrary.decodeResponse(res);
            // (fraudDetected, sender, recipient, witness) = fraudProofTx(res, req);
            require(request.ChannelId == response.ChannelId, "Channel ID must be the same.");
            signedReqBody = response.SignedReqBody;
            resSignature = response.Signature;
            blockNr = response.CurrentBlockHeight;
        } else {
            (responseSP, resHash) = FraudProofDecoderLibrary.decodeResponseSP(res);
            // (fraudDetected, sender, recipient, witness) = fraudProofSP(res, req);
            require(request.ChannelId == responseSP.ChannelId, "Channel ID must be the same.");
            signedReqBody = responseSP.SignedReqBody;
            resSignature = responseSP.Signature;
            blockNr = responseSP.CurrentBlockHeight;
        }

        // Step 3: Get Channel Information
        ChannelInfo memory channelInfo;
        (channelInfo.sender, channelInfo.recipient, channelInfo.status, ) = paychanContract.paychanSelectedArguments(request.ChannelId);
        require(channelInfo.status != 0, "The channel must not be closed.");

        // Step 4: (based on Step 3) verify request and sender from channel information
        address requestSigner = ECDSA.recover(reqHash, signedReqBody);
        require(requestSigner == channelInfo.sender, "It must be a valid request from the light client.");

        // Step 4: (based on Step 3) verify response and sender from channel information
        address responseSigner = ECDSA.recover(resHash, resSignature);
        require(responseSigner == channelInfo.recipient, "It must be a vaid response from the full node.");


        // Step 5: Fraud proof detect
        bytes[] memory key = new bytes[](1);
        (bytes32 headerHash, bytes32 txRoot, bytes32 stateRoot) = HeaderDecoder.decodeHeader(blockHeaderInfo);

        bytes32 blockHash = blockhash(blockNr);
        require(blockHash == headerHash, "Cant trust your root values");


        bool proofStatus;
        if (keccak256(abi.encodePacked(msgType)) == keccak256(abi.encodePacked("response"))) {
            key[0] = response.TxIdx;
            proofStatus = FraudProofHelper.verifyProof(txRoot, response.Proof, key);
        } else {
            key[0] = responseSP.Address;
            emit msgHash(stateRoot);
            emit emitProof(responseSP.Proof);
            emit emitProof(key);
            proofStatus = FraudProofHelper.verifyProof(stateRoot, responseSP.Proof, key);
        }

        // TODO: Change it back to correct 
        emit LogBool(proofStatus);
        require(proofStatus == true, "Fraud proof is valid. Full node is honest.");

        slashWithAddresses(channelInfo.sender, channelInfo.recipient, witness);
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

