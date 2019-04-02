## Get withdrawal history 

```shell
curl -X GET "http://gateway.local/withdrawals?from=1494901162000&to=1499865549600?cex=all"
```

> the above request will return reponse like this:

```json
{
    "huobi": [
        {
            "id": 2272335,
            "created-at": 1525754125590,
            "updated-at": 1525754753403,
            "currency": "ETH",
            "type": "withdraw",
            "amount": 0.48957444,
            "state": "confirmed",
            "fee": 0.01,
            "address": "f6a605cdd9b2471ffdff706f8b7665a12b862158",
            "address-tag": "",
            "tx-hash": "cdef3adad017d9564e62282f5e0f0d87d72b995759f1f7f4e473137cc1b96e56"
        }
    ],
    "binance": [
        {
            "id": "7213fea8e94b4a5593d507237e5a555b",
            "amount": 1,
            "address": "0x6915f16f8791d0a1cc2bf47c13a6b2a92000504b",
            "asset": "ETH",
            "txId": "0xdf33b22bdb2b28b1f75ccd201a4a4m6e7g83jy5fc5d5a9d1340961598cfcb0a1",
            "applyTime": 1525754125591,
            "status": 4
        },
        {
            "id": "7213fea8e94b4a5534ggsd237e5a555b",
            "amount": 1000,
            "Address": "463tWEBn5XZJSxLU34r6g7h8jtxuNcDbjLSjkn3XAXHCbLrTTErJrBWYgHJQyrCwkNgYvyV3z8zctJLPCZy24jvb3NiTcTJ",
            "Asset": "XMR",
            "txId": "b3c6219639c8ae3f9cf010cdc24fw7f7yt8j1e063f9b4bd1a05cb44c4b6e2509",
            "applyTime": 1525754125592,
            "status": 4
        }
    ]
}
```

### HTTP request

`GET http://gateway.local/withdrawals`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | from time to get trades
to | integer | false | now | to time to get trades
cex | string | false | all | valid value: "binance", "huobi"