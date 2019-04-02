# Reserve listed tokens 

Returns list of tokens listed in KyberNetwork, includes any current and historical listed tokens. If a token has
address changed, the old addresses should be included.

## Get listed tokens 

```shell
curl -X GET "http://gateway.local/reserve/tokens?reserve=0x63825c174ab367968EC60f061753D3bbD36A0D8F"
```

> the above request will return reponse like this:

```json
{
    "version": 1,
    "block_number": 7206067,
    "data": [
        {
            "address": "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
            "symbol": "ETH",
            "name": "Ethereum",
            "timestamp": 0
        },
        {
            "address": "0xdd974d5c2e2928dea5f71b9825b8b646686bd200",
            "symbol": "KNC",
            "name": "KyberNetwork",
            "timestamp": 1505194399
        },
        {
            "address": "0x1985365e9f78359a9B6AD760e32412f4a445E862",
            "symbol": "REP",
            "name": "Augur",
            "timestamp": 1531037764,
            "old": [
                {
                    "address": "0xe94327d07fc17907b4db788e5adf2ed424addff6",
                    "timestamp": 1501261709
                }
            ]
        }
    ]
}
```

### HTTP request

`GET http://gateway.local/transactions`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
reserve | string | false | empty | reserve addresses 