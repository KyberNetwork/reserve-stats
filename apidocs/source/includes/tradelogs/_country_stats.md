## Country stats

```shell
curl -X GET "http://gateway.local/country-stats?from=from=1541635200000&to=1541721600000&country=VN&timezone=0"
```

> the above request will return a struct like this

```json
{
    "1541548800000": {
        "eth_volume": 20.621640636171517,
        "usd_volume": 4519.457054957176,
        "burn_fee": 26.842690768136787,
        "total_trade": 15,
        "unique_addresses": 9,
        "kyced_addresses": 3,
        "new_unique_addresses": 3,
        "usd_per_trade": 301.29713699714506,
        "eth_per_trade": 1.3747760424114344
    },
    "1541635200000": {
        "eth_volume": 27.169015642380533,
        "usd_volume": 5950.443010131696,
        "burn_fee": 33.560748924036744,
        "total_trade": 14,
        "unique_addresses": 8,
        "kyced_addresses": 2,
        "new_unique_addresses": 1,
        "usd_per_trade": 425.0316435808354,
        "eth_per_trade": 1.9406439744557524
    }
}
```

Country stats will return stats of a country in provided time range

### HTTP Request

`GET http://gateway.local/country-stats`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | start time range
to | integer | false | now | end time range
country | string | true | coutry code (2 chars)
timezone | integer | 0 | time zone code from -11 to +13


## Heatmap

```shell
curl -X GET "http://gateway.local/heat-map?from=1542067200000&to=1542153600000&asset=0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
```

> the above request will return a struct like this

```json
{
    "unknown": {
        "country": "",
        "total_eth_value": 15.708794993071997,
        "total_token_value": 3011.0842747942766,
        "total_fiat_value": 3277.611123672339,
        "total_burn_fee": 0,
        "total_trade": 0,
        "total_unique_addr": 0,
        "total_kyc_user": 0
    }
}
```

Heatmap will return a sorted array of countries with biggest volume of ETH

### HTTP Request

`GET http://gateway.local/heat-map`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | start time range
to | integer | false | now | end time range 
asset | string | true | nil | address of token to get heatmap for