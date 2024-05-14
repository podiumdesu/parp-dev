// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

contract paychan {

    struct PayChan {
        bytes32 id; // store the payment channel id

        address payable sender; // store the sender
        address payable recipient; // store the recipient

        // funds
        uint senderDeposit;

        // start time of the channel
        uint startTime;

        uint status;     // 1 = open, 0 = closed, 2 = waitingClosure/Closing

        uint oweAmount;     // the amount of money owed by the sender

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

    function openChan(address to, uint duration, uint senderDeposit) public payable returns(bytes32 channelId) {
        // require(duration > waitingPeriod, "Duration must be greater than waiting period");
        require(msg.sender != to, "Sender cannot be recipient");
        
        bytes32 id = keccak256(abi.encodePacked(msg.sender, to));
        require(paychans[id].status >= 0, "Channel already exists");

        require(msg.value == senderDeposit, "Incorrect sender deposit amount");

        paychans[id] = PayChan({
            id: id,
            sender: payable(msg.sender),
            recipient: payable(to),

            senderDeposit: senderDeposit,   // Sender will have to lock the fund first
            oweAmount: 0,

            startTime: block.timestamp,
            status: 1,      // Channel is open

            disputeStartTime: 0,  // until one party called closeChan
            disputeDuration: 1000,    // default 1000s
            senderConfirm: false,
            recipientConfirm: false
        });

        return id;
    }


    // TODO: turn it to be a multisig smart contract
    // Therefore, it can reconstruct the channel id

    // When channel is open, it can be closed by either party with a correct signature
    // function closeChan(bytes32 channelId, bytes32 h, uint8 v, bytes32 r, bytes32 s, uint value) public onlyWhenChannelOpen(channelId) {

    function closeChan(bytes32 channelId, uint value) public onlyWhenChannelOpen(channelId) {
        require(paychans[channelId].oweAmount < value, "Invalid value, smaller than the previous record");
        // value is the requests sent in this channel

        address signer;
        bytes32 proof;
        // signer = ecrecover(h, v, r, s);
        signer = msg.sender;

        // check the signature
        require((signer == paychans[channelId].sender) || (signer == paychans[channelId].recipient), "Not authorized signer");
        
        proof = keccak256(abi.encodePacked(value));

        // require(proof == h, "Not correct value");

        if (signer == paychans[channelId].sender) {
            paychans[channelId].senderConfirm = true;
        } else {
            paychans[channelId].recipientConfirm = true;
        }

        // Record the final state of the channel
        paychans[channelId].oweAmount = value;
        // Wait until the remaining time
        paychans[channelId].status = 2; // set the status to "under dispute period"
        paychans[channelId].disputeStartTime = block.timestamp;

        // TODO: verify the signature

    }

    // confirmClosure(bytes32 channelId, bytes32 h, uint8 v, bytes32 r, bytes32 s, uint value)
    function confirmClosure(bytes32 channelId, uint value) public payable onlyParticipants(channelId) {
        
        // Basic checks
        require(paychans[channelId].status == 2, "This channel is not under waiting-closure period.");
        require(block.timestamp < paychans[channelId].disputeStartTime + paychans[channelId].disputeDuration, "Dispute period has expired.");
        require(value == paychans[channelId].oweAmount, "It seems there is a dispute over the final state, go report it.");
        
        PayChan memory channel = paychans[channelId];
        address payable sender = channel.sender;
        address payable recipient = channel.recipient;
        address signer;
        // bytes32 proof;
        // TODO: Uncomment
        // signer = ecrecover(h, v, r, s);
        // require(signer == msg.sender);
        signer = msg.sender;

        if (channel.senderConfirm) {
            require(signer == recipient, "Only recipient can confirm closure now.");
        } else {
            require(signer == sender, "Only sender can confirm closure now.");
        }

        // proof = keccak256(abi.encodePacked(value));
        // require(proof == h, "Not correct value");

        // Settle the channel
        sender.transfer(channel.senderDeposit - value);
        recipient.transfer(value);

        // update the status of the channel
        paychans[channelId].status = 0;     // close the channel
    }

    function balance(address addr) public view {
        return balance(addr);
    }

    function paychanCheck(bytes32 channelId) public view returns(PayChan memory) {
        return paychans[channelId];
    }

    function paychanSelectedArguments(bytes32 channelId) public view returns (uint status, uint senderB, uint amt) {
        PayChan memory channel = paychans[channelId];

        return 
            (
                channel.status, 
                channel.senderDeposit,
                channel.oweAmount
            );
    }

}
