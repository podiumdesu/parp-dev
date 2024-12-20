// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.8.2 <0.9.0;

interface IDepositContract {
    function slash(address user, address lc, address witness) external;
    function getDeposit(address user) external view returns (uint256);
}

