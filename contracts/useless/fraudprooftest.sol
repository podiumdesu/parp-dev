// SPDX‑License‑Identifier: MIT

pragma solidity ^0.8.17;

// 1) import just the library you care about
import "./libs/decode/msgDecoding.sol";

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

/// @notice A bare‑bones wrapper that calls your library’s decodeRequest.
///         Deploy this + the Library; no other contracts needed.
contract DecoderTester {
    event SignerRecovered(address signer);

    /// @notice Call the library’s decodeRequest and get back both the struct and the hash
    function testDecodeRequest(
        bytes calldata serializedReqBody,
        bytes calldata res
    )
        external
        returns (
            bytes32 channelId,
            uint256 amount,
            bytes32 localBlockHash,
            bytes memory reqBytes,
            bytes32 messageHash,
            address signer
        )
    {
        // this will still emit your events (DecodeRequest + msgHash)
        (RequestBody memory r, bytes32 h) = FraudProofDecoderLibrary
            .decodeRequest(serializedReqBody);

        (ResponseMsg memory response, bytes32 resHash) = FraudProofDecoderLibrary
            .decodeResponse(res);

        // unpack into Remix‑friendly return values
        channelId = r.ChannelId;
        amount = r.Amount;
        localBlockHash = r.LocalBlockHash;
        reqBytes = r.ReqBytes;
        messageHash = h;

        // Step 4: (based on Step 3) verify request and sender from channel information
        address requestSigner = ECDSA.recover(h, response.SignedReqBody);
        require(
            requestSigner == 0xD25a31702b7b86B2e953Baf9ff88Ef716A5306Cc,
            "It must be a valid request from the light client."
        );
        emit SignerRecovered(requestSigner);

        // Step 4: (based on Step 3) verify response and sender from channel information
        address responseSigner = ECDSA.recover(resHash, response.Signature);
        require(
            responseSigner == 0xC8A7ae3f6Ae079c20BA19164089143F48F7B965f,
            "It must be a vaid response from the full node."
        );
        emit SignerRecovered(responseSigner);

        return (
            channelId,
            amount,
            localBlockHash,
            reqBytes,
            messageHash,
            requestSigner
        );

        // it has been tested, and it has no problem.
    }
}
