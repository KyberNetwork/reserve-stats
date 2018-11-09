# Price Analytic

Price analytic component is for analytic module to update their action of settings rate then reserve stats. Reserve stats then will provide appropriate statistic number based on those information

## Update price analytic data

```shell
curl -X POST "http://gateway.local/price-analytics-data" \
-H 'Content-Type: application/json' \
-d '{
    "timestamp": 1541478232000,
    "block_expiration": true,
    "triggering_tokens_list": [
        {
            "token": "KNC",
            "ask_price": 0.123,
            "bid_price": 0.125,
            "mid_afp_price": 0.124,
            "mid_afp_old_price": 0.12,
            "min_spread": 0.002,
            "trigger_update": true
        },
        {
            "token": "OMG",
            "ask_price": 0.123,
            "bid_price": 0.125,
            "mid_afp_price": 0.124,
            "mid_afp_old_price": 0.12,
            "min_spread": 0.002,
            "trigger_update": false
        }
    ]
}'
```

> above request will response a struct like this
on success: http code 200

on error:

```json
{
    "error": <error>
}
```

### HTTP Request

`POST http://gateway.local/price-analytics`


## Get price analytic data

```shell
curl -X GET "http://gateway.local/price-analytic-data?from=1522753160000&to=1522755792000"
```

> sample response:

```json
[
    {
        "timestamp": 1522755271000,
        "block_expiration": false,
        "trigger_price_update": true,
        "triggering_tokens_list": [
            {
            "ask_price": 0.002,
            "bid_price": 0.003,
            "mid_afp_old_price": 0.34555,
            "mid_afp_price": 0.6555,
            "min_spread": 0.233,
            "token": "OMG"
            },
            {
            "ask_price": 0.004,
            "bid_price": 0.005,
            "mid_afp_old_price": 0.21555,
            "mid_afp_price": 0.4355,
            "min_spread": 0.133,
            "token": "KNC"
            }
        ]
    }
]
```

### HTTP Request

`GET http://gateway.local/price-analytics`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour before now | start timestamp to get price analytic data
to | integer | false | now | end timestamp to get price analytic data