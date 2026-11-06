pragma solidity ^0.6.6;

contract CroBridge {

    event __UpsideSendCroToIbc(address sender, string recipient, uint256 amount);

    // Pay the contract a certain Up amount and trigger a Up transfer
    // from the contract to recipient through IBC
    function send_cro_to_crypto_org(string memory recipient) public payable {
        emit __UpsideSendCroToIbc(msg.sender, recipient, msg.value);
    }
}
