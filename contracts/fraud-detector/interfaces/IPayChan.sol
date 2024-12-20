// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.8.2 <0.9.0;

interface IPaychan {
    function paychanSelectedArguments(bytes32 channelId) external view returns (
        address sender,
        address recipient,
        uint status,
        uint fee
    );
}