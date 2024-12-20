// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "./interfaces/IPayChan.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "./libs/rlp/RLPReader.sol";

import "./types/message.sol";
import "./types/channel.sol";

import "./libs/decode/msgDecoding.sol";
import "./libs/ProofHelper.sol";

import "./FraudProofBase.sol";

abstract contract FraudProofTxProcessor is FraudProofBase {

    using RLPReader for bytes;
    using RLPReader for RLPReader.RLPItem;

    using FraudProofDecoderLibrary for bytes;
    using FraudProofHelper for bytes;

    // constructor(address _paychanContractAddress, address _depositContract)
    //     FraudProofBase(_paychanContractAddress, _depositContract)
    // {}

    function fraudProofTx(bytes memory res, bytes memory req) internal returns (bool, address, address, address) {
        (RequestBody memory request, bytes32 reqHash) = FraudProofDecoderLibrary.decodeRequest(req);
        (ResponseMsg memory response, bytes32 resHash) = FraudProofDecoderLibrary.decodeResponse(res);
        emit msgHash(resHash);
        emit msgHash(reqHash);

        ChannelInfo memory channelInfo;
        (channelInfo.sender, channelInfo.recipient, channelInfo.status, ) = paychanContract.paychanSelectedArguments(request.ChannelId);
        
        // just for test 
        // recipient = 0xC8A7ae3f6Ae079c20BA19164089143F48F7B965f;
        // sender = 0xD25a31702b7b86B2e953Baf9ff88Ef716A5306Cc;
        require(channelInfo.status != 0, "The channel must not be closed.");
        
        // string memory msgType = getType(res);
    
        bytes32 txRootHash = response.TxRootHash;
        bytes[] memory key = new bytes[](1);
        key[0] = response.TxIdx;


        address responseSigner = ECDSA.recover(resHash, response.Signature);
        require(responseSigner == channelInfo.recipient, "It must be a vaid response from the full node.");

        address requestSigner = ECDSA.recover(reqHash, response.SignedReqBody);
        require(requestSigner == channelInfo.sender, "It must be a valid request from the light client.");

        require(request.Amount == response.Amount, "The amount needs to be the same");
        
        bool proofStatus = FraudProofHelper.verifyProof(txRootHash, response.Proof, key);
        require(!proofStatus, "Fraud proof is valid. Full node is honest.");
        return (true, channelInfo.sender, channelInfo.recipient, channelInfo.witness);
        // slashWithAddresses(channelInfo.sender, channelInfo.recipient, channelInfo.witness);

        // return (requestSigner, responseSigner);

    }
}