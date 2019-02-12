## Summary

This project provides APIs to access KyberNetwork's reserve related data:

- KyberNetwork's reserve Ethereum addresses
  - reserve
  - pricing operator
  - centralized exchange deposit
- Ethereum transactions of KyberNetwork reserve's addresses
- ERC20 transactions of KyberNetwork reserve's addresses
- Transactions in centralized exchanges (binance, huobi)

## General Conventions

1. For POST/PUT/DELETE requests, the content type is set to: application/json
2. Timestamp format in params: milliseconds since Unix epoch
3. Float amount precision: 6 decimals

## Reserve Addresses

Returns Ethereum [addresses](https://developer.kyber.network/docs/MainnetEnvGuide/) related to KyberNetwork's reserve.
An address will have one of following types:

- reserve: address of KyberNetwork's reserve contract, example: 0x63825c174ab367968EC60f061753D3bbD36A0D8F

- pricing operator: operator addresses of conversion rate contract. Can be get by calling `getOperators` method of [conversion rate contract](https://etherscan.io/address/0x798abda6cc246d0edba912092a2a3dbd3d11191b#readContract), example: 0x8bc3da587def887b5c822105729ee1d6af05a5ca, 0x9224016462b204c57eb70e1d69652f60bcaf53a8

- centralized exchange deposit addresses: Ethereum address to deposit funds to centralized exchanges (binance: 0x44d34a119ba21a42167ff8b77a88f0fc7bb2db90, huobi: 0x0c8fd73eaf6089ef1b91231d0a07d0d2ca2b9d66)

The address list will be maintain by an operator, and is append only. Every time a change is made to the address list, the version is increased by one.

### API Endpoints

#### GET reserve/addresses

Returns all configured addresses.

Example response:

```json
{
  "version": 1,
  "data": [
    {
      "id": 1,
      "address": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
      "type": "reserve",
      "description": "main reserve"
    },
    {
      "id": 2,
      "address": "0x8bC3da587DeF887B5C822105729ee1D6aF05A5ca",
      "type": "operator",
      "description": "pricing operator 1"
    },
    {
      "id": 3,
      "address": "0x9224016462b204c57eb70e1d69652f60bcaf53a8",
      "type": "operator",
      "description": "pricing operator 2"
    },
    {
      "id": 4,
      "address": "0x44d34a119ba21a42167ff8b77a88f0fc7bb2db90",
      "type": "exchange",
      "description": "binance deposit"
    },
    {
      "id": 5,
      "address": "0x0c8fd73eaf6089ef1b91231d0a07d0d2ca2b9d66",
      "type": "exchange",
      "description": "huobi deposit"
    }
  ]
}
```

#### POST reserve/addresses

Append given address to the list, increase version after that.

Example request:

```json
{
  "address": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
  "type": "reserve",
  "description": "main reserve"
}
```

Fields:

- address: required, valid ethereum address, unique
- type: required, matches one of predefined types above
- description: required, unique

Example response:

```json
{
  "id": 20
}

```

Fields:

- id: id of newly created address, to be used in PUT API below

#### PUT reserve/addresses

Update information of an existing address.

Example request.

```json
{
  "id": 20,
  "address": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
  "type": "reserve",
  "description": "main reserve"
}
```

Fields:

- id: required
- address: optional
- type: optional
- description: optional

## Reserve Tokens

Returns list of tokens listed in KyberNetwork, includes any current and historical listed tokens. If a token has 
address changed, the old addresses should be included.

### Fetchers

#### Listed Tokens Fetcher

The service will be run periodically with current block number, result will be stored in database as a append only list.

Input: 

- reserve contract [address](https://etherscan.io/address/0x63825c174ab367968ec60f061753d3bbd36a0d8f)

Output: 

- listed tokens at the current time

Steps: 

- Getting [conversion rate contract](https://etherscan.io/address/0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B#readContract) of given 
reserve contract address by calling `conversionRatesContract` method

- Calling `getListedTokens` method of reserve rate contract to get listed token addresses
  - for each address, extracted following data by calling ERC20 methods
    - "symbol" constant of ERC20 token
    - "name" constant of ERC20 token
- Any time a new token is found, version is increased 1


### API Endpoints

#### GET reserve/tokens

Example response:

```json
{
  "version": 1,
  "block_number": 7206067,
  "data": [
    {
      "address": "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
      "symbol": "ETH",
      "name": "Ethereum"
    },
    {
      "address": "0xdd974d5c2e2928dea5f71b9825b8b646686bd200",
      "symbol": "KNC",
      "name": "KyberNetwork"
    },
    {
      "address": "0x1985365e9f78359a9B6AD760e32412f4a445E862",
      "symbol": "REP",
      "name": "Augur",
      "old_addresses": ["0xe94327d07fc17907b4db788e5adf2ed424addff6"]
    }
  ]
}
```

Fields:

- version: version of token listing
- block_number: last block number where a new token is found
- data
  - address
  - symbol
  - name
  - old_addresses: old address of a token. Two tokens are considered the same if both symbol and name are equal  (there 
  might be additional rules), newer token is current address, older one is added to "old_addresses" field. 
  
## Reserve Transactions

List of Ethereum transactions for configured addresses.
 
### Fetchers
 
#### Pricing Update Transactions Fetcher

Input: pricing operator addresses

Output: transactions of input addresses with following function names:

- setCompactData
- setBaseRate

Extracted data:

- from: address of pricing operator
- to: address of conversion rates contract
- tx_hash: transaction hash
- timestamp: block timestamp
- gas_fee: gas fee of transaction

Example for transaction: 0xdb49f5ff9241653024e10859ca51e086c5b41653f72d166b78eb5982ee780c1c:

- from: 0x8bC3da587DeF887B5C822105729ee1D6aF05A5ca
- to: 0x798abda6cc246d0edba912092a2a3dbd3d11191b
- tx_hash: 0xdb49f5ff9241653024e10859ca51e086c5b41653f72d166b78eb5982ee780c1c
- timestamp: 1548750716
- gas_fee: 1548750716

#### Trade Logs Fetcher

Input: 

- reserve address
- from block / to block

Output: 

- trade logs of given input reserve address from given block range

Relevant events:

1. `TradeExecute` event:

    ```
    TradeExecute(
        index_topic_1 address origin, 
        address src, 
        uint256 srcAmount, 
        address destToken, 
        uint256 destAmount, 
        address destAddress)
    ```

2. ERC 20 `Transfer` event:

    ```
    Transfer(address indexed from, address indexed to, uint tokens)
    ```

A trade is always contains at least one `TradeExecute` and one `Transfer` event. A single transaction might contain 
multiple trades.

1. ETH --> ERC20 token
    - `Transfer` with:
      - emitter: token contract address
      - from: reserve contract address
      - to: KyberNetwork contract address
    - `TradeExecute` with:
      - emitter: reserve contract address
      - src: 0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee
      - destToken: token contract address
2. ERC20 token --> ETH
    - `Transfer` with
      - emitter: token contract address
      - from: KyberNetwork contract address
      - to: reserve contract address
    - `TradeExecute`
      - emitter: reserve contract address
      - src: token contract address
      - destToken: 0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee
3. ERC20 token1 --> ERC20 token2
    - `Transfer` with
      - emitter: token1 contract address
      - from: KyberNetwork contract address
      - to: reserve1 contract address
    - `TradeExecute`
      - emitter: reserve contract address
      - src: token1 contract address
      - destToken: 0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee
    - `Transfer` (optional) with
      - emitter: token1 contract address
      - from: reserve contract address
      - to: reserve2 contract address
    - `TradeExecute` with:
      - emitter: reserve2 contract address
      - src: 0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee
      - destToken: token2 contract address
    - `Transfer` with:
      - emitter: token2 contract address
      - from: KyberNetwork contract
      - to: token2 receiver address
      
    Note: in case reserve1=reserve2, the middle `Transfer` event is omitted.    

Steps:

- listen to `TradeExecute` event of reserve address, record the log index
- get transaction_hash, get transaction receipt
- get attached ERC20 `Transfer` event according to list of events above


Example for transaction: 0x1dfa7900b8ddd0dab58daa8a5f2c6aa1f0be2f7a460095ea797a37d150ac6002 (ETH --> REP):

- `Transfer` event
  - from: 0x26e7a08c746f1bccab8a6e76d6e0ed6905c1813a
  - to: 0x63825c174ab367968ec60f061753d3bbd36a0d8f
  - token: 0x1985365e9f78359a9B6AD760e32412f4a445E862
  - amount: 20.136289
  - tx_hash: 0x1dfa7900b8ddd0dab58daa8a5f2c6aa1f0be2f7a460095ea797a37d150ac6002
  - timestamp: 1548750904
  - gas_fee: 0.002617
- `TradeExecute` event
  - from: 0x63825c174ab367968ec60f061753d3bbd36a0d8f
  - to: 0x26e7a08c746f1bccab8a6e76d6e0ed6905c1813a
  - token: 0x1985365e9f78359a9B6AD760e32412f4a445E862
  - amount: 20.136289
  - tx_hash: 0x1dfa7900b8ddd0dab58daa8a5f2c6aa1f0be2f7a460095ea797a37d150ac6002
  - timestamp: 1548750904
  - gas_fee: 0.002617
  
### API Endpoints

#### GET reserve/transactions?from_time=15XXXXXXXX&to_time=15XXXXXXXX?type=transaction_type

Valid values for transaction_type:

- all (default)
- pricing
- trade

Example response:

```json
{
  "data": [
    {
      "type": "pricing",
      "from": "0x8bC3da587DeF887B5C822105729ee1D6aF05A5ca",
      "to": "0x798abda6cc246d0edba912092a2a3dbd3d11191b",
      "tx_hash": "0xdb49f5ff9241653024e10859ca51e086c5b41653f72d166b78eb5982ee780c1c",
      "timestamp": 1548750716,
      "gas_fee": 0.000434
    },
    {
      "type": "trade",
      "amount": 20.136289,
      "token": "0x1985365e9f78359a9B6AD760e32412f4a445E862",
      "tx_hash": "0x1dfa7900b8ddd0dab58daa8a5f2c6aa1f0be2f7a460095ea797a37d150ac6002",
      "timestamp": 1548750904,
      "gas_fee": 0.002617,
      "transactions": [
        {
          "from": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
          "to": "0x26e7a08c746f1bccab8a6e76d6e0ed6905c1813a",
          "amount": 20.136289,
          "token": "0x1985365e9f78359a9B6AD760e32412f4a445E862"
        },
        {
          "from": "0x26e7a08c746f1bccab8a6e76d6e0ed6905c1813a",
          "to": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
          "amount": 2.321554,
          "token": "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
        }
      ]
    }
  ]
}
```

## CEX Transactions

List of historical trade/withdrawal from centralized exchanges. Supported exchanges:

- binance
- huobi

### Fetchers

#### CEX Trades Fetchers

Returns list of trades for a specific account. The services should returns the data from Binance/Huobi as-is without 
any normalization effort. The handling logic should be implemented in consumer side. 

Input:

- API key
- from time / to time

Output:

- list of historical trades from given account and time range 

1. Binance

    API Endpoint: [GET /api/v3/myTrades](https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md#account-trade-list-user_data)

    The request parameter `symbol` is mandatory, so there is no way to get historical trades in a single API calls. 
    We will have to brute-force the Binance API to make sure no trade is missing.
    
    - get all available tradeing pairs using [GET /v1/exchangeInfo](https://api.binance.com//api/v1/exchangeInfo) API
    - for each pair, calling [GET /api/v3/myTrades](https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md#account-trade-list-user_data) multiple times
      - first time with 
        - `startTime`: given from time
        - `endTime`: given to time
      - next time with
        - `startTime`: timestamp of last trade in result
        - `endTime`: given to time
      - ... until the API call returns empty result
      
      The reason for calling multiple times is [GET /api/v3/myTrades](https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md#account-trade-list-user_data)
      only returns maximum number of trades (configurable by limit parameter) and ignore newer ones if exceed. 
       
    There should be a resource limiter with configurable maxiumum weight limit to make sure we don't hit Binance API 
    rate limiter. We can afford slow fetching time, accuracy over speed.
          
2. Huobi

    API Endpoint: [GET /v1/order/orders](https://github.com/huobiapi/API_Docs_en/wiki/REST_Reference#get-v1orderorders--get-order-list)
    
    Huobi API has the same quirks as Binance, and we need to apply same strategy, using [GET v1/common/symbols](http://api.huobi.pro/v1/common/symbols)
    for fetching all trading pairs.
    
#### CEX Withdrawal Fetchers

Returns list of historical competed withdrawal from centralized exchanges.

Input:

- API key
- from time / to time

Output:

- list of withdrawal trades from given account and time range

1. Binance

    API Endpoint: [GET /wapi/v3/withdrawHistory.html](https://github.com/binance-exchange/binance-official-api-docs/blob/master/wapi-api.md#withdraw-history-user_data).
    
    Params:
    
    - status: 6 (completed) 
    
2. Huobi

    API Endpoint: GET https://www.hbg.com/-/x/pro/v1/query/finances
    
    Notice: this is an unofficial API.
    
#### API Endpoints

##### GET /reserve/cex_trades?from_time=15XXXXXXXX&to_time=15XXXXXXXX?cex=xxx

Valid values for cex parameters:

- all (default)
- binance
- huobi

Example response:

```json
{
  "data": {
    "binance": [
      {
        "symbol": "BNBBTC",
        "id": 28457,
        "orderId": 100234,
        "price": "4.00000100",
        "qty": "12.00000000",
        "commission": "10.10000000",
        "commissionAsset": "BNB",
        "time": 1499865549590,
        "isBuyer": true,
        "isMaker": false,
        "isBestMatch": true
      }
    ],
    "huobi": [
      {
        "id": 59378,
        "symbol": "ethusdt",
        "account-id": 100009,
        "amount": "10.1000000000",
        "price": "100.1000000000",
        "created-at": 1494901162595,
        "type": "buy-limit",
        "field-amount": "10.1000000000",
        "field-cash-amount": "1011.0100000000",
        "field-fees": "0.0202000000",
        "finished-at": 1494901400468,
        "user-id": 1000,
        "source": "api",
        "state": "filled",
        "canceled-at": 0,
        "exchange": "huobi",
        "batch": ""
      }
    ]
  }
}
```

##### GET /reserve/withdrawals?from_time=15XXXXXXXX&to_time=15XXXXXXXX?cex=xxx

Valid values for cex parameters:

- all (default)
- binance
- huobi

Example response:

```json
{
  "data": {
    "binance": [
      [
        {
          "id": "7213fea8e94b4a5593d507237e5a555b",
          "amount": 1,
          "address": "0x6915f16f8791d0a1cc2bf47c13a6b2a92000504b",
          "asset": "ETH",
          "txId": "0xdf33b22bdb2b28b1f75ccd201a4a4m6e7g83jy5fc5d5a9d1340961598cfcb0a1",
          "applyTime": 1508198532000,
          "status": 4
        },
        {
          "id": "7213fea8e94b4a5534ggsd237e5a555b",
          "amount": 1000,
          "address": "463tWEBn5XZJSxLU34r6g7h8jtxuNcDbjLSjkn3XAXHCbLrTTErJrBWYgHJQyrCwkNgYvyV3z8zctJLPCZy24jvb3NiTcTJ",
          "addressTag": "342341222",
          "txId": "b3c6219639c8ae3f9cf010cdc24fw7f7yt8j1e063f9b4bd1a05cb44c4b6e2509",
          "asset": "XMR",
          "applyTime": 1508198532000,
          "status": 4
        }
      ]
    ],
    "houbi": [
      {
        "id": 123,
        "transaction-id": 456,
        "created-at": 1549876997720,
        "updated-at": 1549877232790,
        "candidate-currency": null,
        "currency": "bix",
        "type": "withdraw-virtual",
        "direction": "out",
        "amount": "277.529421600000000000",
        "state": "confirmed",
        "fees": "3.000000000000000000",
        "error-code": "",
        "error-msg": "",
        "to-address": "0x63825c174ab367968ec60f061753d3bbd36a0d8",
        "to-addr-tag": "",
        "tx-hash": "0x9806d75ad08cc8d1d009d016d978619b7bce979c96df4f01a2ddf000f5facb4b",
        "chain": "bix",
        "extra": null
      }
    ]
  }
}
```

## Reserve Rates

Historical rates of:

- ERC20 tokens --> ETH

  KNC/ETH rate means how much KNC I have to sell to get 1 ETH. So if KNC/ETH  bid is 0.0166667 on Kyber, 
  KC/ETH rate is 600

- ETH --> USD

  ETH/USD rate means how much ETH I have to sell to get one USD. So if ETH/USD bid is 100 on coinbase, 
  ETH/USD rate is 0.01

The rates will be fetched daily at the time of last pricing update event of reserve on this day using GMT timezone. 
For example, if last pricing update transaction is mined at 2018-10-14 23:59:50, rates of 2018-10-14 will be rates 
at that block.

### Fetchers

#### Reserve Rates Fetcher

Input:

- date

Output:

- ERC20 token/ETH, ETH/USD rate of given day

From data of Pricing Update Transactions Fetcher, find out the last pricing update event of given date, 
fetching ERC20 token/ETH using `getExpectedRates` method of wrapper contract. The ETH/USD rate will be fetch at the 
same timestamp using CoinGecko API.

### API Endpoints

#### GET reserve/rates?from_time=1539558000&to_time=1539565200


Example responses:

```json
{
  "data": [
    {
      "date": "2018-10-15",
      "ETH": {
        "KNC": 917.431192,
        "ZIL": 5205.351102
      },
      "USD": {
        "ETH": 0.009434
      }
    }
  ]
}
```

## Company Wallet ERC20 Transfers

There will be a preconfigured list of wallet addresses and ERC20 token addresses, returns lift of ERC20 transfers from 
given token addresses that have either source/dest in one of wallet addresses.

### Fetchers

#### Company Wallet ERC20 Transfers Fetcher

Input:

- list of wallet addresses
- list of ERC20 tokens
- from block/to block

Output:

- list of ERC20 transactions

For each token address, listen for `Transfer(address indexed from, address indexed to, uint tokens)` events, filter out
events that have from or to param in given wallet addresses

### API Endpoints

#### wallet/transactions?address=0x3eb01b3391ea15ce752d01cf3d3f09dec596f650&token=0xxxx&from_time=15XXXXXXXX&to_time=15XXXXXXXX

The API should support passing `address`, `token` parameters multiple times in request.

Example response:


```json
{
  "data": [
    {
      "from": "0x3eb01b3391ea15ce752d01cf3d3f09dec596f650",
      "to": "0x8180a5ca4e3b94045e05a9313777955f7518d757",
      "amount": 61875,
      "currency": "0xc86d054809623432210c107af2e3f619dcfbf652",
      "tx_hash": "0xd6afca082a87b9f179d9d6e4bc2657cc90d08e57333cf4418732ac27c95da7ef",
      "timestamp": 1547889183,
      "gas_fee": 0.000792,
      "currency_fee": 0
    }
  ]
}
```