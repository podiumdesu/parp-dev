// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

contract DepositContract {

    mapping(address => uint256) public deposits;

    event Deposited(address indexed user, uint256 amount);

    event Slashed(address indexed user, uint256 amount, uint256 lightClient, uint256 Witness, uint256 burned);

    function deposit() external payable {
        require(msg.value >= 1 ether, "Minimum deposit is 1 Ether");

        deposits[msg.sender] += msg.value;

        emit Deposited(msg.sender, msg.value);
    }

    function getDeposit(address user) external view returns (uint256) {
        return deposits[user];
    }


    function slash(address user, address lc, address witness) external {
        uint256 userDeposit = deposits[user];
        require(userDeposit > 0, "User has no deposit");

        uint256 slashAmount = userDeposit / 2; // Halve the user's deposit
        uint256 oneThird = slashAmount / 3;

        // Deduct the slashed amount from the user's deposit
        deposits[user] -= slashAmount;

        // Transfer 1/3 to light client
        (bool successLc, ) = payable(lc).call{value: oneThird}("");
        require(successLc, "Transfer to light client failed");

        // Transfer 1/3 to witness node
        (bool successWitness, ) = payable(witness).call{value: oneThird}("");
        require(successWitness, "Transfer to witness failed");

        // The remaining 1/3 is burned (left unallocated in the contract)
        emit Slashed(user, slashAmount, oneThird, oneThird, oneThird);
    }
    
    function withdraw(address payable user) external {
        uint256 amount = deposits[user];
        deposits[user] = 0;
        user.transfer(amount);
    }
}