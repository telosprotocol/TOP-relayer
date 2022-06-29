// SPDX-License-Identifier: GPL-3.0-or-later
pragma solidity ^0.8.0;

contract topbridge {
    function sync(bytes memory blockHeader) public returns (bool success) {
        bytes memory payload = abi.encodeWithSignature("sync(bytes)", blockHeader);
        (success,) = msg.sender.call(payload);
        require (success);
    }

    function get_height() public returns (uint64 height) {
        bytes memory payload = abi.encodeWithSignature("get_height()");
        bool success = false;
        bytes memory returnData;
        (success, returnData) = msg.sender.call(payload);
        require(returnData.length >= height + 32);
        assembly {
            height := mload(add(returnData, 0x20))
        }
        require (success);
    }
}