# Centralized exchanges (cex) endpoint

## Get trades history 

```shell
curl -X GET "http://gateway.local/trades?from=1494901162000&to=1499865549600?cex=all"
```

> the above request will return reponse like this:

```json
{
    "binance": {
        "symbol": "BNBBTC",
        "id": 28457,
        "orderId": 100234,
        "price": "4.00000100",
        "qty": "12.00000000",
        "quote_qty": "48.000012",
        "Commission": "10.10000000",
        "commissionAsset": "BNB",
        "time": 1499865549590,
        "isBuyer": true,
        "isMaker": false,
        "isBestMatch": false
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