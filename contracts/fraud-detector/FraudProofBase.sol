// SPDX-License-Identifier: Apache2
pragma solidity ^0.8.17;

import "./interfaces/IPayChan.sol";
import "./interfaces/IDeposit.sol";

abstract contract FraudProofBase {
    IPaychan public paychanContract;
    IDepositContract public depositContract;
    
    constructor(address _paychanContractAddress, address _depositContract) {
        paychanContract = IPaychan(_paychanContractAddress);
        depositContract = IDepositContract(_depositContract);
    }
}