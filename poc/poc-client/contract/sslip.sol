// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

contract paychan {

    event ChannelOpened(bytes32 indexed channelId);
    
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
        require(msg.sender != to, "Sender cannot be recipient");
        
        bytes32 id = keccak256(abi.encodePacked(msg.sender, to));
        // TODO: Wrong code 
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

    // function recipientLockFunds(bytes32 channelId) external payable onlyParticipants(channelId) {
    //     require(msg.sender == paychans[channelId].recipient, "Only recipient");
    //     require(block.timestamp < paychans[channelId].endTime, "Channel has expired");
    //     require(block.timestamp < paychans[channelId].startTime + paychans[channelId].waitingPeriod, "Too late to lock funds");
    //     require(msg.value == paychans[channelId].recipientDeposit, "Incorrect recipient deposit amount");

    //     paychans[channelId].status = 1; // set the channel open
    // }


    // TODO: turn it to be a multisig smart contract
    // Therefore, it can reconstruct the channel id

    // When channel is open, it can be closed by either party with a correct signature
    // function closeChan(bytes32 channelId, bytes32 h, uint8 v, bytes32 r, bytes32 s, uint value) public onlyWhenChannelOpen(channelId) {

    function closeChan(bytes32 channelId, uint value, bytes32 hash, bytes memory signature) public onlyWhenChannelOpen(channelId) {
        require(paychans[channelId].fee <= value, "Invalid value, smaller than the previous record");
        // value is the requests sent in this channel

        // bytes32 proof;
        // signer = ecrecover(h, v, r, s);

        
        // check the signature

        (uint8 v, bytes32 r, bytes32 s) = splitSignature(signature);
        address signer = ecrecover(hash, v, r, s);

        require((signer == paychans[channelId].sender), "Not authorized signer");

        // proof = keccak256(abi.encodePacked(value));

        // Record the final state of the channel
        paychans[channelId].fee = value;
        // Wait until the remaining time
        paychans[channelId].status = 2; // set the status to "under dispute period"
        paychans[channelId].disputeStartTime = block.timestamp;


    }

    function splitSignature(bytes memory signature) public pure returns (uint8 v, bytes32 r, bytes32 s) {
        require(signature.length == 65, "Invalid signature length");

        // Extract v, r, and s from the signature
        assembly {
        r := mload(add(signature, 32))
        s := mload(add(signature, 64))
        v := byte(0, mload(add(signature, 65)))
        }
    }

    // confirmClosure(bytes32 channelId, bytes32 h, uint8 v, bytes32 r, bytes32 s, uint value)
    function confirmClosure(bytes32 channelId, uint value) public payable onlyParticipants(channelId) {
        
        // Basic checks
        require(paychans[channelId].status == 2, "This channel is not under waiting-closure period.");
        require(block.timestamp < paychans[channelId].disputeStartTime + paychans[channelId].disputeDuration, "Dispute period has expired.");
        require(value == paychans[channelId].fee, "It seems there is a dispute over the final state, go report it.");
        
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


        // settle the channel
        // since all the checks have been passed, the channel will calculate the final tokens and settle

        // uint unpaidRequest = channel.requestId - channel.paidRequests;
        // uint amount = unpaidRequest * channel.fee;

        sender.transfer(channel.senderDeposit - paychans[channelId].fee);
        recipient.transfer(paychans[channelId].fee);

        // update the status of the channel
        paychans[channelId].status = 0;     // close the channel
    }

    // TODO: multisig check
    // function closeChan(bytes32 channelId, bytes32 h, uint8 v, bytes32 r, bytes32 s, uint value) public onlyWhenChannelOpen(channelId) {

    // function updateState(bytes32 channelId, uint value) public payable onlyParticipants(channelId) onlyWhenChannelOpen(channelId){
    //     // Verify the signature (if it is from both parties)
    //     require(value > paychans[channelId].fee, "Invalid value");

    //     PayChan memory channel = paychans[channelId];

    //     // update the state
    //     // The value will be recorded

    //     // Use deposit to pay the fee
    //     uint unpaidRequests = value - channel.paidRequests;
    //     paychans[channelId].senderDeposit = channel.senderDeposit - (unpaidRequests * channel.fee);
    //     paychans[channelId].recipientDeposit = channel.recipientDeposit + (unpaidRequests * channel.fee);
        
    //     paychans[channelId].requestId = value;
    //     paychans[channelId].paidRequests = value;

    // }


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
    // function disputeResolve() public {

    // }

    function greeting(address from) public pure returns (string memory) {
        address fromAddress = from; 
        return "hello";
    }
}