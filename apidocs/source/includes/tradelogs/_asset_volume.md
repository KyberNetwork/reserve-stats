## Asset volume

```shell
curl -X GET "http://localhost:8004/asset-volume?asset=ETH&from=1541548800000&to=1541635199999&freq=d"
```

> the above command returns JSON structure like this:

```json

```

This endpoint return the volume of a specific token in a time range **from** a point of time **to** another point of time

### HTTP request

`GET http://localhost:80004/asset-volume`

Params | Type | Required | Default | Description
------ | ---- | -------- | ------- | -----------
asset | string | true | null | 
from | integer | false | one hour before present | from stamp to get asset volume
to | integer | false | now | endpoint timestamp to get asset volume
freq | string | false | d | frequency of aggregation (d for day, and h for hour)