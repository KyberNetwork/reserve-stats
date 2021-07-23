
## Get deposit history 

```shell
curl -X GET "http://gateway.local/deposits?from=1494901162000&to=1499865549600"
```

> the above request will return reponse like this:

```json
{
    "binance": {
        "binance_v2_sub_1": [
            {
                "amount": "99886.26346299",
                "coin": "ENJ",
                "network": "ETH",
                "status": 1,
                "address": "0x306b68c7954de1c14f5e6a1981e98e31aa98aec7",
                "addressTag": "",
                "txId": "0x107aa709846752b8ac1a5bb7ef1de2066dd76eea8b4df939256edce223bbfcc3",
                "insertTime": 1626870757000,
                "transferType": 0,
                "confirmTimes": "12/12"
            },
            {
                "amount": "124072.724806",
                "coin": "MANA",
                "network": "ETH",
                "status": 1,
                "address": "0x306b68c7954de1c14f5e6a1981e98e31aa98aec7",
                "addressTag": "",
                "txId": "0x1fa4273ebd48ce6c94ef35c5e25618f1bf457905576db1ea5a634091581e3622",
                "insertTime": 1626906151000,
                "transferType": 0,
                "confirmTimes": "12/12"
            }
        ]
    }
}
```

### HTTP request

`GET http://gateway.local/deposits`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | from time to get trades
to | integer | false | now | to time to get trades
cex | string | false | all | valid value: "binance", "huobi"