# Centralized exchanges (cex) endpoint

## Get trades history 

```shell
curl -X GET "http://gateway.local/trades?from=1494901162000&to=1499865549600&cex=all"
```

> the above request will return reponse like this:

```json
{
    "binance": {
        "binance_sub_account_1": [
            {
                "symbol": "LRCETH",
                "id": 2152067,
                "orderId": 83194209,
                "price": "0.00037412",
                "qty": "29318.00000000",
                "quote_qty": "",
                "commission": "0.14516094",
                "commissionAsset": "BNB",
                "time": 1597082787593,
                "isBuyer": false,
                "isMaker": false,
                "isBestMatch": true
            }
        ],
        "binance_sub_account_2": [
            {
                "symbol": "BNBUSDT",
                "id": 69038213,
                "orderId": 683611562,
                "price": "22.46610000",
                "qty": "11.00000000",
                "quote_qty": "",
                "commission": "0.00825000",
                "commissionAsset": "BNB",
                "time": 1597086408332,
                "isBuyer": true,
                "isMaker": false,
                "isBestMatch": true
            },
            {
                "symbol": "BNBUSDT",
                "id": 69038233,
                "orderId": 683611816,
                "price": "22.46100000",
                "qty": "11.00000000",
                "quote_qty": "",
                "commission": "0.00825000",
                "commissionAsset": "BNB",
                "time": 1597086416630,
                "isBuyer": true,
                "isMaker": false,
                "isBestMatch": true
            }
        ]
    },
    "huobi": {
        "id": 59378,
        "symbol": "ethusdt",
        "accountID": 100009,
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
}
```

### HTTP request

`GET http://gateway.local/trades`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | from time to get trades
to | integer | false | now | to time to get trades
cex | string | false | all | valid value: "binance", "huobi"


## Get convert to ETH price

```shell
curl -X GET "http://gateway.local/convert_to_eth_price?from=1494901162000&to=1499865549600"
```

> the above request will return reponse like this:

```json
[
    {
        "Symbol": "ETHBTC",
        "Price": 0.027756,
        "Timestamp": 1595471651724
    },
    {
        "Symbol": "ETHBTC",
        "Price": 0.02775,
        "Timestamp": 1595500928416
    }
]
```

### HTTP request

`GET http://gateway.local/convert_to_eth_price`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | true | no | from time to get trades
to | integer | true | now | to time to get trades
