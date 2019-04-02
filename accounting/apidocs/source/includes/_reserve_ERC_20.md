# Wallet ERC20 transfer transactions 

There will be a preconfigured list of wallet addresses and ERC20 token addresses, returns lift of ERC20 transfers from
given token addresses that have either source/dest in one of wallet addresses.

## Get wallet ERC20 transfer transactions

```shell
curl -X GET "http://gateway.local/wallet/transactions?wallet=0x63825c174ab367968EC60f061753D3bbD36A0D8F&token=0xd6Cd31F283d24cfb442cBA1Bcf42290c07C15792"
```

> the above request will return reponse like this:

```json
[
    {
        "timestamp": 1554110762000,
        "hash": "0x531058ad47191a75ab511b7f16aff1dd8f22c19489d89439c522049b730b33d8",
        "from": "0x9ae49C0d7F8F9EF4B864e004FE86Ac8294E20950",
        "contractAddress": "0xdd974D5C2e2928deA5F71b9825b8b646686BD200",
        "to": "0x63825c174ab367968EC60f061753D3bbD36A0D8F",
        "blockNumber": "7481765",
        "value": 141640000000000000000,
        "gas": "700000",
        "gasUsed": "242256",
        "gasPrice": 10000000000
    },
    {
        "timestamp": 1554111480000,
        "hash": "0x6876c84dc6c5f29a0ced835fe694382363d0fef8ec14cc28f68d3a160b34af4e",
        "from": "0x63825c174ab367968EC60f061753D3bbD36A0D8F",
        "contractAddress": "0xdd974D5C2e2928deA5F71b9825b8b646686BD200",
        "to": "0x9ae49C0d7F8F9EF4B864e004FE86Ac8294E20950",
        "blockNumber": "7481822",
        "value": 20024984649382703000,
        "gas": "380000",
        "gasUsed": "239049",
        "gasPrice": 11600000000
    }
]
```

### HTTP request

`GET http://gateway.local/transactions`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | from time to get transactions 
to | integer | false | now | to time to get transactions 
wallet | string | true | empty | wallet address to get transaction
token | string | true | empty | token address to get transactions for 