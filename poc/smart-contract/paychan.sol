// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

contract paychan {

    struct PayChan {
        bytes32 id; // store the payment channel id

        address payable sender; // store the sender
        address payable recipient; // store the recipient

        // funds
        uint senderDeposit;
        // uint recipientDeposit;
        // uint waitingPeriod; // waiting period for the recipient to lock their funds, in seconds

        // expiry date of the channel
        uint startTime;
        // uint duration;
        // uint endTime;
        uint status;     // 1 = open, 0 = closed, 2 = waitingClosure/Closing


        uint oweAmount;     // the amount of money owed by the sender

        // state that will be updated
        // uint requestId;
        // uint paidRequests;
        // uint fee;

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
            // recipientDeposit: recipientDeposit,   // How much money should the recipient lock 
            // waitingPeriod: waitingPeriod,
            // paidRequests: 0,
            oweAmount: 0,

            startTime: block.timestamp,
            status: 1,      // Channel not yet open, waiting for the recipient to lock their funds

            // requestId: 0,  // 0 = no request
            // fee: 1,
            disputeStartTime: 0,  // until one party called closeChan
            disputeDuration: 1000,    // default 1000s
            senderConfirm: false,
            recipientConfirm: false
        });

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

    function closeChan(bytes32 channelId, uint value) public onlyWhenChannelOpen(channelId) {
        // require(paychans[channelId].requestId <= value, "Invalid value, smaller than the previous record");
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
        paychans[channelId].requestId = value;
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
        require(value == paychans[channelId].requestId, "It seems there is a dispute over the final state, go report it.");
        
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

        uint unpaidRequest = channel.requestId - channel.paidRequests;
        uint amount = unpaidRequest * channel.fee;

        sender.transfer(channel.senderDeposit - amount);
        recipient.transfer(channel.recipientDeposit + amount);

        // update the status of the channel
        paychans[channelId].status = 0;     // close the channel
    }

    // TODO: multisig check
    // function closeChan(bytes32 channelId, bytes32 h, uint8 v, bytes32 r, bytes32 s, uint value) public onlyWhenChannelOpen(channelId) {

    function updateState(bytes32 channelId, uint value) public payable onlyParticipants(channelId) onlyWhenChannelOpen(channelId){
        // Verify the signature (if it is from both parties)
        require(value > paychans[channelId].requestId, "Invalid value");

        PayChan memory channel = paychans[channelId];

        // update the state
        // The value will be recorded

        // Use deposit to pay the fee
        uint unpaidRequests = value - channel.paidRequests;
        paychans[channelId].senderDeposit = channel.senderDeposit - (unpaidRequests * channel.fee);
        paychans[channelId].recipientDeposit = channel.recipientDeposit + (unpaidRequests * channel.fee);
        
        paychans[channelId].requestId = value;
        paychans[channelId].paidRequests = value;

    }


    function balance(address addr) public view {
        return balance(addr);
    }

    function paychanCheck(bytes32 channelId) public view returns(PayChan memory) {
        return paychans[channelId];
    }

    function paychanSelectedArguments(bytes32 channelId) public view returns (uint status, uint senderB, uint recB, uint requestId, uint paidRequest) {
        PayChan memory channel = paychans[channelId];

        return 
            (
                channel.status, 
                channel.senderDeposit,
                channel.recipientDeposit,
                channel.requestId,
                channel.paidRequests
            );

    }
    // function disputeResolve() public {

    // }
}
