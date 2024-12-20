// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

struct ChannelInfo {
    address sender;
    address recipient;
    uint status;
    address witness;
}


struct PayChan {
    address sender;
    address recipient;
    uint status;
    uint fee;
}

    event emitPayChan(
        PayChan paychan
    );
