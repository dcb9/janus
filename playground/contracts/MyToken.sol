pragma solidity ^0.4.18;

import "./openzeppelin-solidity/contracts/token/ERC20/CappedToken.sol";

contract MyToken is CappedToken {
    constructor(uint256 _cap) CappedToken(_cap) public {}
}
