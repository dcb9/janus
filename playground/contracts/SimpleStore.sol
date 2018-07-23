pragma solidity ^0.4.18;
contract SimpleStore {
  event UpdateValue(
    address from,
    uint _oldValue,
    uint indexed _newValue
  );

  constructor(uint _value) public {
    value = _value;
  }

    function set(uint newValue) public {
      emit UpdateValue(msg.sender, value, newValue);
        value = newValue;
    }

    function get() public constant returns (uint) {
        return value;
    }

    uint value;
}
