// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

import "./rlp/Helper.sol";  // Assuming you have added an RLP decoding library
import "./common.sol";

import "./fraud-detector/interfaces/IDeposit.sol";

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";


contract paychan {
    using RLPReader for bytes;
    using RLPReader for RLPReader.RLPItem;
    SigHelper sigHelper;

    IDepositContract public depositContract;


    event ChannelOpened(bytes32 indexed channelId);
    event DebugInfo(address indexed recoveredSigner, bytes32 indexed channelId);

    constructor(address _depositContractAddress) {
        depositContract = IDepositContract(_depositContractAddress);
    }

    function isEligible(address user) internal view returns (bool) {
        uint256 userDeposit = depositContract.getDeposit(user);
        return userDeposit >= 1 ether;
    }

    struct PayChan {
        bytes32 id; // store the payment channel id

        address payable sender; // store the sender
        address payable recipient; // store the recipient

        // funds
        uint senderDeposit;

        // expiry date of the channel
        uint startTime;
        // uint duration;
        // uint endTime;
        uint status;     // 1 = open, 0 = closed, 2 = waitingClosure

        // state that will be updated
        uint fee;

        // For closing the channel
        uint disputeStartTime;
        uint disputeDuration;    // default 100s
        bool senderConfirm;
        bool recipientConfirm;

    }

    mapping (bytes32 => PayChan) private paychans;   // store every single payment channel

    modifier onlyParticipants(bytes32 channelId) {
        require(msg.sender == paychans[channelId].sender || msg.sender == paychans[channelId].recipient, "Only participants");
        _;
    }

    modifier onlySender(bytes32 channelId) {
        require(paychans[channelId].sender == msg.sender, "Only sender");
        _;
    }

    modifier onlyRecipient(bytes32 channelId) {
        require(paychans[channelId].recipient == msg.sender, "Only recipient");
        _;
    }

    modifier onlyWhenChannelOpen(bytes32 channelId) {
        require(paychans[channelId].status == 1, "Only when channel is open");
        _;
    }

    function openChan(address to, uint senderDeposit) public payable returns(bytes32 channelId) {
        // require(duration > waitingPeriod, "Duration must be greater than waiting period");
        require(isEligible(to), "Recipient is not eligible");
        require(msg.sender != to, "Sender cannot be recipient");
        
        bytes32 id = keccak256(abi.encodePacked(msg.sender, to, block.timestamp, block.number));
        require(paychans[id].status >= 0, "Channel already exists");

        require(msg.value == senderDeposit, "Incorrect sender deposit amount");

        paychans[id] = PayChan({
            id: id,
            sender: payable(msg.sender),
            recipient: payable(to),

            senderDeposit: senderDeposit,   // Sender will have to lock the fund first

            startTime: block.timestamp,
            // duration: duration,
            // endTime: block.timestamp + duration, 
            status: 1,      // Channel not yet open, waiting for the recipient to lock their funds

            fee: 0,

            disputeStartTime: 0,  // until one party called closeChan
            disputeDuration: 1000,    // default 1000s
            senderConfirm: false,
            recipientConfirm: false
        });
        emit ChannelOpened(id);
        return id;
    }


    function encodeData(uint channelID,  uint amount) public pure returns (bytes32) {
        return bytes32(abi.encodePacked(channelID, amount));
    }


    struct RequestPaymentBody {
        bytes32 ChannelId;
        uint    Amount;
        bytes   Signature;
    }
    event decodedPaymentBody (
        bytes32 ChannelId,
        uint Amount,
        bytes   Signature
    );

    function decodePaymentBody(bytes memory serializedPaymentBody) public returns (RequestPaymentBody memory) {
        RLPReader.RLPItem[] memory items = serializedPaymentBody.toRlpItem().toList();
        require(items.length == 3, "Incorrect number of fields in RLP encoded data");
        RequestPaymentBody memory reqPayment;
        reqPayment.ChannelId = bytes32(items[0].toBytes());
        reqPayment.Amount = items[1].toUint();
        reqPayment.Signature = items[2].toBytes();
        emit decodedPaymentBody(reqPayment.ChannelId, reqPayment.Amount, reqPayment.Signature);

        return reqPayment;
    }

    event msgHash (
        bytes32 msgHash
    );

    function getPaymentHash(
        bytes32 ChannelId,
        uint Amount
    ) public returns (bytes32) {
        bytes32 hash = keccak256(abi.encodePacked(
            ChannelId,
            Amount
        ));
        emit msgHash(hash);
        return keccak256(abi.encodePacked(
            ChannelId,
            Amount
        ));
    }

    event signerAdd (address signer);


    function verifyPayment(bytes32 id, uint amount, bytes memory sig, address signer) public returns (bool) {
        // RequestPaymentBody memory reqPayment;
        // reqPayment = decodePaymentBody(serializedPaymentBody);
        bytes32 resHash = getPaymentHash(id, amount);
        // emit msgHash(resHash);
        address reqSigner = ECDSA.recover(resHash, sig);
        emit signerAdd(signer);
        emit signerAdd(reqSigner);

        if (reqSigner == signer) {
            return true;
        } else {
            return false;
        }
    }

    event LogUint(uint);

    function closeChan(bytes memory serializedPaymentBody) public {
        RequestPaymentBody memory reqPayment;
        reqPayment = decodePaymentBody(serializedPaymentBody);

        bytes32 id = reqPayment.ChannelId;
        // id = 0x7282b44c0d5adacfe20f3703fcf95d9dea56e688cffeef768088761821d8e48e;
        PayChan memory channel = paychans[id];
        require(channel.recipient == msg.sender, "Only recipient");

        uint amount = reqPayment.Amount;
        emit LogUint(amount);
        require(amount <= channel.senderDeposit, "Exceed the maximum that can be retrieved.");

        bytes memory sig = reqPayment.Signature;


        require(verifyPayment(id, amount, sig, channel.sender), "The signature was wrong");
        paychans[id].fee = amount;

        channel.sender.transfer(channel.senderDeposit - paychans[id].fee);
        emit LogUint(channel.senderDeposit - paychans[id].fee);
        channel.recipient.transfer(paychans[id].fee);
        emit LogUint(paychans[id].fee);

        // Indicate that the channel is closed 
        paychans[id].status = 0;
    }



    function confirmClosure(bytes32 channelId, uint value) public payable onlyParticipants(channelId) {
        
        // Basic checks
        require(paychans[channelId].status == 2, "This channel is not under waiting-closure period.");
        require(block.timestamp < paychans[channelId].disputeStartTime + paychans[channelId].disputeDuration, "Dispute period has expired.");
        require(value == paychans[channelId].fee, "It seems there is a dispute over the final state, go report it.");
        
        PayChan memory channel = paychans[channelId];
        address payable sender = channel.sender;
        address payable recipient = channel.recipient;
        address signer;
        signer = msg.sender;

        if (channel.senderConfirm) {
            require(signer == recipient, "Only recipient can confirm closure now.");
        } else {
            require(signer == sender, "Only sender can confirm closure now.");
        }

        // TODO: Now all the token transfer has been done in function `closeChannel`
        //       Move the final settle-up here for a dispute window
        
        // settle the channel
        // since all the checks have been passed, the channel will calculate the final tokens and settle

        // uint unpaidRequest = channel.requestId - channel.paidRequests;
        // uint amount = unpaidRequest * channel.fee;

        sender.transfer(channel.senderDeposit - paychans[channelId].fee);
        recipient.transfer(paychans[channelId].fee);

        // update the status of the channel
        paychans[channelId].status = 0;     // close the channel
    }

    function balance(address addr) public view returns (uint bal) {
        return balance(addr);
    }

    function paychanCheck(bytes32 channelId) public view returns(PayChan memory) {
        return paychans[channelId];
    }

    function paychanSelectedArguments(bytes32 channelId) public view returns (address sender, address rec, uint status, uint senderB, uint fee) {
        PayChan memory channel = paychans[channelId];

        return 
            (
                channel.sender,
                channel.recipient,
                channel.status, 
                channel.senderDeposit,
                channel.fee
            );

    }
}