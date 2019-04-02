# Reserve transaction 

## Get transactions 

```shell
curl -X GET "http://gateway.local/transactions?type=all&from=1554107172000&to=1554107172000"
```

> the above request will return reponse like this:

```json
{
    "normal": [
        {
            "blockNumber": "7481515",
            "hash": "0xd1a055853bc8cec4725470b4db4fcbab675b500e24b1b9d5b02a8be69197f7e9",
            "blockHash": "0xfabc05450cc324fa0c27368547718ff1da8170385cefdae8275e2434d4f94f78",
            "from": "0x5bab5ef16cfac98e216a229db17913454b0f9365",
            "to": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
            "value": 0,
            "gas": "108760",
            "gasUsed": "58760",
            "gasPrice": 50100000000,
            "isError": "0",
            "timestamp": 1554107172000
        }
    ],
    "internal": [
        {
            "blockNumber": "7481515",
            "hash": "0xd1a055853bc8cec4725470b4db4fcbab675b500e24b1b9d5b02a8be69197f7e9",
            "blockHash": "0xfabc05450cc324fa0c27368547718ff1da8170385cefdae8275e2434d4f94f78",
            "from": "0x5bab5ef16cfac98e216a229db17913454b0f9365",
            "to": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
            "value": 0,
            "gas": "108760",
            "gasUsed": "58760",
            "gasPrice": 50100000000,
            "isError": "0",
            "timestamp": 1554107172000
        }
    ],
    "ERC20": [
        {
            "blockNumber": "7481515",
            "hash": "0xd1a055853bc8cec4725470b4db4fcbab675b500e24b1b9d5b02a8be69197f7e9",
            "blockHash": "0xfabc05450cc324fa0c27368547718ff1da8170385cefdae8275e2434d4f94f78",
            "from": "0x5bab5ef16cfac98e216a229db17913454b0f9365",
            "to": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
            "value": 0,
            "gas": "108760",
            "gasUsed": "58760",
            "gasPrice": 50100000000,
            "isError": "0",
            "timestamp": 1554107172000
        }
    ]
}
```

### HTTP request

`GET http://gateway.local/transactions`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | true | one hour from now | from time to get transactions 
to | integer | true | now | to time to get transactions 
type | string | true | all | include: "normal", "internal", "ERC20" 