/**
 *Submitted for verification at Etherscan.io on 2020-11-18
*/

pragma solidity 0.6.8;


interface IERC20 {
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);

    function approve(address _spender, uint256 _value) external returns (bool success);

    function transfer(address _to, uint256 _value) external returns (bool success);

    function transferFrom(
        address _from,
        address _to,
        uint256 _value
    ) external returns (bool success);

    function allowance(address _owner, address _spender) external view returns (uint256 remaining);

    function balanceOf(address _owner) external view returns (uint256 balance);

    function decimals() external view returns (uint8 digits);

    function totalSupply() external view returns (uint256 supply);
}


interface IConversionRates {

    function getRate(
        IERC20 token,
        uint256 currentBlockNumber,
        bool buy,
        uint256 qty
    ) external view returns(uint256);
}


interface IKyberReserve {
    function getConversionRate(
        IERC20 src,
        IERC20 dest,
        uint256 srcQty,
        uint256 blockNumber
    ) external view returns (uint256);

    function conversionRatesContract() external view returns (IConversionRates);
    function sanityRatesContract() external view returns (IKyberSanity);
}


interface IKyberSanity {
    function getSanityRate(IERC20 src, IERC20 dest) external view returns (uint256);
}

contract Utils5 {
    IERC20 internal constant ETH_TOKEN_ADDRESS = IERC20(
        0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE
    );
}

contract KyberRateQueryHelper is Utils5 {

    function getReserveRates(IKyberReserve reserve, IERC20[] calldata srcs, IERC20[] calldata dests)
    external view returns(uint256[] memory pricingRates, uint256[] memory sanityRates)
    {
        require(srcs.length == dests.length, "srcs length != dests");

        pricingRates = new uint256[](srcs.length);
        sanityRates = new uint256[](srcs.length);

        for (uint256 i = 0 ; i < srcs.length ; i++) {

            if (reserve.sanityRatesContract() != IKyberSanity(0x0)) {
                sanityRates[i] = reserve.sanityRatesContract().getSanityRate(srcs[i], dests[i]);
            }

            pricingRates[i] = reserve.conversionRatesContract().getRate(
                srcs[i] == ETH_TOKEN_ADDRESS ? dests[i] : srcs[i],
                block.number,
                srcs[i] == ETH_TOKEN_ADDRESS ? true : false,
                0);
        }

        return (pricingRates,sanityRates);
    }
}