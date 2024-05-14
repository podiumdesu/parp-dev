pragma solidity ^0.8.23;

contract Store {
  event ItemSet(string key, uint256 value);

  string public version;
  mapping (string => uint256) public items;

  constructor(string memory _version) public {
    version = _version;
  }

  function setItem(string memory key, uint256 value) external {
    items[key] = value;
    emit ItemSet(key, value);
  }
}
