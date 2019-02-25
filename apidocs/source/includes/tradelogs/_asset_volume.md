## Asset volume

```shell
curl -X GET "http://gateway.local/asset-volume?asset=0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee&from=1541548800000&to=1541635199999&freq=d"
```

> the above command returns JSON structure like this:

```json
{
    "1541548800000": {
        "eth_amount": 1268.5853507197537,
        "usd_amount": 278024.2908058997,
        "volume": 1268.5853507197537
    },
    "1541635200000": {
        "eth_amount": 1444.1879872733625,
        "usd_amount": 316299.95091842714,
        "volume": 1444.1879872733625
    }
}
```

This endpoint return the volume of a specific token in a time range **from** a point of time **to** another point of time

### HTTP request

`GET http://gateway.local/asset-volume`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
asset | string | true | null | address of token to get asset volume for 
from | integer | false | one hour before present | from stamp to get asset volume
to | integer | false | now | endpoint timestamp to get asset volume
freq | string | false | h | frequency of aggregation (d for day, and h for hour)