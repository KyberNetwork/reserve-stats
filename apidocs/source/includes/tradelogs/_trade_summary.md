## Trade summary

```shell
curl -X GET "http://gateway.local/trade-summary?from=1541635200000&to=1541721600000"
```

> the above request will return a struct like this

```json
{
    "1519862400000": {
        "total_eth_volume": 134.40205219911303,
        "total_usd_amount": 113753.22137532994,
        "total_burn_fee": 137.20624182764925,
        "total_trade": 235,
        "unique_addresses": 134,
        "kyced_addresses": 93,
        "new_unique_addresses": 100,
        "usd_per_trade": 484.05626117161677,
        "eth_per_trade": 0.5719236263792044
    },
    "1519948800000": {
        "total_eth_volume": 130.13009111525548,
        "total_usd_amount": 111349.55757172218,
        "total_burn_fee": 113.19461687215617,
        "total_trade": 230,
        "unique_addresses": 154,
        "kyced_addresses": 54,
        "new_unique_addresses": 104,
        "usd_per_trade": 484.12851118140077,
        "eth_per_trade": 0.5657830048489368
    }
}
```

Trade summary will summary of trade logs **from** a point of time **to** another point of time

### HTTP Request

`GET http://gateway.local/trade-summary`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
from | integer | false | one hour from now | start query time
to | integer | false | now | end query time

